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

	// Check if the register.yaml file exists
	var register models.Register
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

	// Check if the profile already exists in the register
	if _, found := register.Profiles[profile.Project]; found {
		utils.PrintInfo("Profile already registered.")
		return nil
	}

	// Add the new profile entry to the register
	newEntry := models.ProfileEntry{
		Profile: profile.Profile,
		Dir:     profile.Dir,
	}

	// Initialize Profiles map if nil
	if register.Profiles == nil {
		register.Profiles = make(map[string]models.ProfileEntry)
	}

	// Add profile under the project name as the key
	register.Profiles[profile.Project] = newEntry

	// Create or open the register.yaml file
	file, err := os.Create(registerFilePath)
	if err != nil {
		return fmt.Errorf("error creating register.yaml: %v", err)
	}
	defer file.Close()

	// Write the updated register to the file
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	if err := encoder.Encode(register); err != nil {
		return fmt.Errorf("error encoding register.yaml: %v", err)
	}

	utils.PrintSuccess(fmt.Sprintf("Profile for %s has been registered.", profile.Project))
	return nil
}
