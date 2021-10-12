package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iskorotkov/package-manager-cli/internal/keys"
	"github.com/iskorotkov/package-manager-cli/internal/metadata"
	"github.com/iskorotkov/package-manager-cli/pkg/xlog"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

//nolint:gochecknoinits
func init() {
	infoCmd := wrapCommand(&cobra.Command{ //nolint:exhaustivestruct
		Use:   "list",
		Short: "list installed packages",
		Args:  cobra.NoArgs,
		RunE:  list,
	})

	rootCmd.AddCommand(infoCmd)
}

func list(_ *cobra.Command, _ []string) error {
	dir, err := ioutil.ReadDir(keys.MetadataPath)
	if errors.Is(err, os.ErrNotExist) {
		printNoPackagesInstalled()

		return nil
	} else if err != nil {
		return fmt.Errorf("error opening metadata folder: %w", err)
	}

	if len(dir) == 0 {
		printNoPackagesInstalled()

		return nil
	}

	t := createTable()
	t.AppendHeader(table.Row{"repo", "version", "binaries"})

	packages := make([]string, 0, len(dir))

	for _, file := range dir {
		packages = append(packages, file.Name())
	}

	log.Printf("installed packages: %+v", packages)

	for _, name := range packages {
		path := filepath.Join(keys.MetadataPath, name)

		if err := addPackageRow(t, path, name); err != nil {
			return err
		}
	}

	fmt.Printf("%d packages installed\n", len(dir))

	t.Render()

	return nil
}

func addPackageRow(t table.Writer, path string, name string) error {
	xlog.Push(name)
	defer xlog.Pop()

	log.Printf("package metadata at path: %s", path)

	m, err := metadata.Read(path)
	if err != nil {
		return fmt.Errorf("error reading package metadata: %w", err)
	}

	binaries := make([]string, 0, len(m.Installation.Symlinks))

	for _, b := range m.Installation.Symlinks {
		binaries = append(binaries, filepath.Base(b))
	}

	log.Printf("package binaries: %+v", binaries)

	t.AppendRow(table.Row{
		fmt.Sprintf("%s/%s", m.Package.Owner, m.Package.Repo),
		m.Package.Version.Value,
		strings.Join(binaries, ", "),
	})

	return nil
}

func printNoPackagesInstalled() {
	fmt.Println("no packages installed")
	log.Printf("no packages installed")
}
