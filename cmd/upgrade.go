package cmd

import (
	"context"
	"fmt"
	pathpkg "path"
	"path/filepath"
	"slices"

	"github.com/imkk000/gom/internal/bin"
	"github.com/imkk000/gom/internal/store"
	"github.com/urfave/cli/v3"
)

func Upgrade() *cli.Command {
	return &cli.Command{
		Name:      "upgrade",
		Usage:     "reinstall packages at @latest without updating packages.json (all if none given)",
		ArgsUsage: "[name|pkg...]",
		Flags:     []cli.Flag{jobsFlag()},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			pkgs, err := store.Load()
			if err != nil {
				return err
			}
			installed, err := bin.Installed()
			if err != nil {
				return err
			}
			paths := make([]string, 0, len(pkgs))
			for _, p := range pkgs {
				path, _ := store.Split(p)
				paths = append(paths, path)
			}
			pending, err := selectPackages(paths, installed, cmd.Args().Slice())
			if err != nil {
				return err
			}
			parallel(ctx, cmd.Int("jobs"), pending, func(ctx context.Context, p string) {
				if err := reinstall(ctx, p, installed[p].File); err != nil {
					fmt.Printf("[%s] failed: %v\n", p, err)
					return
				}
				fmt.Printf("[%s] done\n", p)
			})
			return nil
		},
	}
}

func selectPackages(paths []string, installed map[string]bin.Binary, names []string) ([]string, error) {
	if len(names) == 0 {
		return paths, nil
	}
	var selected []string
	for _, name := range names {
		var found bool
		for _, p := range paths {
			if p == name || binaryName(p, installed[p]) == name {
				found = true
				if !slices.Contains(selected, p) {
					selected = append(selected, p)
				}
			}
		}
		if !found {
			return nil, fmt.Errorf("package not found in packages.json: %s", name)
		}
	}
	return selected, nil
}

func binaryName(path string, b bin.Binary) string {
	if b.File != "" {
		return filepath.Base(b.File)
	}
	return pathpkg.Base(path)
}
