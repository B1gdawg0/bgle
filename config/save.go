package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"bgle/models"
	"bgle/utils"
)

func SaveProfile(profile models.Profile) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %v", err)
	}

	bgleDir := filepath.Join(homeDir, ".bgle")

	if err := utils.EnsureBgleDirectoryAt(bgleDir); err != nil {
		return fmt.Errorf("error ensuring .bgle directory: %v", err)
	}

	projectDir := filepath.Join(bgleDir, "projects")
	if err := os.MkdirAll(projectDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating project directory: %v", err)
	}

	filePath := filepath.Join(projectDir, fmt.Sprintf("%s.%s.yaml", profile.Project, profile.Profile))

	// Check if the profile already exists
	if _, err := os.Stat(filePath); err == nil {
		// Profile already exists, ask for user confirmation
		utils.PrintWarning(fmt.Sprintf("Profile %s already exists.", filePath))
		overwrite, err := utils.AskConfirmation("Do you want to overwrite it?")
		if err != nil {
			return fmt.Errorf("error reading input: %v", err)
		}

		if !overwrite {
			// User does not want to overwrite
			utils.PrintInfo("Profile not overwritten.")
			return nil
		}
	} else if os.IsNotExist(err) {
		// File doesn't exist, proceed with creation
		utils.PrintInfo(fmt.Sprintf("Creating new profile %s...", filePath))
	} else {
		// If there is another error (e.g., permission issues), return it
		return fmt.Errorf("error checking if file exists: %v", err)
	}

	// Create or overwrite the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	output := models.OutputProfile{
        Branch:      profile.Branch,
        Docker:      profile.Docker,
        EnvFile:     profile.EnvFile,
        EnvVars:     profile.EnvVars,
        Bootstrap:   profile.Bootstrap,
        PreScripts:  profile.PreScripts,
        Scripts:     profile.Scripts,
        PostScripts: profile.PostScripts,
    }
    
    if err := encoder.Encode(output); err != nil {
        return fmt.Errorf("error encoding profile: %v", err)
    }

	if err := UpdateRegisterFile(profile); err != nil {
		return fmt.Errorf("error updating register file: %v", err)
	}

	utils.PrintSuccess(fmt.Sprintf("Profile for %s has been saved and registered.", profile.Project))
	return nil
}