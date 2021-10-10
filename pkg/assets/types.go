package assets

import (
	"github.com/google/go-github/v39/github"
)

type AssetData struct {
	Repository *github.Repository
	Release    *github.RepositoryRelease
	Asset      *github.ReleaseAsset
}
