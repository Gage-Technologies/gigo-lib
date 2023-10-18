package buildinfo

import (
	_ "embed"
	"fmt"
	"regexp"
	"runtime/debug"
	"time"
)

// here we ask git for the current tag at build time
// to ensure that the build is an official release.
// this can be circumvented by creating a local tag
// and never pushing to the repo but hopefully checking
// the commit and clean status will help.
//
// this works by calling the `git describe --tags`
// command at build time via the go:generate system.
// the output (which is the tag) is written to the
// file AUTO-TAG.txt and then embedded into the program
// using go's embed package. we then have to verify that
// the tag is valid so we can trust it's output. we use
// a regular expression to check if the tag fits a
// traditional semantic versioning scheme.
//
// any build that does not have a release grade configuration
// will carry only the git hash and its corresponding build
// info

//go:generate sh -c "printf %s $(git describe --tags) > AUTO-TAG.txt"
//go:embed AUTO-TAG.txt
var tag string

var tagValidator = regexp.MustCompile("^v[0-9]+\\.[0-9]+\\.[0-9]+$")

type BuildInfo struct {
	Release   bool
	Version   string
	BuildTime time.Time
	Commit    string
	Clean     bool
}

func Version() BuildInfo {
	// check if we have a valid version tag
	validTag := tagValidator.MatchString(tag)

	// set default values to retrieve from go's buildinfo system
	commit := ""
	clean := false
	t := time.Unix(0, 0)

	// read build info
	if info, ok := debug.ReadBuildInfo(); ok {
		// iterate build info data
		for _, setting := range info.Settings {
			// set commit for the revision
			if setting.Key == "vcs.revision" {
				commit = setting.Value
			}

			// set time of build
			if setting.Key == "vcs.time" {
				parsed, err := time.Parse(time.RFC3339, setting.Value)
				// silently skip time parse error - this isn't really important just a nice to have
				if err != nil {
					continue
				}
				t = parsed
			}

			// check if this is a dirty build (has uncommited changes)
			if setting.Key == "vcs.modified" {
				clean = setting.Value == "false"
			}
		}
	}

	// handle release
	if clean && validTag {
		return BuildInfo{
			Release:   true,
			Version:   tag,
			BuildTime: t,
			Commit:    commit,
			Clean:     true,
		}
	}

	// handle non release
	cleanStatus := "dirty"
	if clean {
		cleanStatus = "clean"
	}
	shortCommit := "unknown"
	if commit != "" {
		shortCommit = commit[:7]
	}
	return BuildInfo{
		Release:   false,
		Version:   fmt.Sprintf("dev-%s-%s", shortCommit, cleanStatus),
		BuildTime: t,
		Commit:    commit,
		Clean:     clean,
	}
}
