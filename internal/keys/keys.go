package keys

import (
	"os"

	"github.com/iskorotkov/package-manager-cli/pkg/env"
)

//nolint:gochecknoglobals,gomnd,gofumpt
var (
	DownloadsPath = env.Get("PM_DOWNLOADS_PATH", "./local/share/package-manager/downloads")
	SymlinksPath  = env.Get("PM_SYMLINKS_PATH", "./local/bin")

	DownloadsPermissions = os.FileMode(env.GetInt("PM_DOWNLOADS_PERMISSIONS", 0744))
	SymlinksPermissions  = os.FileMode(env.GetInt("PM_SYMLINKS_PERMISSIONS", 0744))
)
