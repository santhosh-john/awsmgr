package aws

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	configFileName      = "config"
	credentialsFileName = "credentials"
)

func CurrentProfile() string {
	if profile := strings.TrimSpace(os.Getenv("AWS_PROFILE")); profile != "" {
		return profile
	}
	return "default"
}

func LoadProfiles(awsDir string) ([]string, error) {
	if awsDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("resolve home directory: %w", err)
		}
		awsDir = filepath.Join(home, ".aws")
	}

	profileSet := make(map[string]struct{})
	for _, source := range []struct {
		path     string
		isConfig bool
	}{
		{path: filepath.Join(awsDir, configFileName), isConfig: true},
		{path: filepath.Join(awsDir, credentialsFileName), isConfig: false},
	} {
		profiles, err := parseProfileFile(source.path, source.isConfig)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}
		for _, profile := range profiles {
			profileSet[profile] = struct{}{}
		}
	}

	profiles := make([]string, 0, len(profileSet))
	for profile := range profileSet {
		profiles = append(profiles, profile)
	}
	sort.Strings(profiles)
	return profiles, nil
}

func parseProfileFile(path string, isConfig bool) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	var profiles []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if !strings.HasPrefix(line, "[") || !strings.Contains(line, "]") {
			continue
		}

		section := strings.TrimSpace(line[1:strings.Index(line, "]")])
		profile := normalizeProfileSection(section, isConfig)
		if profile != "" {
			profiles = append(profiles, profile)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return profiles, nil
}

func normalizeProfileSection(section string, isConfig bool) string {
	if section == "" {
		return ""
	}
	if !isConfig {
		return section
	}
	if section == "default" {
		return section
	}
	if strings.HasPrefix(section, "profile ") {
		return strings.TrimSpace(strings.TrimPrefix(section, "profile "))
	}
	return ""
}
