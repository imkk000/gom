package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/imkk000/gom/internal/bin"
	"github.com/imkk000/gom/internal/exec"
	"github.com/imkk000/gom/internal/store"
	"github.com/urfave/cli/v3"
)

func Uninstall() *cli.Command {
	return &cli.Command{
		Name:      "uninstall",
		Usage:     "remove a package binary and record",
		ArgsUsage: "<pkg>",
		Action: func(_ context.Context, cmd *cli.Command) error {
			pkg := cmd.Args().First()
			if pkg == "" {
				return fmt.Errorf("package argument required")
			}
			base := exec.Strip(pkg)
			path, err := bin.Lookup(base)
			if err != nil {
				fmt.Printf("[%s] failed: %v\n", base, err)
				return err
			}
			if path != "" {
				if err := os.Remove(path); err != nil {
					fmt.Printf("[%s] failed: %v\n", base, err)
					return err
				}
			}
			if err := store.Remove(base); err != nil {
				fmt.Printf("[%s] failed: %v\n", base, err)
				return err
			}
			fmt.Printf("[%s] done\n", base)
			return nil
		},
	}
}
