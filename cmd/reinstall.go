package cmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/imkk000/gom/internal/exec"
)

func reinstall(ctx context.Context, pkg, file string) error {
	if file == "" {
		return exec.Install(ctx, pkg)
	}
	backup := filepath.Join(filepath.Dir(file), "."+filepath.Base(file)+".old")
	if err := os.Rename(file, backup); err != nil {
		return err
	}
	if err := exec.Install(ctx, pkg); err != nil {
		return errors.Join(err, os.Rename(backup, file))
	}
	return os.Remove(backup)
}
