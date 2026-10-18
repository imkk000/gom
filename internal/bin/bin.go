package bin

import (
	"debug/buildinfo"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GopathBin() (string, error) {
	if p := os.Getenv("GOBIN"); p != "" {
		return p, nil
	}
	if out, err := exec.Command("go", "env", "GOBIN").Output(); err == nil {
		if p := strings.TrimSpace(string(out)); p != "" {
			return p, nil
		}
	}
	if p := os.Getenv("GOPATH"); p != "" {
		first := strings.SplitN(p, string(os.PathListSeparator), 2)[0]
		return filepath.Join(first, "bin"), nil
	}
	if out, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		if p := strings.TrimSpace(string(out)); p != "" {
			first := strings.SplitN(p, string(os.PathListSeparator), 2)[0]
			return filepath.Join(first, "bin"), nil
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "go", "bin"), nil
}

type Binary struct {
	File    string
	Version string
}

func (b Binary) Pinnable() bool {
	return b.Version != "" && b.Version != "(devel)" && !strings.Contains(b.Version, "+dirty")
}

func Installed() (map[string]Binary, error) {
	dir, err := GopathBin()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return make(map[string]Binary), nil
		}
		return nil, err
	}
	result := make(map[string]Binary)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		full := filepath.Join(dir, e.Name())
		info, err := buildinfo.ReadFile(full)
		if err != nil {
			continue
		}
		path := strings.TrimSuffix(info.Path, "/...")
		if path == "" || path == "command-line-arguments" {
			path = info.Main.Path
		}
		if path == "" || path == "command-line-arguments" {
			continue
		}
		result[path] = Binary{File: full, Version: info.Main.Version}
	}
	return result, nil
}

func Lookup(pkg string) (Binary, error) {
	base, _, _ := strings.Cut(strings.TrimSuffix(pkg, "/..."), "@")
	installed, err := Installed()
	if err != nil {
		return Binary{}, err
	}
	return installed[base], nil
}
