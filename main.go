// Package main provides the CLI entrypoint for awsmgr.
//
// awsmgr is an operational workflow and environment diagnostics CLI for engineers
// working with AWS SSO, multiple AWS profiles, Kubernetes/EKS, Terraform, and
// enterprise network environments.
//
// Author: Santhosh John
// License: MIT
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/santhosh-john/awsmgr/internal/aws"
	"github.com/santhosh-john/awsmgr/internal/doctor"
	"github.com/santhosh-john/awsmgr/internal/ui"
)

const commandTimeout = 30 * time.Second

var (
	stdout            io.Writer = os.Stdout
	stderr            io.Writer = os.Stderr
	loadProfiles                = aws.LoadProfiles
	currentProfile              = aws.CurrentProfile
	validateProfiles            = aws.ValidateProfiles
	runDoctorChecks             = doctor.Run
	printStatus                 = ui.PrintStatus
	printDoctorReport           = ui.PrintDoctorReport
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(stderr, "awsmgr: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	switch args[0] {
	case "list":
		return runList()
	case "current":
		return runCurrent()
	case "validate":
		return runValidate(ctx)
	case "doctor":
		return runDoctor(ctx)
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runList() error {
	profiles, err := loadProfiles("")
	if err != nil {
		return err
	}
	for _, profile := range profiles {
		fmt.Fprintln(stdout, profile)
	}
	return nil
}

func runCurrent() error {
	fmt.Fprintln(stdout, currentProfile())
	return nil
}

func runValidate(ctx context.Context) error {
	profiles, err := loadProfiles("")
	if err != nil {
		return err
	}
	if len(profiles) == 0 {
		return fmt.Errorf("no AWS profiles found")
	}

	fmt.Fprintln(stdout, "Concurrently validating access sessions...")
	results := validateProfiles(ctx, profiles, aws.DefaultProfileValidationTimeout)
	for _, result := range results {
		if result.Valid {
			printStatus("OK", "AWS", fmt.Sprintf("profile=%s account=%s", result.Profile, result.AccountID))
			continue
		}
		printStatus("FAIL", "AWS", fmt.Sprintf("profile=%s invalid", result.Profile))
	}
	fmt.Fprintln(stdout)
	return nil
}

func runDoctor(ctx context.Context) error {
	report := runDoctorChecks(ctx, currentProfile())
	printDoctorReport(report)

	if !report.OK() {
		return fmt.Errorf("environment diagnostics failed")
	}
	return nil
}

func printUsage() {
	fmt.Fprintf(stdout, `%sawsmgr%s - Operational workflow and environment diagnostics.

Usage:
  awsmgr list
  awsmgr current
  awsmgr validate
  awsmgr doctor
`, ui.ColorCyan+ui.ColorBold, ui.ColorReset)
}
