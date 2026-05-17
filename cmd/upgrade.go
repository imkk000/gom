package cmd

import (
	"context"
	"fmt"

	"github.com/imkk000/gom/internal/exec"
	"github.com/imkk000/gom/internal/store"
	"github.com/urfave/cli/v3"
)

func Upgrade() *cli.Command {
	return &cli.Command{
		Name:  "upgrade",
		Usage: "reinstall all packages at @latest",
		Action: func(ctx context.Context, _ *cli.Command) error {
			pkgs, err := store.Load()
			if err != nil {
				return err
			}
			for _, p := range pkgs {
				if err := exec.InstallLatest(ctx, p); err != nil {
					fmt.Printf("[%s] failed: %v\n", p, err)
					continue
				}
				fmt.Printf("[%s] done\n", p)
			}
			return nil
		},
	}
}
