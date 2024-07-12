package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	listFlag    bool
	versionFlag bool
	version     = "1.0.0"
)

func init() {
	flag.BoolVar(&listFlag, "list", false, "List all available commands")
	flag.BoolVar(&listFlag, "l", false, "List all available commands (short flag)")
	flag.BoolVar(&versionFlag, "version", false, "Show the program's version")
	flag.BoolVar(&versionFlag, "v", false, "Show the program's version (short flag)")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if versionFlag {
		printVersion()
		return
	}

	if listFlag {
		listCommands()
		return
	}

	if len(flag.Args()) < 1 {
		flag.Usage()
		return
	}

	commandID := flag.Arg(0)
	commandPath := getCommandFilePath(commandID)

	command, err := readCommandFromFile(commandPath)
	if err != nil {
		logError("Error reading command:", err)
		return
	}

	executeCommand(command)
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: fabrun [options] <command_name>\n")
	fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nExample usage:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  fabrun --list               # List all available commands\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  fabrun -l                   # List all available commands (short flag)\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  fabrun git_change_2weeks    # Run the git_change_2weeks command\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  fabrun --version            # Show the program's version\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  fabrun -v                   # Show the program's version (short flag)\n")
}

func printVersion() {
	fmt.Println("fabrun version", version)
}

func listCommands() {
	commandsDir := getCommandsDirectory()
	dirs, err := os.ReadDir(commandsDir)
	if err != nil {
		logError("Error reading commands directory:", err)
		return
	}

	fmt.Println("Available commands:")
	for _, dir := range dirs {
		if dir.IsDir() {
			fmt.Println(" -", dir.Name())
		}
	}
}

func getCommandsDirectory() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "fabrun", "commands")
}

func getCommandFilePath(commandID string) string {
	return filepath.Join(getCommandsDirectory(), commandID, "command.md")
}

func readCommandFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var commandBuilder strings.Builder
	for scanner.Scan() {
		commandBuilder.WriteString(scanner.Text())
	}

	return commandBuilder.String(), scanner.Err()
}

func executeCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logError("Command execution failed with error:", err)
	}
}

func logError(message string, err error) {
	log.Println(message, err)
}
