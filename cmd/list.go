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
		Action: func(_ context.Context, _ *cli.Command) error {
			pkgs, err := store.Load()
			if err != nil {
				return err
			}
			installed, err := bin.Installed()
			if err != nil {
				return err
			}
			for _, p := range pkgs {
				path, version := store.Split(p)
				b, ok := installed[path]
				switch {
				case !ok:
					fmt.Printf("[missing] %s\n", p)
				case version != "" && version != b.Version:
					fmt.Printf("[mismatch] %s (installed %s)\n", p, b.Version)
				default:
					fmt.Printf("[ok] %s\n", p)
				}
			}
			return nil
		},
	}
}
