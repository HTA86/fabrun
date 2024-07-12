package main

import (
	"log"
	"os"
	"path/filepath"
)

func GetCommandsDirectory() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "fabrun", "commands")
}

func GetCommandFilePath(commandID string) string {
	return filepath.Join(GetCommandsDirectory(), commandID, "command.md")
}

func GetCommandDirectories() ([]os.DirEntry, error) {
	commandsDir := GetCommandsDirectory()
	return os.ReadDir(commandsDir)
}

func LogError(message string, err error) {
	log.Println(message, err)
}
