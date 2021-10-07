package assets

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v39/github"
)

type Downloader struct {
	client *github.Client
}

func NewDownloader(client *github.Client) *Downloader {
	return &Downloader{client: client}
}

func (d Downloader) Download(
	ctx context.Context,
	repo *github.Repository,
	asset *github.ReleaseAsset,
	dest string,
) error {
	rc, _, err := d.client.Repositories.DownloadReleaseAsset(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		asset.GetID(),
		http.DefaultClient,
	)
	if err != nil {
		return fmt.Errorf("error downloading release asset: %w", err)
	}

	defer func(rc io.ReadCloser) {
		_ = rc.Close()
	}(rc)

	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if _, err := io.Copy(file, rc); err != nil {
		return fmt.Errorf("error copying file contents to the dest folder: %w", err)
	}

	return nil
}
