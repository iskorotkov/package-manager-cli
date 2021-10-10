package keys

import (
	"os"

	"github.com/iskorotkov/package-manager-cli/pkg/env"
)

//nolint:gochecknoglobals,gomnd,gofumpt
var (
	DownloadsPath = env.Get("PM_DOWNLOADS_PATH", "~/.local/share/package-manager/downloads")
	PackagesPath  = env.Get("PM_PACKAGES_PATH", "~/.local/share/package-manager/packages")
	MetadataPath  = env.Get("PM_METADATA_PATH", "~/.local/share/package-manager/metadata")
	SymlinksPath  = env.Get("PM_SYMLINKS_PATH", "~/.local/bin")

	DownloadsPermissions = os.FileMode(env.GetInt("PM_DOWNLOADS_PERMISSIONS", 0744))
	PackagesPermissions  = os.FileMode(env.GetInt("PM_PACKAGES_PERMISSIONS", 0744))
	MetadataPermissions  = os.FileMode(env.GetInt("PM_METADATA_PERMISSIONS", 0744))
	SymlinksPermissions  = os.FileMode(env.GetInt("PM_SYMLINKS_PERMISSIONS", 0744))
)
