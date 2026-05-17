package exec

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

func Install(ctx context.Context, pkg string) error {
	target := pkg
	if !strings.Contains(target, "@") {
		target += "@latest"
	}
	cmd := exec.CommandContext(ctx, "go", "install", "-trimpath", "-ldflags=-s -w", target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InstallLatest(ctx context.Context, pkg string) error {
	base, _, _ := strings.Cut(pkg, "@")
	return Install(ctx, base+"@latest")
}

func Strip(pkg string) string {
	base, _, _ := strings.Cut(pkg, "@")
	return base
}
