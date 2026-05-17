package kube

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestCurrentContextReportsMissingKubectl(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	result := CurrentContext(context.Background(), time.Second)
	if result.OK {
		t.Fatal("CurrentContext().OK = true, want false")
	}
	if result.Message != "kubectl not found" {
		t.Fatalf("CurrentContext().Message = %q, want kubectl not found", result.Message)
	}
}

func TestCurrentContextReturnsActiveContext(t *testing.T) {
	binDir := t.TempDir()
	writeExecutable(t, filepath.Join(binDir, executableName("kubectl")), `if [ "$1" = "config" ] && [ "$2" = "current-context" ]; then
  printf 'prod-cluster\n'
  exit 0
fi
exit 1
`, `@echo off
if "%1"=="config" if "%2"=="current-context" (
  echo prod-cluster
  exit /b 0
)
exit /b 1
`)
	t.Setenv("PATH", binDir+string(filepath.ListSeparator)+os.Getenv("PATH"))

	result := CurrentContext(context.Background(), time.Second)
	if !result.OK {
		t.Fatalf("CurrentContext().OK = false, want true; result=%#v", result)
	}
	if result.Context != "prod-cluster" {
		t.Fatalf("CurrentContext().Context = %q, want prod-cluster", result.Context)
	}
	if result.Message != "" {
		t.Fatalf("CurrentContext().Message = %q, want empty", result.Message)
	}
}

func TestCurrentContextRejectsEmptyContext(t *testing.T) {
	binDir := t.TempDir()
	writeExecutable(t, filepath.Join(binDir, executableName("kubectl")), `if [ "$1" = "config" ] && [ "$2" = "current-context" ]; then
  printf '\n'
  exit 0
fi
exit 1
`, `@echo off
if "%1"=="config" if "%2"=="current-context" (
  echo.
  exit /b 0
)
exit /b 1
`)
	t.Setenv("PATH", binDir+string(filepath.ListSeparator)+os.Getenv("PATH"))

	result := CurrentContext(context.Background(), time.Second)
	if result.OK {
		t.Fatal("CurrentContext().OK = true, want false")
	}
	if result.Message != "kubectl context unavailable" {
		t.Fatalf("CurrentContext().Message = %q, want kubectl context unavailable", result.Message)
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
