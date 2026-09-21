package openstack

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const cloudsFileName = "clouds.yaml"

func DiscoverCloudsFiles(explicitPath string) ([]string, error) {
	if explicitPath != "" {
		path, err := validateCloudsFile(explicitPath)
		if err != nil {
			return nil, err
		}

		return []string{path}, nil
	}

	paths := []string{
		"/etc/openstack/" + cloudsFileName,
	}

	configDir, err := os.UserConfigDir()
	if err == nil {
		paths = append(paths, filepath.Join(configDir, "openstack", cloudsFileName))
	}

	paths = append(paths, cloudsFileName)

	discovered := make([]string, 0, len(paths))

	for _, path := range paths {
		info, err := os.Stat(path)

		if err == nil {
			if info.IsDir() {
				continue
			}

			absolutePath, err := filepath.Abs(path)
			if err != nil {
				return nil, fmt.Errorf("resolving clouds file %q: %w", path, err)
			}

			discovered = append(discovered, absolutePath)
			continue
		}

		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("checking clouds file %q: %w", path, err)
		}
	}

	if len(discovered) == 0 {
		return nil, fmt.Errorf("clouds.yaml not found")
	}

	return discovered, nil
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
