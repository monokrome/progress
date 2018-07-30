package progress

import (
	"os"
	"strings"
)

const dataDirectorySuffix = "progress"

// getDataHome gets the root of all configuration files. Not only
// configuration files for this project, but *all* configuration files. For
// instance, it could be /etc/ as opposed /etc/progress/
func getDataHome() string {
	homeDirectory := getHomePath()

	// We will attempt to respect XDG settings, although this may be implemented
	// wrong. Please file an issue ticket if it is - I'm more than happy to fix!`
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	xdgDataPath := os.Getenv("XDG_DATA_DIRS")

	// If we have set config home, then we're positive that's where this is going
	if xdgDataHome != "" {
		return xdgDataHome
	}

	if xdgDataPath != "" {
		xdgDataPaths := strings.Split(":", xdgDataPath)

		for _, currentPath := range xdgDataPaths {
			if currentPath[:len(homeDirectory)] == homeDirectory {
				return currentPath
			}
		}
	}

	return homeDirectory
}

// getDataBasePath returns the path where all project-specific
// configurations should live
func getDataBasePath() string {
	var basePathSegments []string

	// Allow users to choose an exact config root if they want to skip this
	// guessing process
	userRoot := os.Getenv("PRG_DATA_PATH")

	if userRoot != "" {
		return userRoot
	}

	basePathSegments = append(basePathSegments, getDataHome())
	basePathSegments = append(basePathSegments, configDirectorySuffix)

	return strings.Join(basePathSegments, string(os.PathSeparator))
}

// GetDataPath returns the name of the directory where
func GetDataPath(pathSegments ...string) string {
	var allPathSegments []string
	allPathSegments = append(allPathSegments, getDataBasePath())

	for _, path := range pathSegments {
		allPathSegments = append(allPathSegments, path)
	}

	return strings.Join(allPathSegments, string(os.PathSeparator))
}
