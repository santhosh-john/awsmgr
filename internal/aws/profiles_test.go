package aws

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadProfilesDeduplicatesAndSorts(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "config"), `
[default]
region = us-east-1

[profile prod]
sso_session = corp

[profile dev]
region = us-west-2

[sso-session corp]
sso_region = us-east-1

[services corp]
s3 =
`)
	writeFile(t, filepath.Join(dir, "credentials"), `
[prod]
aws_access_key_id = example

[sandbox]
aws_access_key_id = example
`)

	got, err := LoadProfiles(dir)
	if err != nil {
		t.Fatalf("LoadProfiles returned error: %v", err)
	}

	want := []string{"default", "dev", "prod", "sandbox"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadProfiles() = %v, want %v", got, want)
	}
}

func TestCurrentProfileFallsBackToDefault(t *testing.T) {
	t.Setenv("AWS_PROFILE", "")

	if got := CurrentProfile(); got != "default" {
		t.Fatalf("CurrentProfile() = %q, want default", got)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}
