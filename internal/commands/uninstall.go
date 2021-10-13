package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/iskorotkov/package-manager-cli/internal/keys"
	"github.com/iskorotkov/package-manager-cli/internal/metadata"
	"github.com/iskorotkov/package-manager-cli/pkg/packages"
	"github.com/iskorotkov/package-manager-cli/pkg/xlog"
	"github.com/spf13/cobra"
)

//nolint:gochecknoinits
func init() {
	uninstallCmd := wrapCommand(&cobra.Command{ //nolint:exhaustivestruct
		Use:   "uninstall",
		Short: "uninstall package",
		Args:  cobra.MinimumNArgs(1),
		RunE:  uninstall,
	})

	rootCmd.AddCommand(uninstallCmd)
}

func uninstall(_ *cobra.Command, args []string) error {
	packageName := args[0]

	xlog.Push(packageName)
	defer xlog.Pop()

	log.Printf("package name: %s", packageName)

	dir, err := ioutil.ReadDir(keys.MetadataPath)
	if errors.Is(err, os.ErrNotExist) {
		printPackageNotInstalled(packageName)

		return nil
	} else if err != nil {
		return fmt.Errorf("error opening metadata folder: %w", err)
	}

	if len(dir) == 0 {
		printPackageNotInstalled(packageName)

		return nil
	}

	pkg, err := packages.ParsePackage(packageName)
	if err != nil {
		return fmt.Errorf("error parsing package name: %w", err)
	}

	log.Printf("pkg package name as package metadata: %+v", pkg)

	packageMetadata, path, err := findPackageMetadata(pkg, dir)
	if err != nil {
		return err
	}

	if packageMetadata == nil {
		printPackageNotInstalled(packageName)

		return nil
	}

	if path == "" {
		return fmt.Errorf("empty metadata file path")
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error removing metadata file: %w", err)
	}

	log.Printf("removing package: %+v", packageMetadata)

	if err := removePackage(packageMetadata); err != nil {
		return err
	}

	fmt.Printf("uninstalled package '%s/%s'", packageMetadata.Package.Owner, packageMetadata.Package.Repo)

	return nil
}

func findPackageMetadata(pkg packages.Package, files []fs.FileInfo) (*packages.Metadata, string, error) {
	for _, file := range files {
		if file.Name() != pkg.Repo {
			log.Printf("skipping file due to filename mismatch: %s", file.Name())

			continue
		}

		path := filepath.Join(keys.MetadataPath, file.Name())

		meta, err := metadata.Read(path)
		if err != nil {
			return nil, "", fmt.Errorf("error reading package metadata: %w", err)
		}

		if meta.Package.Repo != pkg.Repo {
			log.Printf("skipping file due to package metadata mismatch: %+v", meta.Package)

			continue
		}

		return &meta, path, nil
	}

	return nil, "", nil
}

func removePackage(packageMetadata *packages.Metadata) error {
	if err := os.RemoveAll(packageMetadata.Installation.Package); err != nil {
		return fmt.Errorf("error removing package folder: %w", err)
	}

	for _, symlink := range packageMetadata.Installation.Symlinks {
		if err := os.Remove(symlink); err != nil {
			return fmt.Errorf("error removing symlink: %w", err)
		}
	}

	return nil
}

func printPackageNotInstalled(name string) {
	fmt.Printf("package '%s' isn't installed\n", name)
	log.Printf("package isn't installed: %s", name)
}
