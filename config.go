package progress

import (
	"os"
	"strings"
)

const configDirectorySuffix = "progress"

// getConfigurationHome gets the root of all configuration files. Not only
// configuration files for this project, but *all* configuration files. For
// instance, it could be /etc/ as opposed /etc/progress/
func getConfigurationHome() string {
	homeDirectory := getHomePath()

	// We will attempt to respect XDG settings, although this may be implemented
	// wrong. Please file an issue ticket if it is - I'm more than happy to fix!`
	xdgConfigurationHome := os.Getenv("XDG_CONFIG_HOME")
	xdgConfigurationPath := os.Getenv("XDG_CONFIG_DIRS")

	// If we have set config home, then we're positive that's where this is going
	if xdgConfigurationHome != "" {
		return xdgConfigurationHome
	}

	if xdgConfigurationPath != "" {
		xdgConfigurationPaths := strings.Split(":", xdgConfigurationPath)

		for _, currentPath := range xdgConfigurationPaths {
			if len(currentPath) < len(homeDirectory) {
				continue
			}

			if currentPath[:len(homeDirectory)] == homeDirectory {
				return currentPath
			}
		}
	}

	return homeDirectory
}

// getConfigurationBasePath returns the path where all project-specific
// configurations should live
func getConfigurationBasePath() string {
	var basePathSegments []string

	// Allow users to choose an exact config root if they want to skip this
	// guessing process
	userRoot := os.Getenv("PRG_CONFIG_PATH")

	if userRoot != "" {
		return userRoot
	}

	basePathSegments = append(basePathSegments, getConfigurationHome())
	basePathSegments = append(basePathSegments, configDirectorySuffix)

	return strings.Join(basePathSegments, string(os.PathSeparator))
}

// GetConfigurationPath gets a path to the configuration directory with the
// given path segments appended to it. If not segments are given, it will
// return the base path wherein process configurations are stored.
func GetConfigurationPath(pathSegments ...string) string {
	var allPathSegments []string

	allPathSegments = append(allPathSegments, getConfigurationBasePath())

	for _, path := range pathSegments {
		allPathSegments = append(allPathSegments, path)
	}

	return strings.Join(allPathSegments, string(os.PathSeparator))
}
