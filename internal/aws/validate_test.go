package aws

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestValidateProfileSuccess(t *testing.T) {
	binDir := t.TempDir()
	writeExecutable(t, filepath.Join(binDir, executableName("aws")), `printf '{"Account":"123456789012"}\n'
`, `@echo off
echo {"Account":"123456789012"}
`)
	t.Setenv("PATH", binDir+string(filepath.ListSeparator)+os.Getenv("PATH"))

	result := ValidateProfile(context.Background(), "prod", time.Second)
	if !result.Valid {
		t.Fatalf("ValidateProfile().Valid = false, want true; result=%#v", result)
	}
	if result.AccountID != "123456789012" {
		t.Fatalf("ValidateProfile().AccountID = %q, want 123456789012", result.AccountID)
	}
	if result.Message != "" {
		t.Fatalf("ValidateProfile().Message = %q, want empty", result.Message)
	}
	if result.Err != nil {
		t.Fatalf("ValidateProfile().Err = %v, want nil", result.Err)
	}
}

func TestValidateProfileReportsMissingAWSCLI(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	result := ValidateProfile(context.Background(), "prod", time.Second)
	if result.Valid {
		t.Fatal("ValidateProfile().Valid = true, want false")
	}
	if result.Message != "aws CLI not found" {
		t.Fatalf("ValidateProfile().Message = %q, want aws CLI not found", result.Message)
	}
	if result.Err == nil {
		t.Fatal("ValidateProfile().Err = nil, want error")
	}
}

func TestValidateProfileReportsExpiredSSO(t *testing.T) {
	binDir := t.TempDir()
	writeExecutable(t, filepath.Join(binDir, executableName("aws")), `echo 'The SSO session associated with this profile has expired or is otherwise invalid' >&2
exit 255
`, `@echo off
echo The SSO session associated with this profile has expired or is otherwise invalid 1>&2
exit /b 255
`)
	t.Setenv("PATH", binDir+string(filepath.ListSeparator)+os.Getenv("PATH"))

	result := ValidateProfile(context.Background(), "prod", time.Second)
	if result.Valid {
		t.Fatal("ValidateProfile().Valid = true, want false")
	}
	if result.Message != "profile=prod AWS SSO session expired or invalid" {
		t.Fatalf("ValidateProfile().Message = %q, want expired SSO message", result.Message)
	}
	if result.Err == nil {
		t.Fatal("ValidateProfile().Err = nil, want error")
	}
}

func TestValidateProfileReportsMissingProfile(t *testing.T) {
	binDir := t.TempDir()
	writeExecutable(t, filepath.Join(binDir, executableName("aws")), `echo 'The config profile (prod) could not be found' >&2
exit 255
`, `@echo off
echo The config profile (prod) could not be found 1>&2
exit /b 255
`)
	t.Setenv("PATH", binDir+string(filepath.ListSeparator)+os.Getenv("PATH"))

	result := ValidateProfile(context.Background(), "prod", time.Second)
	if result.Valid {
		t.Fatal("ValidateProfile().Valid = true, want false")
	}
	if result.Message != "profile=prod not found in AWS CLI config" {
		t.Fatalf("ValidateProfile().Message = %q, want missing profile message", result.Message)
	}
	if result.Err == nil {
		t.Fatal("ValidateProfile().Err = nil, want error")
	}
}

func writeExecutable(t *testing.T, path, unixBody, windowsBody string) {
	t.Helper()

	body := "#!/bin/sh\n" + unixBody
	if runtime.GOOS == "windows" {
		body = windowsBody
	}
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
}

func executableName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".bat"
	}
	return name
}
