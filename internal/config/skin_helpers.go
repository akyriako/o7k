package config

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed skins/*.yaml
var embeddedSkins embed.FS

func installSkins(configDir string) error {
	skinsDir := filepath.Join(configDir, "skins")

	if err := os.MkdirAll(skinsDir, 0755); err != nil {
		return fmt.Errorf("creating skins directory: %w", err)
	}

	files, err := fs.Glob(embeddedSkins, "skins/*.yaml")
	if err != nil {
		return fmt.Errorf("listing embedded skins: %w", err)
	}

	for _, file := range files {
		data, err := embeddedSkins.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading embedded skin %s: %w", file, err)
		}

		path := filepath.Join(skinsDir, filepath.Base(file))

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if os.IsExist(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("creating skin %s: %w", path, err)
		}

		_, writeErr := f.Write(data)
		closeErr := f.Close()

		if writeErr != nil {
			return fmt.Errorf("writing skin %s: %w", path, writeErr)
		}
		if closeErr != nil {
			return fmt.Errorf("closing skin %s: %w", path, closeErr)
		}
	}

	return nil
}

func loadSkin(configDir, name string) (Skin, error) {
	var skin Skin

	if name == "" || name != filepath.Base(name) ||
		strings.ContainsAny(name, `/\`) || name == "." || name == ".." {
		return skin, fmt.Errorf("invalid skin name %q", name)
	}

	path := filepath.Join(configDir, "skins", name+".yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return skin, fmt.Errorf("reading skin %q: %w", name, err)
	}

	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	decoder.KnownFields(true)

	if err := decoder.Decode(&skin); err != nil {
		return skin, fmt.Errorf("parsing skin %q: %w", name, err)
	}

	return skin, nil
}
