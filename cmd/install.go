package cmd

import (
	"context"
	"fmt"

	"github.com/imkk000/gom/internal/exec"
	"github.com/imkk000/gom/internal/store"
	"github.com/urfave/cli/v3"
)

func Install() *cli.Command {
	return &cli.Command{
		Name:      "install",
		Usage:     "install a package and record it",
		ArgsUsage: "<pkg>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			pkg := cmd.Args().First()
			if pkg == "" {
				return fmt.Errorf("package argument required")
			}
			base := exec.Strip(pkg)
			if err := exec.Install(ctx, pkg); err != nil {
				fmt.Printf("[%s] failed: %v\n", base, err)
				return err
			}
			if err := store.Add(base); err != nil {
				fmt.Printf("[%s] failed: %v\n", base, err)
				return err
			}
			fmt.Printf("[%s] done\n", base)
			return nil
		},
	}
}
