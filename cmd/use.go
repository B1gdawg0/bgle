package cmd

import (
	"bgle/config"
	"bgle/models"
	"bgle/utils"
	"fmt"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var useCmd = &cobra.Command{
	Use:   "use [project:profile]",
	Short: "🔧 Use a specific profile from the .bgle directory or register.yaml",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		currentDir, err := os.Getwd()
		if err != nil {
			utils.PrintError(fmt.Sprintf("Unable to get current directory: %v", err))
			return
		}
		name, profileName, err := utils.ParseNameProfile(args[0])
		if err != nil {
			utils.PrintError("Invalid format. Use project:profile - " + err.Error())
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			utils.PrintError("Error getting home directory: " + err.Error())
			return
		}

		bgleDir := os.Getenv("BGLE_HOME")
		if bgleDir == "" {
			bgleDir = filepath.Join(homeDir, ".bgle")
		}

		err = utils.EnsureBgleDirectoryAt(bgleDir)
		if err != nil {
			utils.PrintError("Error ensuring .bgle directory: " + err.Error())
			return
		}

		registerFile := filepath.Join(bgleDir, "register.yaml")
		profileEntity, err := loadProfileFromRegister(registerFile, name, profileName)
		if err != nil {
			// utils.PrintError("Profile not found in register.yaml, falling back to .bgle directory...")
			// // Fallback to loading profile from .bgle directory
			// filename := filepath.Join(bgleDir, fmt.Sprintf("%s.%s.yaml", name, profileName))
			// profile, err = loadProfile(filename)
			// if err != nil {
			// 	utils.PrintError(fmt.Sprintf("Error loading profile from %s: %v", filename, err))
			// 	return
			// }

			utils.PrintError(fmt.Sprintf("Error: %v", err))
			return
		}

		filename := filepath.Join(bgleDir, fmt.Sprintf("projects/%s.%s.yaml", name, profileName))
		profile, err := loadProfile(filename)

		profile.Project = name
		profile.Profile = profileEntity.Profile
		profile.Dir = profileEntity.Dir
		profile.Bootstrap.Enabled = profile.Dir == "."

		if err != nil {
			utils.PrintError(fmt.Sprintf("Error loading profile from %s: %v", filename, err))
			return
		}

		// Apply the profile settings
		err = config.ApplyProfileSettings(profile)
		if err != nil {
			utils.PrintError("Error applying profile: " + err.Error())
			return
		}

		utils.PrintSuccess(fmt.Sprintf("Profile %s:%s applied successfully.", name, profileName))

		utils.CompareCurrentDirWithProfileDir(currentDir, profile.Dir)
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func loadProfile(filePath string) (*models.Profile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var profile models.Profile
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func loadProfileFromRegister(registerFile, projectName, profileName string) (*models.ProfileEntry, error) {
	// Load the register.yaml file
	file, err := os.Open(registerFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var register struct {
		Profiles map[string]models.ProfileEntry `yaml:"profiles"`
	}

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&register)
	if err != nil {
		return nil, err
	}

	// Check if the projectName exists in the register
	profileEntry, found := register.Profiles[projectName]
	if !found {
		return nil, fmt.Errorf("project %s not found in register.yaml", projectName)
	}

	// Now, check if the profile matches within the project
	if profileEntry.Profile != profileName {
		return nil, fmt.Errorf("profile %s not found under project %s in register.yaml", profileName, projectName)
	}

	return &profileEntry, nil
}
