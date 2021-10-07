package packages

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major  int
	Minor  *int
	Patch  *int
	Suffix string
}

type Package struct {
	Owner, Repo string
	Version     *Version
}

//nolint:gomnd
//goland:noinspection GoUnusedExportedFunction
func ParsePackage(name string) (Package, error) {
	var username, repo, versionStr string

	ss := strings.Split(name, "@")
	if len(ss) > 2 {
		return Package{}, fmt.Errorf("package name can't contain more than one @")
	}

	if len(ss) == 2 {
		versionStr = ss[1]
	}

	ss = strings.Split(ss[0], "/")
	if len(ss) > 2 {
		return Package{}, fmt.Errorf("package name can't contain more than one /")
	}

	username = ss[0]

	if len(ss) == 2 {
		repo = ss[1]
	} else {
		// Use the same name for username and repo (e. g. my-package/my-package).
		repo = ss[0]
	}

	version, err := ParseVersion(versionStr)
	if err != nil {
		return Package{}, fmt.Errorf("error parsing version: %w", err)
	}

	return Package{
		Owner:   username,
		Repo:    repo,
		Version: version,
	}, nil
}

//nolint:gomnd
func ParseVersion(version string) (*Version, error) {
	if version == "" {
		return nil, nil
	}

	trimmed := strings.TrimPrefix(version, "v")
	parts := strings.Split(trimmed, "-")

	mainPart, suffix := "", ""
	if len(parts) > 1 {
		suffix = strings.Join(parts[1:], "-")
	}

	mainParts := strings.Split(mainPart, ".")

	major, err := strconv.Atoi(mainParts[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing major version: %w", err)
	}

	var minor, patch *int

	if len(mainParts) >= 2 {
		i, err := strconv.Atoi(mainParts[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing minor version: %w", err)
		}

		minor = &i
	}

	if len(mainParts) >= 2 {
		i, err := strconv.Atoi(mainParts[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing patch version: %w", err)
		}

		patch = &i
	}

	return &Version{
		Major:  major,
		Minor:  minor,
		Patch:  patch,
		Suffix: suffix,
	}, nil
}
