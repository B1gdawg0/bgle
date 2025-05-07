package cmd

import (
	"bgle/config"
	"bgle/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var useCmd = &cobra.Command{
	Use:   "use [project:profile]",
	Short: "🔧 Use a specific profile from the .bgle directory",
	Args:  cobra.ExactArgs(1), // Require exactly one argument
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]

		s := strings.Split(profileName, ":")

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
		if err != nil{
			utils.PrintError("Error getting .bgle directory: " + err.Error())
			return
		}

		filename := filepath.Join(bgleDir, strings.Join(s, ".")+".yaml")

		profile, err := loadProfile(filename)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Error loading profile from %s: %v", filename, err))
			return
		}

		err = applyProfileSettings(profile)
		if err != nil {
			utils.PrintError("Error: "+err.Error())
			return
		}

		utils.PrintSuccess(fmt.Sprintf("Profile %s applied successfully.", profileName))
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func loadProfile(filePath string) (*config.Profile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var profile config.Profile
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func applyProfileSettings(profile *config.Profile) error {
	if profile.Dir != "" {
		err := os.Chdir(profile.Dir)
		if err != nil {
			return fmt.Errorf("error changing directory: %v", err)
		}
		utils.PrintInfo(fmt.Sprintf("📦 Moving to %s...", shortenPath(profile.Dir)))
	}

	if profile.Branch != "" {
		err := runGitCheckout(profile.Branch)
		if err != nil {
			return fmt.Errorf("error checking out branch: %v", err)
		}
		utils.PrintInfo(fmt.Sprintf("Checked out branch: %s", profile.Branch))
	}

	if profile.Docker.Enabled && profile.Docker.Up {
		utils.PrintInfo("🏃Bringing up Docker Compose...")
		err := runDockerCompose(profile.Docker.ComposeFile, profile.EnvFile)
		if err != nil {
			return fmt.Errorf("error running docker-compose: %v", err)
		}
		utils.PrintSuccess("Docker Compose is up.")
	}

	for _, script := range profile.Scripts {
		utils.PrintInfo(fmt.Sprintf("Running script: %s", script))
		err := runScript(script)
		if err != nil {
			return fmt.Errorf("error running script '%s': %v", script, err)
		}
		utils.PrintSuccess(fmt.Sprintf("%s Done!", script))
	}

	return nil
}

func shortenPath(fullPath string) string {
	safePath := filepath.ToSlash(fullPath)
	parts := strings.Split(safePath, "/")
	if len(parts) >= 2 {
		return filepath.Join(parts[len(parts)-2], parts[len(parts)-1])
	}
	return fullPath
}

func runGitCheckout(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDockerCompose(composeFile string, envFile string) error {
	args := []string{"compose", "-f", composeFile}

	if envFile != "" {
		args = append(args, "--env-file", envFile)
	}

	args = append(args, "up", "-d")

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runScript(script string) error {
	cmd := exec.Command("sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
