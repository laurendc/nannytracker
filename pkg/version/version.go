package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the current version of the application
	// This should be set during the build process
	Version = "dev"

	// BuildTime is the time when the binary was built
	// This should be set during the build process
	BuildTime = "unknown"

	// GitCommit is the git commit hash
	// This should be set during the build process
	GitCommit = "unknown"
)

// Info contains version information
type Info struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// Get returns the current version information
func Get() Info {
	return Info{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String returns a formatted version string
func String() string {
	info := Get()
	return fmt.Sprintf("NannyTracker %s (%s) built on %s",
		info.Version, info.GitCommit, info.BuildTime)
}

// FullString returns a detailed version string
func FullString() string {
	info := Get()
	return fmt.Sprintf("NannyTracker %s\n"+
		"  Build Time: %s\n"+
		"  Git Commit: %s\n"+
		"  Go Version: %s\n"+
		"  OS/Arch: %s/%s",
		info.Version, info.BuildTime, info.GitCommit,
		info.GoVersion, info.OS, info.Arch)
}
