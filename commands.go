package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ListCommands() {
	dirs, err := GetCommandDirectories()
	if err != nil {
		LogError("Error reading commands directory:", err)
		return
	}

	fmt.Println("Available commands:")
	for _, dir := range dirs {
		if dir.IsDir() {
			fmt.Println(" -", dir.Name())
		}
	}
}

func ReadCommandFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var cmdBuilder strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cmdBuilder.WriteString(scanner.Text())
	}

	return cmdBuilder.String(), scanner.Err()
}

func ExecuteCommand(cmd string) {
	c := buildCommand(cmd)
	runCommand(c)
}

func buildCommand(cmd string) *exec.Cmd {
	return exec.Command("bash", "-c", cmd)
}

func runCommand(c *exec.Cmd) {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		LogError("Command execution failed with error:", err)
	}
}
