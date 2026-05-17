# gom

Go package manager. Track Go-installed tools in a single packages file and keep them in sync across machines.

## Install

```sh
go install github.com/imkk000/gom@latest
```

## Usage

```sh
gom install <pkg>      Install a package and record it
gom uninstall <pkg>    Remove a package binary and record
gom list               List packages with install status
gom sync               Mirror packages.json: install missing, remove extras
gom upgrade            Reinstall all packages at @latest
gom export             Scan GOPATH/bin and write all Go packages into packages.json (with diff)
```

Every install runs `go install -trimpath -ldflags="-s -w"`.

## Packages file

Location:

- `$XDG_CONFIG_HOME/gom/packages.json`
- fallback `~/.config/gom/packages.json`

Format:

```json
["github.com/foo/bar", "github.com/baz/qux"]
```

## Example

```sh
gom install github.com/golangci/golangci-lint/v2/cmd/golangci-lint
gom list
gom sync
gom upgrade
```
