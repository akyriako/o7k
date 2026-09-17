package openstack

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const cloudsFileName = "clouds.yaml"

func DiscoverCloudsFile(explicitPath string) (string, error) {
	if explicitPath != "" {
		return validateCloudsFile(explicitPath)
	}

	paths := []string{
		"/etc/openstack/" + cloudsFileName,
	}

	configDir, err := os.UserConfigDir()
	if err == nil {
		paths = append(paths, filepath.Join(configDir, "openstack", cloudsFileName))
	}

	paths = append(paths, cloudsFileName)

	for _, path := range paths {
		info, err := os.Stat(path)

		if err == nil {
			if info.IsDir() {
				continue
			}

			return filepath.Abs(path)
		}

		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("checking clouds file %q: %w", path, err)
		}
	}

	return "", fmt.Errorf("clouds.yaml not found")
}

func validateCloudsFile(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("invalid clouds file %q: %w", path, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("clouds file %q is a directory", path)
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolving clouds file %q: %w", path, err)
	}

	return absolutePath, nil
}
