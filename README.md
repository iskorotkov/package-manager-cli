# Package Manager CLI

Package manager for installing GitHub releases.

It's a CLI tool that lets you search and install other CLI tools from GitHub from your console.

If you need to get the latest version of `k9s`, `buf`, `kubectx`, `minikube`, `legendary` or other package, you don't have to make changes to system installation (which may not be possible) or search for the latest release on GitHub, select a correct file for your OS/architecture, manually extract and install it.

All you have to do is to use `pmcli install minikube`, and it's done!

## Installation

1. Download `pmcli` from releases page and place it in `~/.local/share/bin` folder.
2. Run `chmod +x ~/.local/share/bin/pmcli` to make file executable.
3. Run `pmcli` to make sure everything is working.

NOTE: Make sure that `~/.local/share/bin` is added to `PATH` environment variable.

## Usage

Search for minikube on GitHub:

```shell
pmcli search minikube
```

---

See info about minikube package:

```shell
pmcli info minikube
```

NOTE: It internally makes a search request, takes the first match (i. e. what GitHub considers to be the best match) and shows info about it.

---

Install minikube for your OS/arch:

```shell
pmcli install minikube
```

NOTE: It internally makes a search request, takes the first match (i. e. what GitHub considers to be the best match) and installs the latest release.

NOTE: It extracts .tar.gz archive if necessary, then creates symlinks to binaries in ~/.local/share/bin folder. It can skip a binary file if it already exist there.

NOTE: You can use `pmcli install {owner}/{repo}` instead of shorter version `pmcli install {repo}` if the latter doesn't pick the correct repo.

---

Uninstall minikube package (if it's installed):

```shell
pmcli uninstall minikube
```

---

List installed packages:

```shell
pmcli list
```

## Configuration

See `internal/keys/keys.go` for all values that can be configured via environment variables.

## Bugs, ideas and contribution

If you've encountered a bug or want to request a feature, feel free to open an issue on GitHub.
