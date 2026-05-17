package terraform

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestAvailableReportsMissingTerraform(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	result := Available(context.Background(), time.Second)
	if result.OK {
		t.Fatal("Available().OK = true, want false")
	}
	if result.Message != "terraform not found" {
		t.Fatalf("Available().Message = %q, want terraform not found", result.Message)
	}
}

func TestAvailableReturnsOKWhenTerraformResponds(t *testing.T) {
	binDir := t.TempDir()
	writeExecutable(t, filepath.Join(binDir, executableName("terraform")), `if [ "$1" = "version" ] && [ "$2" = "-json" ]; then
  printf '{"terraform_version":"1.9.0"}\n'
  exit 0
fi
exit 1
`, `@echo off
if "%1"=="version" if "%2"=="-json" (
  echo {"terraform_version":"1.9.0"}
  exit /b 0
)
exit /b 1
`)
	t.Setenv("PATH", binDir+string(filepath.ListSeparator)+os.Getenv("PATH"))

	result := Available(context.Background(), time.Second)
	if !result.OK {
		t.Fatalf("Available().OK = false, want true; result=%#v", result)
	}
	if result.Message != "" {
		t.Fatalf("Available().Message = %q, want empty", result.Message)
	}
}

func TestAvailableReportsFailureWhenTerraformErrors(t *testing.T) {
	binDir := t.TempDir()
	writeExecutable(t, filepath.Join(binDir, executableName("terraform")), "exit 1\n", "@echo off\r\nexit /b 1\r\n")
	t.Setenv("PATH", binDir+string(filepath.ListSeparator)+os.Getenv("PATH"))

	result := Available(context.Background(), time.Second)
	if result.OK {
		t.Fatal("Available().OK = true, want false")
	}
	if result.Message != "terraform unavailable" {
		t.Fatalf("Available().Message = %q, want terraform unavailable", result.Message)
	}
	if result.Err == nil {
		t.Fatal("Available().Err = nil, want error")
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
