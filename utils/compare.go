package utils

import (
	"fmt"
	"path/filepath"
)

func CompareCurrentDirWithProfileDir(currentDir, profileDir string) {
	absCurrentDir, err := filepath.Abs(currentDir)
	if err != nil {
		PrintError(fmt.Sprintf("Unable to resolve current directory: %v", err))
		return
	}

	absProfileDir, err := filepath.Abs(profileDir)
	if err != nil {
		PrintError(fmt.Sprintf("Unable to resolve profile directory: %v", err))
		return
	}

	if absCurrentDir == absProfileDir {
		return
	}

	rel, err := filepath.Rel(absCurrentDir, absProfileDir)
	if err != nil || rel == "" {
		PrintBoxedWarning(fmt.Sprintf("You're not in the project directory.\nTo go there:\n  cd %s", absProfileDir))
	} else {
		PrintBoxedWarning(fmt.Sprintf("You're not in the project directory.\nTo go there:\n  cd %s", rel))
	}
}
