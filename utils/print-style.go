package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PrintInfo(msg string) {
	fmt.Printf("\033[1;34m%s\033[0m\n", msg) // Blue
}

func PrintSuccess(msg string) {
	fmt.Printf("\033[1;32m%s\033[0m\n", msg) // Green
}

func PrintError(msg string) {
	fmt.Printf("\033[1;31m%s\033[0m\n", msg) // Red
}

func PrintWarning(msg string) {
	fmt.Printf("\033[1;33m%s\033[0m\n", msg) // Yellow
}

func PrintBoxedWarning(message string) {
	border := "╔" + strings.Repeat("═", 50) + "╗"
	footer := "╚" + strings.Repeat("═", 50) + "╝"

	fmt.Println("\033[1;33m" + border + "\033[0m") // Yellow border
	for _, line := range strings.Split(message, "\n") {
		fmt.Printf("\033[1;33m║\033[0m %-48s \033[1;33m║\033[0m\n", line)
	}
	fmt.Println("\033[1;33m" + footer + "\033[0m")
}

func AskConfirmation(question string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + " (y/n): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = input[:len(input)-1]

	if input == "y" || input == "Y" {
		return true, nil
	}
	return false, nil
}