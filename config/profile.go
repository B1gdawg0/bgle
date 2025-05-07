package config

import (
	"bgle/utils"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Profile struct {
    Project string   `yaml:"project"`
    Profile string   `yaml:"profile"`
    Dir     string   `yaml:"directory"`
    Branch  string   `yaml:"branch,omitempty"`
    Docker  Docker   `yaml:"docker"`
    EnvFile string   `yaml:"env_file,omitempty"`
    Scripts []string `yaml:"scripts"`
}

type Docker struct {
    Enabled     bool   `yaml:"enabled"`
    ComposeFile string `yaml:"compose_file"`
    Up          bool   `yaml:"up"`
    Build       bool   `yaml:"build"`
}

// Function to scan for docker-compose files in the current directory
func selectDockerComposeFile() string {
    // Scan the directory for docker-compose*.yml files
    files, _ := filepath.Glob("docker-compose*.yml")
    if len(files) == 0 {
        fmt.Println("No docker-compose files found.")
        return ""
    }

    if len(files) == 1 {
        fmt.Printf("Using Docker Compose file: %s\n", files[0])
        return files[0]
    }

    fmt.Println("Select a Docker Compose file:")
    for i, file := range files {
        fmt.Printf("%d) %s\n", i+1, file)
    }

    var choice int
    fmt.Print("Choose one (number): ")
    fmt.Scanln(&choice)

    if choice < 1 || choice > len(files) {
        fmt.Println("Invalid choice. Defaulting to first.")
        return files[0]
    }

    return files[choice-1]
}

func InteractiveProfileCreation(name string, profile string) Profile {
    utils.PrintInfo("Starting profile creation...")
    reader := bufio.NewReader(os.Stdin)

    // Default to current directory for the project
    dir, _ := os.Getwd()
    // fmt.Print("Project directory (default: current): ")
    // dir, _ := reader.ReadString('\n')
    // dir = strings.TrimSpace(dir)
    // if dir == "" {
    //     dir, _ = os.Getwd()
    // }

    // Git branch (optional)
    fmt.Print("Git branch (optional): ")
    branch, _ := reader.ReadString('\n')
    branch = strings.TrimSpace(branch)

    // Ask if Docker should be used
    fmt.Print("Use Docker? (y/n): ")
    dockerInput, _ := reader.ReadString('\n')
    docker := strings.TrimSpace(dockerInput) == "y"

    var dockerConfig Docker
    if docker {
        // Allow user to select the Docker Compose file
        file := selectDockerComposeFile()
        if file == "" {
            fmt.Println("No Docker Compose file selected. Proceeding without Docker.")
        }

        // Ask for docker-compose up and build
        fmt.Print("Run docker-compose up -d? (y/n): ")
        upInput, _ := reader.ReadString('\n')
        up := strings.TrimSpace(upInput) == "y"

        fmt.Print("Run docker build? (y/n): ")
        buildInput, _ := reader.ReadString('\n')
        build := strings.TrimSpace(buildInput) == "y"

        dockerConfig = Docker{
            Enabled:     true,
            ComposeFile: file,
            Up:          up,
            Build:       build,
        }
    }

    // Ask for path to .env file (optional)
    fmt.Print("Path to .env file (optional): ")
    envFile, _ := reader.ReadString('\n')
    envFile = strings.TrimSpace(envFile)

    // Ask for scripts
    fmt.Println("Enter startup scripts one by one. Type 'exit' to finish:")
    var scripts []string
    for {
        fmt.Print("> ")
        script, _ := reader.ReadString('\n')
        script = strings.TrimSpace(script)
        if script == "exit" {
            break
        }
        scripts = append(scripts, script)
    }

    profileConfig := Profile{
        Project: name,
        Profile: profile,
        Dir:     dir,
        Branch:  branch,
        Docker:  dockerConfig,
        EnvFile: envFile,
        Scripts: scripts,
    }

    if err := SaveProfile(profileConfig); err != nil {
        utils.PrintError("Error saving profile:"+err.Error())
    } else {
        utils.PrintSuccess("Profile saved successfully!")
    }

    return profileConfig
}