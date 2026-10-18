package store_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/imkk000/gom/internal/store"
)

func TestSaveWritesThroughSymlink(t *testing.T) {
	mustNoErr := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}

	root := t.TempDir()
	dotfiles := filepath.Join(root, "dotfiles")
	mustNoErr(os.MkdirAll(dotfiles, 0o750))
	realFile := filepath.Join(dotfiles, "packages.json")
	mustNoErr(os.WriteFile(realFile, []byte(`["github.com/old/pkg"]`), 0o600))

	cfg := filepath.Join(root, "config")
	gomDir := filepath.Join(cfg, "gom")
	mustNoErr(os.MkdirAll(gomDir, 0o750))
	link := filepath.Join(gomDir, "packages.json")
	mustNoErr(os.Symlink(realFile, link))

	t.Setenv("XDG_CONFIG_HOME", cfg)
	mustNoErr(store.Add("github.com/new/pkg"))

	fi, err := os.Lstat(link)
	mustNoErr(err)
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatal("packages.json is no longer a symlink after Save")
	}

	pkgs, err := store.Load()
	mustNoErr(err)
	want := []string{"github.com/new/pkg", "github.com/old/pkg"}
	if !slices.Equal(pkgs, want) {
		t.Fatalf("got %v, want %v", pkgs, want)
	}

	info, err := os.Stat(realFile)
	mustNoErr(err)
	if info.Size() == 0 {
		t.Fatal("dotfiles target was not written through")
	}
}

func TestSplit(t *testing.T) {
	cases := []struct{ in, path, version string }{
		{"github.com/a/b", "github.com/a/b", ""},
		{"github.com/a/b@latest", "github.com/a/b", ""},
		{"github.com/a/b@v1.2.3", "github.com/a/b", "v1.2.3"},
	}
	for _, c := range cases {
		path, version := store.Split(c.in)
		if path != c.path || version != c.version {
			t.Errorf("Split(%q) = %q, %q; want %q, %q", c.in, path, version, c.path, c.version)
		}
	}
}

func TestAddReplacesPinAndRemoveDropsIt(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	for _, e := range []string{"github.com/a/b@v1.0.0", "github.com/c/d", "github.com/a/b@v2.0.0"} {
		if err := store.Add(e); err != nil {
			t.Fatal(err)
		}
	}
	pkgs, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"github.com/a/b@v2.0.0", "github.com/c/d"}; !slices.Equal(pkgs, want) {
		t.Fatalf("got %v, want %v", pkgs, want)
	}
	if err := store.Remove("github.com/a/b"); err != nil {
		t.Fatal(err)
	}
	pkgs, err = store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"github.com/c/d"}; !slices.Equal(pkgs, want) {
		t.Fatalf("got %v, want %v", pkgs, want)
	}
}
