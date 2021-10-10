package binaries

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func AddSymlinks(src, dest string, permissions os.FileMode) ([]string, error) {
	if err := os.MkdirAll(dest, permissions); err != nil {
		return nil, fmt.Errorf("error creating folder '%s' for symlinks: %w", dest, err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return nil, fmt.Errorf("error reading contents of src folder: %w", err)
	}

	var created []string

	for _, entry := range entries {
		entrySymlinks, err := analyzeSourceEntry(src, dest, permissions, entry)
		if err != nil {
			return created, err
		}

		created = append(created, entrySymlinks...)
	}

	return created, nil
}

func analyzeSourceEntry(src string, dest string, permissions os.FileMode, entry os.DirEntry) ([]string, error) {
	ext := filepath.Ext(entry.Name())
	lowerName := strings.ToLower(entry.Name())

	var created []string

	switch {
	case entry.IsDir() && entry.Name() == "bin":
		binSymlinks, err := addSymlinksToBin(src, dest, entry, permissions)
		if err != nil {
			return created, err
		}

		created = append(created, binSymlinks...)
	case entry.IsDir(), shouldSkipExtension(ext), shouldSkipFile(lowerName):
	default:
		src := filepath.Join(src, entry.Name())
		dest := filepath.Join(dest, entry.Name())

		if err := createSymlink(src, dest, permissions); errors.Is(err, os.ErrExist) {
			printSymlinkExists(entry)
		} else if err != nil {
			return created, err
		}

		created = append(created, dest)
	}

	return created, nil
}

func shouldSkipExtension(ext string) bool {
	return ext != "" && ext != ".sh"
}

func shouldSkipFile(lowerName string) bool {
	return strings.Contains(lowerName, "readme") || strings.Contains(lowerName, "license")
}

func addSymlinksToBin(src string, dest string, dir os.DirEntry, permissions os.FileMode) ([]string, error) {
	binFiles, err := os.ReadDir(dir.Name())
	if err != nil {
		return nil, fmt.Errorf("error reading contents of bin folder: %w", err)
	}

	var created []string //nolint:prealloc

	for _, entry := range binFiles {
		if entry.IsDir() {
			continue
		}

		src := filepath.Join(src, "bin", entry.Name())
		dest := filepath.Join(dest, entry.Name())

		if err := createSymlink(src, dest, permissions); errors.Is(err, os.ErrExist) {
			printSymlinkExists(dir)
		} else if err != nil {
			return created, err
		}

		created = append(created, dest)
	}

	return created, nil
}

func printSymlinkExists(entry os.DirEntry) {
	fmt.Printf("%s: symlink already exists\n", entry.Name())
}

func createSymlink(src string, dest string, permissions os.FileMode) error {
	src, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("error getting abs path for src: %w", err)
	}

	if err := os.Chmod(src, permissions); err != nil {
		return fmt.Errorf("error changing permissions for file '%s': %w", src, err)
	}

	if err := os.Symlink(src, dest); errors.Is(err, os.ErrExist) {
		return fmt.Errorf("symlink '%s' already exists: %w", dest, err)
	} else if err != nil {
		return fmt.Errorf("error creating symlink from '%s' to '%s': %w", src, dest, err)
	}

	return nil
}
