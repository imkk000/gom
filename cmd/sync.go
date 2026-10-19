package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/imkk000/gom/internal/bin"
	"github.com/imkk000/gom/internal/store"
	"github.com/urfave/cli/v3"
)

func Sync() *cli.Command {
	return &cli.Command{
		Name:  "sync",
		Usage: "mirror packages.json: install missing and remove extras",
		Flags: []cli.Flag{jobsFlag()},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			pkgs, err := store.Load()
			if err != nil {
				return err
			}
			installed, err := bin.Installed()
			if err != nil {
				return err
			}
			wanted := make(map[string]struct{}, len(pkgs))
			var pending []string
			for _, p := range pkgs {
				path, version := store.Split(p)
				wanted[path] = struct{}{}
				if b, ok := installed[path]; ok && (version == "" || version == b.Version) {
					continue
				}
				pending = append(pending, p)
			}
			parallel(ctx, cmd.Int("jobs"), pending, func(ctx context.Context, p string) {
				path, _ := store.Split(p)
				if err := reinstall(ctx, p, installed[path].File); err != nil {
					fmt.Printf("[%s] failed: %v\n", p, err)
					return
				}
				fmt.Printf("[%s] installed\n", p)
			})
			for path, b := range installed {
				if _, ok := wanted[path]; ok {
					continue
				}
				if err := os.Remove(b.File); err != nil {
					fmt.Printf("[%s] failed: %v\n", path, err)
					continue
				}
				fmt.Printf("[%s] removed\n", path)
			}
			return nil
		},
	}
}
