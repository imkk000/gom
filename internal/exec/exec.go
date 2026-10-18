package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func Install(ctx context.Context, pkg string) error {
	target := pkg
	if !strings.Contains(target, "@") {
		target += "@latest"
	}
	out, err := exec.CommandContext(ctx, "go", "install", "-trimpath", "-ldflags=-s -w", target).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, bytes.TrimSpace(out))
	}
	return nil
}
