package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "0.8.1"
	GitCommit = "dev"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)

// Full returns a formatted multi-line version and build info string.
func Full() string {
	return fmt.Sprintf("spp-shadowsocks-plugin version %s\ngit commit: %s\nbuild time: %s\ngo version: %s\nos/arch:    %s/%s",
		Version, GitCommit, BuildTime, GoVersion, runtime.GOOS, runtime.GOARCH)
}

// Short returns a concise single-line version string.
func Short() string {
	return fmt.Sprintf("spp-shadowsocks-plugin %s (%s, %s/%s)", Version, GitCommit, runtime.GOOS, runtime.GOARCH)
}
