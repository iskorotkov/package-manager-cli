package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/iskorotkov/package-manager-cli/internal/keys"
	"github.com/iskorotkov/package-manager-cli/internal/metadata"
	"github.com/iskorotkov/package-manager-cli/pkg/archives"
	"github.com/iskorotkov/package-manager-cli/pkg/assets"
	"github.com/iskorotkov/package-manager-cli/pkg/binaries"
	"github.com/iskorotkov/package-manager-cli/pkg/xlog"
	"github.com/spf13/cobra"
)

//nolint:gochecknoinits
func init() {
	installCmd := wrapCommand(&cobra.Command{ //nolint:exhaustivestruct
		Use:   "install",
		Short: "install package",
		Args:  cobra.MinimumNArgs(1),
		RunE:  install,
	})

	rootCmd.AddCommand(installCmd)
}

func install(_ *cobra.Command, args []string) error {
	packageName := args[0]

	xlog.Push(packageName)
	defer xlog.Pop()

	log.Printf("package name: %s", packageName)

	client := github.NewClient(nil)

	asset, err := selectAsset(client, packageName)
	if err != nil {
		return err
	}

	log.Printf("selected repo: %s", asset.Repository.GetFullName())
	log.Printf("selected release: %s", asset.Release.GetTagName())
	log.Printf("selected asset: %s", asset.Asset.GetName())

	downloadPath := filepath.Join(keys.DownloadsPath, asset.Asset.GetName())

	log.Printf("downloading package to: %s", downloadPath)

	if err := downloadAsset(client, asset, downloadPath); err != nil {
		return err
	}

	defer cleanupFile(downloadPath)

	log.Printf("downloaded to: %s", downloadPath)

	packagePath := filepath.Join(keys.PackagesPath, asset.Repository.GetName())

	log.Printf("moving package to: %s", packagePath)

	if err := moveToPackageFolder(asset, downloadPath, packagePath); err != nil {
		return err
	}

	log.Printf("creating symlinks at: %s", keys.SymlinksPath)

	symlinks, err := binaries.AddSymlinks(packagePath, keys.SymlinksPath, keys.SymlinksPermissions)
	if err != nil {
		return fmt.Errorf("error adding package to path: %w", err)
	}

	log.Printf("saved symlinks: %+v", symlinks)

	log.Printf("saving metadata to: %s", keys.MetadataPath)

	err = metadata.Save(packagePath, keys.MetadataPath, asset, symlinks, keys.MetadataPermissions)
	if err != nil {
		return fmt.Errorf("error saving package metadata: %w", err)
	}

	return nil
}

func moveToPackageFolder(asset assets.AssetData, downloadPath string, packagePath string) error {
	if strings.HasSuffix(asset.Asset.GetName(), ".tar.gz") {
		xlog.Push("archive")
		defer xlog.Pop()

		log.Printf("extracting .tar.gz archive")

		if err := archives.ExtractTarGz(downloadPath, packagePath, keys.PackagesPermissions); err != nil {
			return fmt.Errorf("error extracting tar.gz file: %w", err)
		}
	} else {
		xlog.Push("binary")
		defer xlog.Pop()

		log.Printf("moving binary file to package folder")

		err := moveFileToPackageFolder(downloadPath, packagePath, keys.PackagesPermissions, asset.Repository)
		if err != nil {
			return err
		}
	}

	return nil
}

// cleanupFile removes file if it still exists.
// It is useful to call after package installation,
// and it will ignore cases where downloaded file was moved somewhere else.
func cleanupFile(file string) {
	if err := os.Remove(file); err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("error removing downloaded file '%s': %v\n", file, err)
	}
}

func moveFileToPackageFolder(src string, dest string, permissions os.FileMode, repo *github.Repository) error {
	if err := os.MkdirAll(dest, permissions); err != nil {
		return fmt.Errorf("error creating package folder: %w", err)
	}

	if err := os.Rename(src, filepath.Join(dest, repo.GetName())); err != nil {
		return fmt.Errorf("error moving file to package folder: %w", err)
	}

	return nil
}

func selectAsset(client *github.Client, packageName string) (assets.AssetData, error) {
	result, _, err := client.Search.Repositories(context.Background(), packageName, nil)
	if err != nil {
		return assets.AssetData{}, fmt.Errorf("error searching repositories: %w", err)
	}

	if len(result.Repositories) == 0 {
		return assets.AssetData{}, fmt.Errorf("no results")
	}

	repo := result.Repositories[0]

	releases, _, err := client.Repositories.ListReleases(
		context.Background(),
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		nil,
	)
	if err != nil {
		return assets.AssetData{}, fmt.Errorf("error getting releases: %w", err)
	}

	if len(releases) == 0 {
		return assets.AssetData{}, fmt.Errorf("no releases available")
	}

	release := releases[0]

	asset, err := assets.ForPlatform(release.Assets, getPlatforms())
	if err != nil {
		return assets.AssetData{}, fmt.Errorf("no assets available: %w", err)
	}

	return assets.AssetData{
		Repository: repo,
		Release:    release,
		Asset:      asset,
	}, nil
}

func downloadAsset(client *github.Client, asset assets.AssetData, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), keys.DownloadsPermissions); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("error creating folder for downloads: %w", err)
	}

	downloader := assets.NewDownloader(client)
	if err := downloader.Download(context.Background(), asset.Repository, asset.Asset, dest); err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}

	return nil
}

func getPlatforms() []assets.Platform {
	return []assets.Platform{
		{OS: assets.OSLinux, Arch: assets.ArchX64},
		{OS: assets.OSLinux, Arch: assets.ArchX86},
		{OS: assets.OSLinux, Arch: assets.ArchUnknown},
		{OS: assets.OSUnknown, Arch: assets.ArchUnknown},
	}
}
