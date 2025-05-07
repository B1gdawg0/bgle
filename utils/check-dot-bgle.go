package utils

import (
	"os"
	"path/filepath"
)

// func EnsureBgleDirectory() error {
//     dir := ".bgle"
//     if _, err := os.Stat(dir); os.IsNotExist(err) {
//         fmt.Println("Creating .bgle directory...")
//         return os.Mkdir(dir, os.ModePerm)
//     }
//     return nil
// }

func EnsureBgleDirectoryAt(path string) error {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return os.MkdirAll(path, 0755)
    }
    return nil
}

func GetBgleDir() (string, error) {
    if custom := os.Getenv("BGLE_HOME"); custom != "" {
        return custom, nil
    }
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(home, ".bgle"), nil
}