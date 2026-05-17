package main

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/santhoshjohn/awsmgr/internal/aws"
	"github.com/santhoshjohn/awsmgr/internal/doctor"
)

func TestRunPrintsUsageWhenNoArgs(t *testing.T) {
	stdoutBuf, restore := stubCLI(t)
	defer restore()

	if err := run(nil); err != nil {
		t.Fatalf("run(nil) error = %v, want nil", err)
	}

	got := stdoutBuf.String()
	if !strings.Contains(got, "Usage:") || !strings.Contains(got, "awsmgr doctor") {
		t.Fatalf("run(nil) output = %q, want usage text", got)
	}
}

func TestRunPrintsUsageForHelp(t *testing.T) {
	stdoutBuf, restore := stubCLI(t)
	defer restore()

	if err := run([]string{"--help"}); err != nil {
		t.Fatalf("run(--help) error = %v, want nil", err)
	}

	if got := stdoutBuf.String(); !strings.Contains(got, "awsmgr current") {
		t.Fatalf("run(--help) output = %q, want usage text", got)
	}
}

func TestRunReturnsUnknownCommandError(t *testing.T) {
	_, restore := stubCLI(t)
	defer restore()

	err := run([]string{"bogus"})
	if err == nil || err.Error() != `unknown command "bogus"` {
		t.Fatalf("run(bogus) error = %v, want unknown command error", err)
	}
}

func TestRunListPrintsProfiles(t *testing.T) {
	stdoutBuf, restore := stubCLI(t)
	defer restore()

	loadProfiles = func(string) ([]string, error) {
		return []string{"default", "prod"}, nil
	}

	if err := run([]string{"list"}); err != nil {
		t.Fatalf("run(list) error = %v, want nil", err)
	}

	if got := stdoutBuf.String(); got != "default\nprod\n" {
		t.Fatalf("run(list) output = %q, want profile list", got)
	}
}

func TestRunCurrentPrintsActiveProfile(t *testing.T) {
	stdoutBuf, restore := stubCLI(t)
	defer restore()

	currentProfile = func() string { return "prod" }

	if err := run([]string{"current"}); err != nil {
		t.Fatalf("run(current) error = %v, want nil", err)
	}

	if got := stdoutBuf.String(); got != "prod\n" {
		t.Fatalf("run(current) output = %q, want active profile", got)
	}
}

func TestRunValidateReturnsErrorForNoProfiles(t *testing.T) {
	_, restore := stubCLI(t)
	defer restore()

	loadProfiles = func(string) ([]string, error) {
		return nil, nil
	}

	err := run([]string{"validate"})
	if err == nil || err.Error() != "no AWS profiles found" {
		t.Fatalf("run(validate) error = %v, want no profiles error", err)
	}
}

func TestRunValidatePrintsProfileStatuses(t *testing.T) {
	_, restore := stubCLI(t)
	defer restore()

	loadProfiles = func(string) ([]string, error) {
		return []string{"default", "prod"}, nil
	}
	validateProfiles = func(_ context.Context, profiles []string, timeout time.Duration) []aws.ValidationResult {
		if !reflect.DeepEqual(profiles, []string{"default", "prod"}) {
			t.Fatalf("validate profiles = %v, want [default prod]", profiles)
		}
		if timeout != aws.DefaultProfileValidationTimeout {
			t.Fatalf("validate timeout = %v, want %v", timeout, aws.DefaultProfileValidationTimeout)
		}
		return []aws.ValidationResult{
			{Profile: "default", Valid: true, AccountID: "123456789012"},
			{Profile: "prod", Valid: false, Err: errors.New("expired token")},
		}
	}

	var got []string
	printStatus = func(status, area, message string) {
		got = append(got, status+"|"+area+"|"+message)
	}

	if err := run([]string{"validate"}); err != nil {
		t.Fatalf("run(validate) error = %v, want nil", err)
	}

	want := []string{
		"OK|AWS|profile=default account=123456789012",
		"FAIL|AWS|profile=prod invalid",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("printed statuses = %v, want %v", got, want)
	}
}

func TestRunDoctorPrintsReportAndReturnsNilOnSuccess(t *testing.T) {
	_, restore := stubCLI(t)
	defer restore()

	currentProfile = func() string { return "prod" }
	wantReport := doctor.Report{Checks: []doctor.Check{{Name: "AWS", Status: doctor.StatusOK, OK: true, Message: "profile=prod account=123456789012"}}}
	runDoctorChecks = func(_ context.Context, profile string) doctor.Report {
		if profile != "prod" {
			t.Fatalf("doctor profile = %q, want prod", profile)
		}
		return wantReport
	}

	var printed doctor.Report
	printDoctorReport = func(report doctor.Report) {
		printed = report
	}

	if err := run([]string{"doctor"}); err != nil {
		t.Fatalf("run(doctor) error = %v, want nil", err)
	}
	if len(printed.Checks) != 1 {
		t.Fatalf("len(printed.Checks) = %d, want 1", len(printed.Checks))
	}
	if got := printed.Checks[0]; got.Name != wantReport.Checks[0].Name || got.Status != wantReport.Checks[0].Status || got.OK != wantReport.Checks[0].OK || got.Message != wantReport.Checks[0].Message || got.Err != nil {
		t.Fatalf("printed doctor check = %#v, want %#v", got, wantReport.Checks[0])
	}
}

func TestRunDoctorDoesNotFailOnWarnings(t *testing.T) {
	_, restore := stubCLI(t)
	defer restore()

	currentProfile = func() string { return "prod" }
	runDoctorChecks = func(context.Context, string) doctor.Report {
		return doctor.Report{Checks: []doctor.Check{{Name: "Alignment", Status: doctor.StatusWarn, OK: true, Message: "context=dev-cluster does not appear to match profile=prod or account=123456789012"}}}
	}

	if err := run([]string{"doctor"}); err != nil {
		t.Fatalf("run(doctor) error = %v, want nil for warning-only report", err)
	}
}

func TestRunDoctorReturnsErrorOnFailedDiagnostics(t *testing.T) {
	_, restore := stubCLI(t)
	defer restore()

	currentProfile = func() string { return "prod" }
	runDoctorChecks = func(context.Context, string) doctor.Report {
		return doctor.Report{Checks: []doctor.Check{{Name: "AWS", Message: "profile=prod invalid"}}}
	}

	err := run([]string{"doctor"})
	if err == nil || err.Error() != "environment diagnostics failed" {
		t.Fatalf("run(doctor) error = %v, want diagnostics failure", err)
	}
}

func stubCLI(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()

	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}
	prevStdout := stdout
	prevStderr := stderr
	prevLoadProfiles := loadProfiles
	prevCurrentProfile := currentProfile
	prevValidateProfiles := validateProfiles
	prevRunDoctorChecks := runDoctorChecks
	prevPrintStatus := printStatus
	prevPrintDoctorReport := printDoctorReport

	stdout = stdoutBuf
	stderr = stderrBuf
	loadProfiles = aws.LoadProfiles
	currentProfile = aws.CurrentProfile
	validateProfiles = aws.ValidateProfiles
	runDoctorChecks = doctor.Run
	printStatus = func(status, area, message string) {}
	printDoctorReport = func(report doctor.Report) {}

	return stdoutBuf, func() {
		stdout = prevStdout
		stderr = prevStderr
		loadProfiles = prevLoadProfiles
		currentProfile = prevCurrentProfile
		validateProfiles = prevValidateProfiles
		runDoctorChecks = prevRunDoctorChecks
		printStatus = prevPrintStatus
		printDoctorReport = prevPrintDoctorReport
	}
}
