package utils

import "fmt"

func PrintInfo(msg string) {
	fmt.Printf("\033[1;34m%s\033[0m\n", msg) // Blue
}

func PrintSuccess(msg string) {
	fmt.Printf("\033[1;32m%s\033[0m\n", msg) // Green
}

func PrintError(msg string) {
	fmt.Printf("\033[1;31m%s\033[0m\n", msg) // Red
}
