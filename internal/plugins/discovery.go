package plugins

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func Discover() ([]string, error) {
	dir, err := useUserPluginsDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}

		return nil, err
	}

	var paths []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.Mode().IsRegular() && info.Mode().Perm()&0111 != 0 {
			paths = append(paths, filepath.Join(dir, entry.Name()))
		}
	}

	return paths, nil
}
