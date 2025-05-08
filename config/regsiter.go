package config

import (
	"bgle/models"
	"bgle/utils"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func UpdateRegisterFile(profile models.Profile) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	bgleDir := filepath.Join(homeDir, ".bgle")
	registerFilePath := filepath.Join(bgleDir, "register.yaml")

	// Ensure .bgle directory exists
	if err := utils.EnsureBgleDirectoryAt(bgleDir); err != nil {
		return fmt.Errorf("failed to create .bgle directory: %v", err)
	}

	var register models.Register

	// Load existing register.yaml if it exists
	if _, err := os.Stat(registerFilePath); err == nil {
		file, err := os.Open(registerFilePath)
		if err != nil {
			return fmt.Errorf("error opening register.yaml: %v", err)
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&register); err != nil {
			return fmt.Errorf("error decoding register.yaml: %v", err)
		}
	}

	// Initialize map if nil
	if register.Profiles == nil {
		register.Profiles = make(map[string]models.ProfileEntry)
	}

	// Update or add the profile entry
	register.Profiles[profile.Project] = models.ProfileEntry{
		Profile: profile.Profile,
		Dir:     profile.Dir,
	}

	// Save back to register.yaml
	file, err := os.Create(registerFilePath)
	if err != nil {
		return fmt.Errorf("error creating register.yaml: %v", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	if err := encoder.Encode(register); err != nil {
		return fmt.Errorf("error encoding register.yaml: %v", err)
	}

	utils.PrintSuccess(fmt.Sprintf("Register updated for project %s.", profile.Project))
	return nil
}
