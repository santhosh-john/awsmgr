package kube

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

type ContextResult struct {
	OK      bool
	Context string
	Message string
	Err     error
}

func CurrentContext(ctx context.Context, timeout time.Duration) ContextResult {
	if _, err := exec.LookPath("kubectl"); err != nil {
		return ContextResult{Message: "kubectl not found", Err: err}
	}

	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	output, err := exec.CommandContext(cmdCtx, "kubectl", "config", "current-context").Output()
	if err != nil {
		return ContextResult{Message: "kubectl context unavailable", Err: err}
	}

	currentContext := strings.TrimSpace(string(output))
	if currentContext == "" {
		return ContextResult{Message: "kubectl context unavailable", Err: errors.New("empty kubectl context")}
	}

	return ContextResult{OK: true, Context: currentContext}
}
