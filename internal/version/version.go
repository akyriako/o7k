package version

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const (
	defaultVersion   = "v0.0.0-dev"
	defaultCommit    = "none"
	defaultBuildDate = "unknown"
	defaultValue     = "unknown"
)

var (
	Version   = defaultVersion
	Commit    = defaultCommit
	BuildDate = defaultBuildDate
)

type VersionInfo struct {
	GoVersion string `json:"goVersion"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
	Modified  string `json:"modified"`
}

func GetBuildInfo() VersionInfo {
	binfo, ok := debug.ReadBuildInfo()

	if Version != defaultVersion || Commit != defaultCommit || BuildDate != defaultBuildDate {
		goVer := defaultValue
		modified := defaultValue

		if ok {
			goVer = binfo.GoVersion
			modified = getBuildInfoSetting(binfo, "vcs.modified")
		}

		return VersionInfo{
			GoVersion: goVer,
			Version:   Version,
			Commit:    Commit,
			BuildDate: BuildDate,
			Modified:  modified,
		}
	}

	if !ok {
		return VersionInfo{
			GoVersion: defaultValue,
			Version:   defaultVersion,
			Commit:    defaultCommit,
			BuildDate: defaultBuildDate,
			Modified:  defaultValue,
		}
	}

	return VersionInfo{
		GoVersion: binfo.GoVersion,
		Version:   binfo.Main.Version,
		Commit:    getBuildInfoSetting(binfo, "vcs.revision"),
		BuildDate: getBuildInfoSetting(binfo, "vcs.time"),
		Modified:  getBuildInfoSetting(binfo, "vcs.modified"),
	}
}

func getBuildInfoSetting(info *debug.BuildInfo, key string) string {
	for _, setting := range info.Settings {
		if setting.Key == key {
			return setting.Value
		}
	}

	return defaultValue
}

const (
	githubOwner = "akyriako"
	githubRepo  = "o7k"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

type UpdateInfo struct {
	Available      bool
	CurrentVersion string
	LatestVersion  string
	URL            string
}

func CheckForUpdate(ctx context.Context) (UpdateInfo, error) {
	current := Version

	//if current == defaultVersion {
	//	return UpdateInfo{
	//		CurrentVersion: current,
	//	}, nil
	//}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", githubOwner, githubRepo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return UpdateInfo{}, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", githubRepo)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return UpdateInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UpdateInfo{}, fmt.Errorf("github returned status %s", resp.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return UpdateInfo{}, err
	}

	return UpdateInfo{
		Available:      compareVersions(current, release.TagName),
		CurrentVersion: current,
		LatestVersion:  release.TagName,
		URL:            release.HTMLURL,
	}, nil
}

func compareVersions(current, latest string) bool {
	current = normalizeSemver(current)
	latest = normalizeSemver(latest)

	if !semver.IsValid(current) || !semver.IsValid(latest) {
		return false
	}

	return semver.Compare(latest, current) > 0
}

func normalizeSemver(v string) string {
	v = strings.TrimSpace(v)
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	return v
}
