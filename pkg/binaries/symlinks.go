package binaries

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func AddSymlinks(src, dest string, permissions os.FileMode) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("error reading contents of src folder: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == "bin" {
				if err := addSymlinksToBin(src, dest, entry, permissions); err != nil {
					return err
				}
			}

			continue
		}

		ext := filepath.Ext(entry.Name())
		if shouldSkipExtension(ext) {
			continue
		}

		lowerName := strings.ToLower(entry.Name())
		if shouldSkipFile(lowerName) {
			continue
		}

		if err := createSymlink(src, dest, entry, permissions); errors.Is(err, os.ErrExist) {
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

func addSymlinksToBin(src string, dest string, entry os.DirEntry, permissions os.FileMode) error {
	binFiles, err := os.ReadDir(entry.Name())
	if err != nil {
		return fmt.Errorf("error reading contents of bin folder: %w", err)
	}

	for _, binEntry := range binFiles {
		if binEntry.IsDir() {
			continue
		}

		if err := createSymlink(src, filepath.Join(dest, "bin"), binEntry, permissions); errors.Is(err, os.ErrExist) {
			printSymlinkExists(entry)
		} else if err != nil {
			return err
		}
	}

	return nil
}

func printSymlinkExists(entry os.DirEntry) {
	fmt.Printf("%s: symlink already exists", entry.Name())
}

func createSymlink(src string, dest string, entry os.DirEntry, permissions os.FileMode) error {
	oldName, err := filepath.Abs(filepath.Join(src, entry.Name()))
	if err != nil {
		return fmt.Errorf("error getting abs path for src: %w", err)
	}

	if err := os.Chmod(oldName, permissions); err != nil {
		return fmt.Errorf("error changing permissions for file '%s': %w", oldName, err)
	}

	newName := filepath.Join(dest, entry.Name())

	if err := os.Symlink(oldName, newName); errors.Is(err, os.ErrExist) {
		return fmt.Errorf("symlink '%s' already exists: %w", newName, err)
	} else if err != nil {
		return fmt.Errorf("error creating symlink from '%s' to '%s': %w", oldName, newName, err)
	}

	return nil
}
