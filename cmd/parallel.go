package cmd

import (
	"context"
	"errors"
	"sync"

	"github.com/urfave/cli/v3"
)

const defaultJobs = 4

func jobsFlag() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    "jobs",
		Aliases: []string{"j"},
		Value:   defaultJobs,
		Usage:   "maximum concurrent installs",
		Validator: func(n int) error {
			if n < 1 {
				return errors.New("jobs must be at least 1")
			}
			return nil
		},
	}
}

func parallel[T any](ctx context.Context, limit int, items []T, fn func(context.Context, T)) {
	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup
	defer wg.Wait()
	for _, item := range items {
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			return
		}
		wg.Go(func() {
			defer func() { <-sem }()
			fn(ctx, item)
		})
	}
}
