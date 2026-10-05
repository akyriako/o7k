package plugins

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const maxPluginSize = 100 << 20 // 100 MiB

func Install(ctx context.Context, source string) error {
	u, err := url.Parse(source)
	if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		return installRemote(ctx, u)
	}

	return installLocal(source)
}

func installLocal(source string) error {
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("opening plugin %q: %w", source, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat plugin %q: %w", source, err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("plugin %q is not a regular file", source)
	}

	return install(filepath.Base(source), file)
}

func installRemote(ctx context.Context, source *url.URL) error {
	body, err := download(ctx, source.String())
	if err != nil {
		return err
	}
	defer body.Close()

	name := filepath.Base(source.Path)
	if name == "." || name == "/" || name == "" {
		return fmt.Errorf("plugin URL does not contain a filename")
	}

	return install(name, body)
}

func install(name string, source io.Reader) error {
	pluginDir, err := useUserPluginsDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("creating plugin directory: %w", err)
	}

	tempPath, err := writePluginTemp(pluginDir, source)
	if err != nil {
		return err
	}
	defer os.Remove(tempPath)

	binaryName := pluginBinaryName(name)

	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		return fmt.Errorf("reading plugin directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == name {
			continue
		}

		if pluginBinaryName(entry.Name()) != binaryName {
			continue
		}

		if err := os.Remove(filepath.Join(pluginDir, entry.Name())); err != nil {
			return fmt.Errorf("removing previous plugin version: %w", err)
		}
	}

	destination := filepath.Join(pluginDir, name)

	if err := os.Rename(tempPath, destination); err != nil {
		return fmt.Errorf("installing plugin: %w", err)
	}

	fmt.Printf("installed plugin %s\n", destination)

	return nil
}

func Remove(name string) error {
	if name == "" || filepath.Base(name) != name {
		return fmt.Errorf("invalid plugin name %q", name)
	}

	pluginDir, err := useUserPluginsDir()
	if err != nil {
		return err
	}

	path := filepath.Join(pluginDir, name)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("removing plugin %q: %w", name, err)
	}

	fmt.Printf("removed plugin %s\n", path)

	return nil
}

func download(ctx context.Context, source string) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("downloading plugin: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		response.Body.Close()
		return nil, fmt.Errorf("downloading plugin: %s", response.Status)
	}

	return response.Body, nil
}

func writePluginTemp(pluginDir string, source io.Reader) (string, error) {
	temp, err := os.CreateTemp(pluginDir, ".o7k-plugin-*")
	if err != nil {
		return "", fmt.Errorf("creating temporary plugin: %w", err)
	}

	tempPath := temp.Name()

	written, err := io.Copy(temp, io.LimitReader(source, maxPluginSize+1))
	if err != nil {
		temp.Close()
		os.Remove(tempPath)
		return "", fmt.Errorf("writing plugin: %w", err)
	}

	if written > maxPluginSize {
		temp.Close()
		os.Remove(tempPath)
		return "", fmt.Errorf("plugin exceeds maximum size of %d MiB", maxPluginSize>>20)
	}

	if err := temp.Chmod(0755); err != nil {
		temp.Close()
		os.Remove(tempPath)
		return "", fmt.Errorf("making plugin executable: %w", err)
	}

	if err := temp.Close(); err != nil {
		os.Remove(tempPath)
		return "", fmt.Errorf("closing plugin: %w", err)
	}

	return tempPath, nil
}
