// Package embed provides embedded assets for TAW.
package embed

import (
	"embed"
)

//go:embed assets/*
var Assets embed.FS

// GetHelp returns the help content.
func GetHelp() (string, error) {
	data, err := Assets.ReadFile("assets/HELP.md")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
