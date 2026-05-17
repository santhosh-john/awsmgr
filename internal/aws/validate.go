package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"
)

const DefaultProfileValidationTimeout = 12 * time.Second

// MaxConcurrentWorkers prevents OS process starvation when thousands of profiles exist
const MaxConcurrentWorkers = 10

var execCommandContext = exec.CommandContext

type ValidationResult struct {
	Profile   string
	Valid     bool
	AccountID string
	Message   string
	Err       error
}

type callerIdentity struct {
	Account string `json:"Account"`
}

func ValidateProfiles(ctx context.Context, profiles []string, timeout time.Duration) []ValidationResult {
	results := make(chan ValidationResult, len(profiles)) // Buffered to avoid goroutine blockages
	var wg sync.WaitGroup

	// Semaphore channel to limit max parallel processes
	sem := make(chan struct{}, MaxConcurrentWorkers)

	for _, profile := range profiles {
		wg.Add(1)
		sem <- struct{}{} // Acquire token

		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }() // Release token

			res := ValidateProfile(ctx, p, timeout)
			results <- res // FIXED: Send to channel using <- operator instead of bracket indexing
		}(profile)
	}

	// Wait asynchronously for workers to wrap up
	wg.Wait()
	close(results)

	validationResults := make([]ValidationResult, 0, len(profiles))
	for result := range results {
		validationResults = append(validationResults, result)
	}

	sort.Slice(validationResults, func(i, j int) bool {
		return validationResults[i].Profile < validationResults[j].Profile
	})
	return validationResults
}

func ValidateProfile(ctx context.Context, profile string, timeout time.Duration) ValidationResult {
	profileCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := execCommandContext(profileCtx, "aws", "sts", "get-caller-identity", "--profile", profile, "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ValidationResult{Profile: profile, Message: validationFailureMessage(profile, err, output), Err: err}
	}

	var identity callerIdentity
	if err := json.Unmarshal(output, &identity); err != nil {
		return ValidationResult{Profile: profile, Message: fmt.Sprintf("profile=%s returned invalid AWS CLI response", profile), Err: err}
	}
	if identity.Account == "" {
		return ValidationResult{Profile: profile, Message: fmt.Sprintf("profile=%s AWS CLI response missing account ID", profile), Err: errMissingAccountID}
	}

	return ValidationResult{
		Profile:   profile,
		Valid:     true,
		AccountID: identity.Account,
	}
}

func validationFailureMessage(profile string, err error, output []byte) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Sprintf("profile=%s validation timed out", profile)
	}

	var execErr *exec.Error
	if errors.As(err, &execErr) && errors.Is(execErr.Err, exec.ErrNotFound) {
		return "aws CLI not found"
	}

	message := summarizeCommandOutput(string(output))
	if message == "" {
		return fmt.Sprintf("profile=%s validation failed", profile)
	}

	lower := strings.ToLower(message)
	switch {
	case strings.Contains(lower, "sso session associated with this profile has expired"),
		strings.Contains(lower, "error loading sso token"),
		strings.Contains(lower, "token has expired"):
		return fmt.Sprintf("profile=%s AWS SSO session expired or invalid", profile)
	case strings.Contains(lower, "could not be found") && strings.Contains(lower, "profile"):
		return fmt.Sprintf("profile=%s not found in AWS CLI config", profile)
	case strings.Contains(lower, "unable to locate credentials"):
		return fmt.Sprintf("profile=%s credentials unavailable", profile)
	default:
		return fmt.Sprintf("profile=%s AWS CLI error: %s", profile, message)
	}
}

func summarizeCommandOutput(output string) string {
	lines := strings.Split(output, "\n")
	var parts []string
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}
