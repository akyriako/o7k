package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/mod/semver"
)

type UpdateResult struct {
	Updated     bool
	FromVersion string
	ToVersion   string
}

func Update(ctx context.Context, plugin Info) (UpdateResult, error) {
	result := UpdateResult{
		FromVersion: plugin.Version,
		ToVersion:   plugin.Version,
	}

	if plugin.Path == "" {
		return result, fmt.Errorf("plugin path is empty")
	}

	if plugin.URL == "" {
		return result, fmt.Errorf("plugin %q does not provide an update URL", plugin.Name)
	}

	repository, err := parseGitHubRepository(plugin.URL)
	if err != nil {
		return result, err
	}

	release, err := getLatestGitHubRelease(ctx, repository)
	if err != nil {
		return result, err
	}

	newer, err := hasNewerVersion(plugin.Version, release.TagName)
	if err != nil {
		return result, err
	}

	if !newer {
		return result, nil
	}

	result.ToVersion = strings.TrimPrefix(release.TagName, "v")

	asset, err := findReleaseAsset(release)
	if err != nil {
		return result, err
	}

	body, err := download(ctx, asset.BrowserDownloadURL)
	if err != nil {
		return result, err
	}
	defer body.Close()

	pluginDir, err := useUserPluginsDir()
	if err != nil {
		return result, err
	}

	tempPath, err := writePluginTemp(pluginDir, body)
	if err != nil {
		return result, err
	}
	defer os.Remove(tempPath)

	backupDir, err := useUserPluginsBackupDir()
	if err != nil {
		return result, err
	}

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return result, fmt.Errorf("creating plugin backup directory: %w", err)
	}

	binaryName := pluginBinaryName(asset.Name)

	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		return result, fmt.Errorf("reading plugin directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Join(pluginDir, entry.Name()) == plugin.Path {
			continue
		}

		if pluginBinaryName(entry.Name()) != binaryName {
			continue
		}

		if err := os.Remove(filepath.Join(pluginDir, entry.Name())); err != nil {
			return result, fmt.Errorf("removing previous plugin version: %w", err)
		}
	}

	backupPath := filepath.Join(backupDir, filepath.Base(plugin.Path))
	destination := filepath.Join(pluginDir, asset.Name)

	if err := os.Remove(backupPath); err != nil && !os.IsNotExist(err) {
		return result, fmt.Errorf("removing previous plugin backup: %w", err)
	}

	if err := os.Rename(plugin.Path, backupPath); err != nil {
		return result, fmt.Errorf("backing up plugin: %w", err)
	}

	if err := os.Rename(tempPath, destination); err != nil {
		if restoreErr := os.Rename(backupPath, plugin.Path); restoreErr != nil {
			return result, fmt.Errorf("installing plugin: %w; restoring backup: %v", err, restoreErr)
		}

		return result, fmt.Errorf("installing plugin: %w", err)
	}

	result.Updated = true

	return result, nil
}

type githubRepository struct {
	Owner string
	Name  string
}

type githubRelease struct {
	TagName string               `json:"tag_name"`
	Assets  []githubReleaseAsset `json:"assets"`
}

type githubReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func parseGitHubRepository(rawURL string) (githubRepository, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return githubRepository{}, fmt.Errorf("parsing plugin URL %q: %w", rawURL, err)
	}

	if u.Scheme != "https" || u.Host != "github.com" {
		return githubRepository{}, fmt.Errorf("unsupported plugin URL %q", rawURL)
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return githubRepository{}, fmt.Errorf("invalid GitHub repository URL %q", rawURL)
	}

	return githubRepository{
		Owner: parts[0],
		Name:  strings.TrimSuffix(parts[1], ".git"),
	}, nil
}

func getLatestGitHubRelease(ctx context.Context, repository githubRepository) (githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repository.Owner, repository.Name)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return githubRelease{}, fmt.Errorf("creating GitHub release request: %w", err)
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("User-Agent", "o7k")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return githubRelease{}, fmt.Errorf("getting latest GitHub release: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return githubRelease{}, fmt.Errorf("getting latest GitHub release: %s", response.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(response.Body).Decode(&release); err != nil {
		return githubRelease{}, fmt.Errorf("decoding latest GitHub release: %w", err)
	}

	return release, nil
}

func normalizeVersion(version string) string {
	version = strings.TrimSpace(version)

	if version != "" && !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	return version
}

func hasNewerVersion(current, latest string) (bool, error) {
	currentVersion := normalizeVersion(current)
	latestVersion := normalizeVersion(latest)

	if !semver.IsValid(currentVersion) {
		return false, fmt.Errorf("invalid installed plugin version %q", current)
	}

	if !semver.IsValid(latestVersion) {
		return false, fmt.Errorf("invalid release version %q", latest)
	}

	return semver.Compare(latestVersion, currentVersion) > 0, nil
}

func findReleaseAsset(release githubRelease) (githubReleaseAsset, error) {
	suffix := fmt.Sprintf("_%s_%s", runtime.GOOS, runtime.GOARCH)

	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, suffix) {
			return asset, nil
		}
	}

	return githubReleaseAsset{}, fmt.Errorf("release %q has no asset for %s/%s", release.TagName, runtime.GOOS, runtime.GOARCH)
}

func pluginBinaryName(name string) string {
	suffix := fmt.Sprintf("_%s_%s", runtime.GOOS, runtime.GOARCH)
	name = strings.TrimSuffix(name, suffix)

	if i := strings.LastIndex(name, "_"); i >= 0 {
		name = name[:i]
	}

	return name
}
