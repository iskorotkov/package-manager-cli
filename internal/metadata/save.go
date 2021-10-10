package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iskorotkov/package-manager-cli/pkg/assets"
	"github.com/iskorotkov/package-manager-cli/pkg/packages"
)

func Save(src, dest string, asset assets.AssetData, symlinks []string, permissions os.FileMode) error {
	if err := os.MkdirAll(dest, permissions); err != nil {
		return fmt.Errorf("error creating folder '%s' for metadata: %w", dest, err)
	}

	m := packages.Metadata{
		Package: packages.Package{
			Owner: asset.Repository.GetOwner().GetLogin(),
			Repo:  asset.Repository.GetName(),
			Version: packages.Version{ //nolint:exhaustivestruct
				// TODO: Parse version values.
				Value: asset.Release.GetTagName(),
			},
		},
		Installation: packages.Installation{
			Package:  src,
			Symlinks: symlinks,
		},
	}

	b, err := json.MarshalIndent(&m, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling package metadata: %w", err)
	}

	metadataPath := filepath.Join(dest, asset.Repository.GetName())

	if err := os.WriteFile(metadataPath, b, permissions); err != nil {
		return fmt.Errorf("error writing metadata file '%s': %w", metadataPath, err)
	}

	return nil
}
