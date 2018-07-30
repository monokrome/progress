package progress

import "os"

func getHomePath() string {
	homeDirectory := os.Getenv("HOME")

	if homeDirectory != "" {
		return homeDirectory
	}

	// TODO: Support Windows paths
	return "/"
}
