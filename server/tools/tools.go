package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"log/slog"
)

// searchRootDirectory returns the root directory of the project.
func SearchRootDirectory() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check if the go.mod file exists in the current directory.
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			slog.Debug("Found project root directory.", "Root Directory", currentDir)
			return currentDir, nil
		}

		// Move one directory up.
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// We have reached the root directory and haven't found go.mod.
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("unable to search root directory")
}

func GetCurrentDirectory() string {
	cwd, err := os.Getwd()
	if err != nil {
		slog.Error("Error getting the current working directory:", err)
		return ""
	}
	return cwd
}
