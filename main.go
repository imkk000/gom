package main

import (
	"context"
	"fmt"
	"os"

	"github.com/imkk000/gom/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:                  "gom",
		Usage:                 "Go package manager",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			cmd.Install(),
			cmd.Uninstall(),
			cmd.List(),
			cmd.Sync(),
			cmd.Upgrade(),
			cmd.Export(),
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
