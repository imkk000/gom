package exec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/build"
	"os/exec"
	"path/filepath"
	"strings"
)

func IsLocal(pkg string) bool {
	return build.IsLocalImport(pkg) || filepath.IsAbs(pkg)
}

func Install(ctx context.Context, pkg string) error {
	target := pkg
	if !IsLocal(target) && !strings.Contains(target, "@") {
		target += "@latest"
	}
	out, err := exec.CommandContext(ctx, "go", "install", "-trimpath", "-ldflags=-s -w", target).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, bytes.TrimSpace(out))
	}
	return nil
}

func ImportPath(ctx context.Context, pkg string) (string, error) {
	out, err := exec.CommandContext(ctx, "go", "list", "-f", "{{.ImportPath}}", pkg).Output()
	if err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			return "", fmt.Errorf("%w: %s", err, bytes.TrimSpace(exitErr.Stderr))
		}
		return "", err
	}
	paths := strings.Fields(string(out))
	if len(paths) != 1 {
		return "", fmt.Errorf("%s matches %d packages, expected 1", pkg, len(paths))
	}
	return paths[0], nil
}
