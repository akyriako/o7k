package version

import "runtime/debug"

const (
	defaultVersion   = "dev"
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
