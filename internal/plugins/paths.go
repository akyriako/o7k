package plugins

import (
	"fmt"
	"os"
	"path/filepath"
)

func useUserPluginsDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("getting user config directory: %w", err)
	}

	return filepath.Join(configDir, "o7k", "plugins"), nil
}

func useUserPluginsBackupDir() (string, error) {
	dir, err := useUserPluginsDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "backup"), nil
}
