package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/imkk000/gom/internal/bin"
	"github.com/imkk000/gom/internal/exec"
	"github.com/imkk000/gom/internal/store"
	"github.com/urfave/cli/v3"
)

func Sync() *cli.Command {
	return &cli.Command{
		Name:  "sync",
		Usage: "mirror packages.json: install missing and remove extras",
		Action: func(ctx context.Context, _ *cli.Command) error {
			pkgs, err := store.Load()
			if err != nil {
				return err
			}
			installed, err := bin.Installed()
			if err != nil {
				return err
			}
			for _, p := range pkgs {
				if _, ok := installed[p]; ok {
					continue
				}
				if err := exec.InstallLatest(ctx, p); err != nil {
					fmt.Printf("[%s] failed: %v\n", p, err)
					continue
				}
				fmt.Printf("[%s] installed\n", p)
			}
			for pkg, path := range installed {
				if slices.Contains(pkgs, pkg) {
					continue
				}
				if err := os.Remove(path); err != nil {
					fmt.Printf("[%s] failed: %v\n", pkg, err)
					continue
				}
				fmt.Printf("[%s] removed\n", pkg)
			}
			return nil
		},
	}
}
