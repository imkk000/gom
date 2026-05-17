package cmd

import (
	"context"
	"fmt"

	"github.com/imkk000/gom/internal/bin"
	"github.com/imkk000/gom/internal/store"
	"github.com/urfave/cli/v3"
)

func List() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list packages with install status",
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
				status := "missing"
				if _, ok := installed[p]; ok {
					status = "ok"
				}
				fmt.Printf("[%s] %s\n", status, p)
			}
			return nil
		},
	}
}
