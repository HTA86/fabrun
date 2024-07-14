package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	commandFilePath := filepath.Join(homeDir, ".config", "fabrun", "commands", commandName, "command.md")

	file, err := os.Open(commandFilePath)
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
