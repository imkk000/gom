package cmd

import (
	"context"
	"fmt"
	"slices"

	"github.com/imkk000/gom/internal/bin"
	"github.com/imkk000/gom/internal/store"
	"github.com/urfave/cli/v3"
)

func Export() *cli.Command {
	return &cli.Command{
		Name:  "export",
		Usage: "scan GOPATH/bin and write all Go packages with their installed versions into packages.json",
		Action: func(_ context.Context, _ *cli.Command) error {
			current, err := store.Load()
			if err != nil {
				return err
			}
			installed, err := bin.Installed()
			if err != nil {
				return err
			}
			next := make([]string, 0, len(installed))
			for path, b := range installed {
				var version string
				if b.Pinnable() {
					version = b.Version
				}
				next = append(next, store.Join(path, version))
			}
			slices.Sort(next)

			added, removed := diff(current, next)
			for _, p := range added {
				fmt.Printf("[+] %s\n", p)
			}
			for _, p := range removed {
				fmt.Printf("[-] %s\n", p)
			}
			if len(added) == 0 && len(removed) == 0 {
				fmt.Println("no changes")
			}
			return store.Save(next)
		},
	}
}

func diff(old, next []string) (added, removed []string) {
	for _, p := range next {
		if !slices.Contains(old, p) {
			added = append(added, p)
		}
	}
	for _, p := range old {
		if !slices.Contains(next, p) {
			removed = append(removed, p)
		}
	}
	return added, removed
}
