package binaries

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func AddSymlinks(src, dest string, permissions os.FileMode) error {
	if err := os.MkdirAll(dest, permissions); err != nil {
		return fmt.Errorf("error creating folder '%s' for symlinks: %w", dest, err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("error reading contents of src folder: %w", err)
	}

	for _, entry := range entries {
		if err := analyzeSourceEntry(src, dest, permissions, entry); err != nil {
			return err
		}
	}

	return nil
}

func analyzeSourceEntry(src string, dest string, permissions os.FileMode, entry os.DirEntry) error {
	ext := filepath.Ext(entry.Name())
	lowerName := strings.ToLower(entry.Name())

	switch {
	case entry.IsDir() && entry.Name() == "bin":
		if err := addSymlinksToBin(src, dest, entry, permissions); err != nil {
			return err
		}
	case entry.IsDir(), shouldSkipExtension(ext), shouldSkipFile(lowerName):
	default:
		err := createSymlink(filepath.Join(src, entry.Name()), filepath.Join(dest, entry.Name()), permissions)
		if errors.Is(err, os.ErrExist) {
			printSymlinkExists(entry)
		} else if err != nil {
			return err
		}
	}

	return nil
}

func shouldSkipExtension(ext string) bool {
	return ext != "" && ext != ".sh"
}

func shouldSkipFile(lowerName string) bool {
	return strings.Contains(lowerName, "readme") || strings.Contains(lowerName, "license")
}

func addSymlinksToBin(src string, dest string, dir os.DirEntry, permissions os.FileMode) error {
	binFiles, err := os.ReadDir(dir.Name())
	if err != nil {
		return fmt.Errorf("error reading contents of bin folder: %w", err)
	}

	for _, entry := range binFiles {
		if entry.IsDir() {
			continue
		}

		err := createSymlink(filepath.Join(src, entry.Name()), filepath.Join(dest, "bin", entry.Name()), permissions)
		if errors.Is(err, os.ErrExist) {
			printSymlinkExists(dir)
		} else if err != nil {
			return err
		}
	}

	return nil
}

func printSymlinkExists(entry os.DirEntry) {
	fmt.Printf("%s: symlink already exists", entry.Name())
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
