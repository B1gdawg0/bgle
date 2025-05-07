package utils

import (
	"fmt"
	"strings"
)

func ParseNameProfile(input string) (string, string, error) {
    parts := strings.Split(input, ":")

    if len(parts) > 2 {
        return "", "", fmt.Errorf("invalid input: more than one ':' found")
    }

    name := parts[0]
    profile := "default" // default profile

    if len(parts) == 2 {
        profile = parts[1]
    }

    return name, profile, nil
}