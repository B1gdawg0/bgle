package config

import (
	"bgle/models"
	"bgle/utils"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func SaveProfile(profile models.Profile) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    bgleDir := filepath.Join(homeDir, ".bgle")

    if err := utils.EnsureBgleDirectoryAt(bgleDir); err != nil {
        return err
    }

    filePath := filepath.Join(bgleDir, fmt.Sprintf("%s.%s.yaml", profile.Project, profile.Profile))
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := yaml.NewEncoder(file)
    encoder.SetIndent(2)
    return encoder.Encode(profile)
}
