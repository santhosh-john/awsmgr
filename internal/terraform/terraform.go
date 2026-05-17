package terraform

import (
	"context"
	"os/exec"
	"time"
)

type AvailabilityResult struct {
	OK      bool
	Message string
	Err     error
}

func Available(ctx context.Context, timeout time.Duration) AvailabilityResult {
	path, err := exec.LookPath("terraform")
	if err != nil {
		return AvailabilityResult{Message: "terraform not found", Err: err}
	}

	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := exec.CommandContext(cmdCtx, path, "version", "-json").Run(); err != nil {
		return AvailabilityResult{Message: "terraform unavailable", Err: err}
	}

	return AvailabilityResult{OK: true}
}
