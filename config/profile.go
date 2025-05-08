package config

import (
	"bgle/models"
	"bgle/utils"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func selectDockerComposeFiles() []string {
    files, _ := filepath.Glob("docker-compose*.yml")
    if len(files) == 0 {
        fmt.Println("No docker-compose files found.")
        return nil
    }

    fmt.Println("Select Docker Compose files (comma-separated numbers):")
    for i, file := range files {
        fmt.Printf("%d) %s\n", i+1, file)
    }

    fmt.Print("Your choice: ")
    var input string
    fmt.Scanln(&input)

    var selected []string
    indexes := strings.Split(input, ",")
    for _, idx := range indexes {
        i, err := strconv.Atoi(strings.TrimSpace(idx))
        if err == nil && i >= 1 && i <= len(files) {
            selected = append(selected, files[i-1])
        }
    }

    return selected
}

func InteractiveProfileCreation(name string, profile string) models.Profile {
    utils.PrintInfo("Starting profile creation...")
    reader := bufio.NewReader(os.Stdin)
    dir, _ := os.Getwd()

	fmt.Print("Would you like to set up a bootstrap script? (y/n): ")
	bootstrapScript, _ := reader.ReadString('\n')
	wannaSetup := strings.TrimSpace(strings.ToLower(bootstrapScript)) == "y"

	var bootstrap models.Bootstrap

	if wannaSetup {
		fmt.Print("Enter Git repository URL: ")
		repo, _ := reader.ReadString('\n')
		bootstrap.Repo_URL = strings.TrimSpace(repo)

		fmt.Println("Enter bootstrap scripts (type 'exit' to finish):")
		for {
			fmt.Print("> ")
			entry, _ := reader.ReadString('\n')
			entry = strings.TrimSpace(entry)

			if entry == "exit" {
				break
			}

			if entry != "" {
				bootstrap.Scripts = append(bootstrap.Scripts, entry)
			}
		}
	}
	

    fmt.Print("Git branch (optional): ")
    branch, _ := reader.ReadString('\n')
    branch = strings.TrimSpace(branch)

    fmt.Print("Use Docker? (y/n): ")
    dockerInput, _ := reader.ReadString('\n')
    docker := strings.TrimSpace(dockerInput) == "y"

    var dockerConfig models.Docker
    if docker {
        files := selectDockerComposeFiles()
        fmt.Print("Run docker-compose up -d? (y/n): ")
        upInput, _ := reader.ReadString('\n')
        up := strings.TrimSpace(upInput) == "y"

        // fmt.Print("Run docker build? (y/n): ")
        // buildInput, _ := reader.ReadString('\n')
        // build := strings.TrimSpace(buildInput) == "y"

        dockerConfig = models.Docker{
            Enabled:      true,
            ComposeFiles: files,
            Up:           up,
            // Build:        build,
        }
    }

    fmt.Print("Path to .env file (optional): ")
    envFile, _ := reader.ReadString('\n')
    envFile = strings.TrimSpace(envFile)

    // ENV VARS
    fmt.Println("Add env vars (KEY=VALUE). Type 'exit' to stop:")
    envVars := make(map[string]string)
    for {
        fmt.Print("> ")
        entry, _ := reader.ReadString('\n')
        entry = strings.TrimSpace(entry)
        if entry == "exit" {
            break
        }
        parts := strings.SplitN(entry, "=", 2)
        if len(parts) == 2 {
            envVars[parts[0]] = parts[1]
        }
    }

    // SCRIPTS
    readScripts := func(title string) []string {
        fmt.Printf("Enter %s one by one. Type 'exit' to finish:\n", title)
        var s []string
        for {
            fmt.Print("> ")
            line, _ := reader.ReadString('\n')
            line = strings.TrimSpace(line)
            if line == "exit" {
                break
            }
            if line != ""{
				s = append(s, line)
			}
        }
        return s
    }

    preScripts := readScripts("pre-scripts")
    var mainScripts []string

    if !dockerConfig.Enabled {
        mainScripts = readScripts("main scripts")
    }
    postScripts := readScripts("post-scripts")

    profileConfig := models.Profile{
        Project:     name,
        Profile:     profile,
        Dir:         dir,
        Branch:      branch,
        Docker:      dockerConfig,
        EnvFile:     envFile,
        EnvVars:     envVars,
		Bootstrap:   bootstrap,
        PreScripts:  preScripts,
        Scripts:     mainScripts,
        PostScripts: postScripts,
    }

    return profileConfig
}

func ApplyProfileSettings(profile *models.Profile) error {
	// 0. Set up project if it's don't exist
	if profile.Bootstrap.Enabled && profile.Bootstrap.Repo_URL != ""{
		utils.PrintInfo("Start cloning repository...")
		if err := runScript("git clone "+profile.Bootstrap.Repo_URL); err != nil{
			return fmt.Errorf("can't clone project from '%s' repository: %v", profile.Bootstrap.Repo_URL, err)
		}

		dirLocal, err := createDirNameFromRepo(profile.Bootstrap.Repo_URL)

		if err != nil{
			return fmt.Errorf("%s", "Error: "+err.Error())
		}

		profile.Dir = dirLocal

		for _, script := range profile.Bootstrap.Scripts {
			utils.PrintInfo(fmt.Sprintf("Running bootstrap: %s", script))
			if err := runScript(script); err != nil {
				return fmt.Errorf("error running bootstrap '%s': %v", script, err)
			}
		}
	}

	// 1. Change directory
	if profile.Dir != "" {
		if profile.Dir == "."{ return fmt.Errorf("no directory to access. please complete 'bgle sync %s:%s' first", profile.Project, profile.Profile)}
		err := os.Chdir(profile.Dir)
		if err != nil {
			return fmt.Errorf("error changing directory: %v", err)
		}
		utils.PrintInfo(fmt.Sprintf("📦 Changed directory to %s", shortenPath(profile.Dir)))
	}

	// 2. Git checkout
	if profile.Branch != "" {
		err := runGitCheckout(profile.Branch)
		if err != nil {
			return fmt.Errorf("error checking out branch: %v", err)
		}
		utils.PrintInfo(fmt.Sprintf("Checked out branch: %s", profile.Branch))
	}

	if profile.Bootstrap.Enabled {
		if err := runGitPullOrigin(profile.Branch); err != nil{
			return fmt.Errorf("error pull repository url from branch %s\nError: %v",profile.Branch,err)
		}

		profile.Bootstrap.Enabled = false

		if err := UpdateRegisterFile(*profile); err != nil{
			return err
		}
	}
	
	// 3. Check if .env file exists, otherwise generate using env_vars
	if profile.EnvFile != "" {
		if _, err := os.Stat(profile.EnvFile); os.IsNotExist(err) && len(profile.EnvVars) > 0 {
			utils.PrintInfo("env file not found. Generating from env_vars...")
			err := createEnvFile(profile.EnvFile, profile.EnvVars)
			if err != nil {
				return fmt.Errorf("failed to create env file: %v", err)
			}
			utils.PrintSuccess("env file generated.")
		}
	}

	// 4. Run pre-scripts
	for _, script := range profile.PreScripts {
		utils.PrintInfo(fmt.Sprintf("Running pre-script: %s", script))
		if err := runScript(script); err != nil {
			return fmt.Errorf("error running pre-script '%s': %v", script, err)
		}
	}
    fmt.Println("")

	// 5. Docker compose (optional)
	if profile.Docker.Enabled && profile.Docker.Up {
		utils.PrintInfo("Bringing up Docker Compose...\n")
		for _, file := range profile.Docker.ComposeFiles {
			utils.PrintInfo(fmt.Sprintf("📄 Using compose file: %s", file))
		}
		err := runDockerCompose(profile.Docker.ComposeFiles, profile.EnvFile)
		if err != nil {
			return fmt.Errorf("error running docker-compose: %v", err)
		}
		utils.PrintSuccess("Docker Compose is up.\n")
	}

	// 6. Main scripts
	for _, script := range profile.Scripts {
		utils.PrintInfo(fmt.Sprintf("Running script: %s", script))
		if err := runScript(script); err != nil {
			return fmt.Errorf("error running script '%s': %v", script, err)
		}
		utils.PrintSuccess(fmt.Sprintf("%s done", script))
	}

	// 7. Post scripts
	for _, script := range profile.PostScripts {
		utils.PrintInfo(fmt.Sprintf("🧹 Running post-script: %s", script))
		if err := runScript(script); err != nil {
			return fmt.Errorf("error running post-script '%s': %v", script, err)
		}
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

func runGitPullOrigin(branch string) error {
	cmd := exec.Command("git", "pull","origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDockerCompose(composeFiles []string, envFile string) error {
	args := []string{"compose"}

	for _, file := range composeFiles {
		args = append(args, "-f", file)
	}

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

func createEnvFile(filePath string, envVars map[string]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for k, v := range envVars {
		_, err := fmt.Fprintf(file, "%s=%s\n", k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func getGitRepoNameFromURL(repo string) (string, error) {
	url := strings.TrimSpace(repo)
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid origin URL: %s", url)
	}

	return parts[len(parts)-1], nil
}

func createDirNameFromRepo(repo string)(string, error){
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	repoName, err := getGitRepoNameFromURL(repo)
	if err != nil {
		return "", err
	}

	return cwd+"/"+repoName, nil;
}