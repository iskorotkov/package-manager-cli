package archives

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractTarGz(src string, dest string, permissions os.FileMode) error {
	if err := os.MkdirAll(dest, permissions); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("error creating folder for downloads: %w", err)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening archive file: %w", err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(srcFile)

	gzipReader, err := gzip.NewReader(srcFile)
	if err != nil {
		return fmt.Errorf("error creating archive reader: %w", err)
	}

	defer func(r *gzip.Reader) {
		_ = r.Close()
	}(gzipReader)

	tarReader := tar.NewReader(gzipReader)

	if err := extractEntries(tarReader, dest, permissions); err != nil {
		return err
	}

	return nil
}

func extractEntries(tarReader *tar.Reader, dest string, permissions os.FileMode) error {
	for {
		h, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("error extracting file header: %w", err)
		}

		entryPath, err := sanitizeExtractPath(dest, h.Name)
		if err != nil {
			return err
		}

		switch h.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(entryPath, permissions); err != nil {
				return fmt.Errorf("error creating folder: %w", err)
			}
		case tar.TypeReg:
			if err := extractFile(tarReader, entryPath); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected type flag '%v' when extracting file '%s' in tar archive ",
				h.Typeflag, h.Name)
		}
	}

	return nil
}

func extractFile(tarReader *tar.Reader, dest string) error {
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	if _, err := io.Copy(destFile, tarReader); err != nil {
		return fmt.Errorf("error copying file contents to the dest folder: %w", err)
	}

	return nil
}

// sanitizeExtractPath helps to avoid Zip Slip vulnerability (https://snyk.io/research/zip-slip-vulnerability).
func sanitizeExtractPath(folder string, file string) (string, error) {
	p := filepath.Join(folder, file)

	if !strings.HasPrefix(p, filepath.Clean(folder)+string(os.PathSeparator)) {
		return "", fmt.Errorf("archive traversal vulnerability detected")
	}

	return p, nil
}
