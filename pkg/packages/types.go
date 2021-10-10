package packages

type Components struct {
	Major  int    `json:"major"`
	Minor  *int   `json:"minor"`
	Patch  *int   `json:"patch"`
	Suffix string `json:"suffix"`
}

type Version struct {
	Value      string      `json:"value"`
	Components *Components `json:"components"`
}

type Package struct {
	Owner   string  `json:"owner"`
	Repo    string  `json:"repo"`
	Version Version `json:"version"`
}

type Installation struct {
	Package  string   `json:"package"`
	Symlinks []string `json:"symlink"`
}

type Metadata struct {
	Package      Package      `json:"package"`
	Installation Installation `json:"installation"`
}
