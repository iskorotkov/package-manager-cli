package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/iskorotkov/package-manager-cli/internal/keys"
	"github.com/iskorotkov/package-manager-cli/pkg/archives"
	"github.com/iskorotkov/package-manager-cli/pkg/assets"
	"github.com/iskorotkov/package-manager-cli/pkg/binaries"
	"github.com/spf13/cobra"
)

//nolint:gochecknoinits
func init() {
	installCmd := &cobra.Command{ //nolint:exhaustivestruct
		Use:   "install",
		Short: "install package",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]

			client := github.NewClient(nil)

			repo, asset, err := selectAsset(client, packageName)
			if err != nil {
				return err
			}

			downloadDest := filepath.Join(keys.DownloadsPath, asset.GetName())
			if err := downloadAsset(client, repo, asset, downloadDest); err != nil {
				return err
			}

			defer cleanup(downloadDest)

			extractDest := filepath.Join(keys.DownloadsPath, repo.GetName())

			if strings.HasSuffix(asset.GetName(), ".tar.gz") {
				if err := archives.ExtractTarGz(downloadDest, extractDest, keys.DownloadsPermissions); err != nil {
					return fmt.Errorf("error extracting tar.gz file: %w", err)
				}
			} else {
				if err := moveFileToPackageFolder(downloadDest, extractDest, repo); err != nil {
					return err
				}
			}

			if err := os.MkdirAll(keys.SymlinksPath, keys.DownloadsPermissions); err != nil {
				return fmt.Errorf("error creating folder '%s' for symlinks: %w", keys.SymlinksPath, err)
			}

			if err := binaries.AddSymlinks(extractDest, keys.SymlinksPath, keys.SymlinksPermissions); err != nil {
				return fmt.Errorf("error adding package to path: %w", err)
			}

			return nil
		},
	}

	rootCmd.AddCommand(installCmd)
}

func cleanup(downloadDest string) {
	func() {
		if err := os.Remove(downloadDest); err != nil {
			fmt.Printf("error removing downloaded file '%s': %v\n", downloadDest, err)
		}
	}()
}

func moveFileToPackageFolder(src string, dest string, repo *github.Repository) error {
	if err := os.Mkdir(dest, keys.DownloadsPermissions); err != nil {
		return fmt.Errorf("error creating folder for package: %w", err)
	}

	from, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening downloaded file: %w", err)
	}

	to, err := os.Create(filepath.Join(dest, repo.GetName()))
	if err != nil {
		return fmt.Errorf("error creating a new file in dest folder: %w", err)
	}

	if _, err := io.Copy(to, from); err != nil {
		return fmt.Errorf("error copying downloaded file in package folder: %w", err)
	}

	return nil
}

func selectAsset(client *github.Client, packageName string) (*github.Repository, *github.ReleaseAsset, error) {
	result, _, err := client.Search.Repositories(context.Background(), packageName, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error searching repositories: %w", err)
	}

	if len(result.Repositories) == 0 {
		return nil, nil, fmt.Errorf("no results")
	}

	repo := result.Repositories[0]

	releases, _, err := client.Repositories.ListReleases(
		context.Background(),
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting releases: %w", err)
	}

	if len(releases) == 0 {
		return nil, nil, fmt.Errorf("no releases available")
	}

	release := releases[0]

	asset, err := assets.ForPlatform(release.Assets, getPlatforms())
	if err != nil {
		return nil, nil, fmt.Errorf("no assets available: %w", err)
	}

	return repo, asset, nil
}

func downloadAsset(client *github.Client, repo *github.Repository, asset *github.ReleaseAsset, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), keys.DownloadsPermissions); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("error creating folder for downloads: %w", err)
	}

	downloader := assets.NewDownloader(client)
	if err := downloader.Download(context.Background(), repo, asset, dest); err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}

	return nil
}

func getPlatforms() []assets.Platform {
	return []assets.Platform{
		{OS: assets.OSLinux, Arch: assets.ArchX64},
		{OS: assets.OSLinux, Arch: assets.ArchX86},
		{OS: assets.OSLinux, Arch: assets.ArchUnknown},
	}
}
