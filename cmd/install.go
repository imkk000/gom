package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/imkk000/gom/internal/bin"
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
				return errors.New("package argument required")
			}
			path, _ := store.Split(pkg)
			if err := exec.Install(ctx, pkg); err != nil {
				fmt.Printf("[%s] failed: %v\n", path, err)
				return err
			}
			b, err := bin.Lookup(path)
			if err != nil {
				fmt.Printf("[%s] failed: %v\n", path, err)
				return err
			}
			entry := path
			if b.Pinnable() {
				entry = store.Join(path, b.Version)
			}
			if err := store.Add(entry); err != nil {
				fmt.Printf("[%s] failed: %v\n", path, err)
				return err
			}
			fmt.Printf("[%s] done\n", entry)
			return nil
		},
	}
}
