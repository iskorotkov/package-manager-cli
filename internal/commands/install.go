package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/iskorotkov/package-manager-cli/pkg/archives"
	"github.com/iskorotkov/package-manager-cli/pkg/assets"
	"github.com/iskorotkov/package-manager-cli/pkg/binaries"
	"github.com/spf13/cobra"
)

//nolint:gofumpt
const (
	permissions   = 0744
	downloadsDest = "./temp/downloads"
	symlinkDest   = "./temp/symlinks"
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

			downloadDest := filepath.Join(downloadsDest, asset.GetName())
			if err := downloadAsset(client, repo, asset, downloadDest); err != nil {
				return err
			}

			extractDest := filepath.Join(downloadsDest, repo.GetName())

			if strings.HasSuffix(asset.GetName(), ".tar.gz") {
				if err := archives.ExtractTarGz(downloadDest, extractDest, permissions); err != nil {
					return fmt.Errorf("error extracting tar.gz file: %w", err)
				}
			}

			if err := os.MkdirAll(symlinkDest, permissions); err != nil {
				return fmt.Errorf("error creating folder '%s' for symlinks: %w", symlinkDest, err)
			}

			if err := binaries.AddSymlinks(extractDest, symlinkDest); err != nil {
				return fmt.Errorf("error adding package to path: %w", err)
			}

			return nil
		},
	}

	rootCmd.AddCommand(installCmd)
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
	if err := os.MkdirAll(filepath.Dir(dest), permissions); err != nil && !errors.Is(err, os.ErrExist) {
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
