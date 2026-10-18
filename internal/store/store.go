package store

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func Path() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "gom", "packages.json"), nil
}

func Load() ([]string, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []string{}, nil
	}
	var pkgs []string
	if err := json.Unmarshal(data, &pkgs); err != nil {
		return nil, err
	}
	return pkgs, nil
}

func resolveLink(path string) string {
	for range 32 {
		fi, err := os.Lstat(path)
		if err != nil || fi.Mode()&os.ModeSymlink == 0 {
			return path
		}
		dest, err := os.Readlink(path)
		if err != nil {
			return path
		}
		if !filepath.IsAbs(dest) {
			dest = filepath.Join(filepath.Dir(path), dest)
		}
		path = dest
	}
	return path
}

func Save(pkgs []string) error {
	path, err := Path()
	if err != nil {
		return err
	}
	target := resolveLink(path)
	dir := filepath.Dir(target)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	out := slices.Clone(pkgs)
	slices.Sort(out)
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "packages-*.json")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		return errors.Join(err, tmp.Close(), os.Remove(tmpName))
	}
	if err := tmp.Close(); err != nil {
		return errors.Join(err, os.Remove(tmpName))
	}
	if err := os.Chmod(tmpName, 0o600); err != nil {
		return errors.Join(err, os.Remove(tmpName))
	}
	if err := os.Rename(tmpName, target); err != nil {
		return errors.Join(err, os.Remove(tmpName))
	}
	return nil
}

func Split(entry string) (path, version string) {
	path, version, _ = strings.Cut(entry, "@")
	if version == "latest" {
		version = ""
	}
	return path, version
}

func Join(path, version string) string {
	if version == "" {
		return path
	}
	return path + "@" + version
}

func samePath(path string) func(string) bool {
	return func(entry string) bool {
		p, _ := Split(entry)
		return p == path
	}
}

func Add(entry string) error {
	pkgs, err := Load()
	if err != nil {
		return err
	}
	if slices.Contains(pkgs, entry) {
		return nil
	}
	path, _ := Split(entry)
	pkgs = slices.DeleteFunc(pkgs, samePath(path))
	pkgs = append(pkgs, entry)
	return Save(pkgs)
}

func Remove(path string) error {
	pkgs, err := Load()
	if err != nil {
		return err
	}
	return Save(slices.DeleteFunc(pkgs, samePath(path)))
}
