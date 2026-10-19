package cmd

import (
	"context"
	"fmt"

	"github.com/imkk000/gom/internal/bin"
	"github.com/imkk000/gom/internal/store"
	"github.com/urfave/cli/v3"
)

func Upgrade() *cli.Command {
	return &cli.Command{
		Name:  "upgrade",
		Usage: "reinstall all packages at @latest without updating packages.json",
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
			pending := make([]string, 0, len(pkgs))
			for _, p := range pkgs {
				path, _ := store.Split(p)
				pending = append(pending, path)
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
