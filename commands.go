package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	homeDir         string
	commandBasePath string
)

func init() {
	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	commandBasePath = filepath.Join(homeDir, ".config", "fabrun", "commands")
}

func executeCommand(command string) error {
	// Create a new command execution
	cmd := exec.Command("sh", "-c", command)

	// Set standard output and standard error to the current process's standard output and standard error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	return cmd.Run()
}

func readCommandFile(commandName string) (string, error) {
	filePath := filepath.Join(commandBasePath, commandName, "command.md")
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var command strings.Builder
	for scanner.Scan() {
		command.WriteString(scanner.Text())
		command.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.TrimSpace(command.String()), nil
}

func getAllCommandNames() ([]string, error) {
	entries, err := os.ReadDir(commandBasePath)
	if err != nil {
		return nil, err
	}

	var commandNames []string
	for _, entry := range entries {
		if entry.IsDir() {
			commandNames = append(commandNames, entry.Name())
		}
	}

	return commandNames, nil
}
