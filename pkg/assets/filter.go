package assets

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v39/github"
)

const (
	OSAny     = OS("*")
	OSUnknown = OS("?")
	OSWindows = OS("windows")
	OSLinux   = OS("linux")
	OSMac     = OS("mac")

	ArchAny     = Arch("*")
	ArchUnknown = Arch("?")
	ArchARM86   = Arch("arm")
	ArchARM64   = Arch("arm64")
	ArchX86     = Arch("x86")
	ArchX64     = Arch("x64")
	ArchPPC64   = Arch("ppc64")
	ArchPPC64LE = Arch("ppc64le")
)

type OS string

type Arch string

type Platform struct {
	OS
	Arch
}

type assetMetadata struct {
	*github.ReleaseAsset
	Platform
}

func ForPlatform(assets []*github.ReleaseAsset, platforms []Platform) (*github.ReleaseAsset, error) {
	metadata := make([]assetMetadata, 0, len(assets))

	for _, a := range assets {
		name := strings.ToLower(a.GetName())

		m := assetMetadata{
			ReleaseAsset: a,
			Platform: Platform{
				OS:   selectOS(name),
				Arch: selectArch(name),
			},
		}

		metadata = append(metadata, m)
	}

	var filtered []assetMetadata

	for _, p := range platforms {
		for _, m := range metadata {
			if p.OS != OSAny && p.OS != m.OS {
				continue
			}

			if p.Arch != ArchAny && p.Arch != m.Arch {
				continue
			}

			filtered = append(filtered, m)
		}

		if len(filtered) > 0 {
			return filtered[0].ReleaseAsset, nil
		}
	}

	return nil, fmt.Errorf("no assets available for this platform and arch")
}

func selectArch(name string) Arch {
	switch {
	case strings.Contains(name, "arm64"):
		return ArchARM64
	case strings.Contains(name, "arm"):
		return ArchARM86
	case strings.Contains(name, "ppc64le"):
		return ArchPPC64LE
	case strings.Contains(name, "ppc64"):
		return ArchPPC64
	case strings.Contains(name, "x64"), strings.Contains(name, "x86_64"), strings.Contains(name, "x86-64"):
		return ArchX64
	case strings.Contains(name, "x86"):
		return ArchX86
	}

	return ArchUnknown
}

func selectOS(name string) OS {
	switch {
	case strings.Contains(name, "linux"):
		return OSLinux
	case strings.Contains(name, "mac"), strings.Contains(name, "osx"), strings.Contains(name, "darwin"):
		return OSMac
	case strings.Contains(name, "win"):
		return OSWindows
	}

	return OSUnknown
}
