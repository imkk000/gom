# gom

Go package manager. Track Go-installed tools in a single packages file and keep them in sync across machines.

## Install

```sh
go install github.com/imkk000/gom@latest
```

## Usage

```sh
gom install <pkg>[@version]  Install a package and record it with its installed version
gom install .                Install the package in the current directory (any local path works)
gom uninstall <pkg>          Remove a package binary and record
gom list                     List packages with install status
gom sync [-j N]              Mirror packages.json: install missing or mismatched versions, remove extras
gom upgrade [-j N]           Reinstall all packages at @latest (packages.json is not updated)
gom export                   Scan GOPATH/bin and write all Go packages with versions into packages.json (with diff)
```

`sync` and `upgrade` run up to `-j/--jobs` installs concurrently (default 4).

`packages.json` is the source of truth for `sync`. After `upgrade`, run `export` to record the new versions; running `sync` first reinstalls the recorded versions, which is a downgrade. `list` shows such packages as `[mismatch]`.

Binaries built from a local checkout (`(devel)` or `+dirty`) are recorded without a version and track `@latest`.

Every install runs `go install -trimpath -ldflags="-s -w"`.

`gom install` always writes a fresh binary. Plain `go install` skips writing when the existing binary is up to date, so `gom` moves the existing binary aside first, then deletes it on success or restores it if the install fails.

Local paths (`.`, `./cmd/foo`, absolute paths) are resolved to their import path with `go list` and recorded under that path.

## Packages file

Location:

- `$XDG_CONFIG_HOME/gom/packages.json`
- fallback `~/.config/gom/packages.json`

Format:

```json
["github.com/foo/bar", "github.com/baz/qux@v1.2.3"]
```

Entries without a version track `@latest`.

## Example

```sh
gom install github.com/golangci/golangci-lint/v2/cmd/golangci-lint
gom list
gom sync
gom upgrade
```
