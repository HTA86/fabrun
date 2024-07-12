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
	commandDirs, err := getCommandDirectories()
	if err != nil {
		logError("Error reading commands directory:", err)
		return
	}

	fmt.Println("Available commands:")
	for _, dir := range commandDirs {
		if dir.IsDir() {
			fmt.Println(" -", dir.Name())
		}
	}
}

func getCommandDirectories() ([]os.DirEntry, error) {
	commandsDir := getCommandsDirectory()
	return os.ReadDir(commandsDir)
}

func getCommandsDirectory() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "fabrun", "commands")
}

func getCommandFilePath(commandID string) string {
	return filepath.Join(getCommandsDirectory(), commandID, "command.md")
}

func readCommandFromFile(filePath string) (string, error) {
	file, err := openFile(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return readCommand(file)
}

func openFile(filePath string) (*os.File, error) {
	return os.Open(filePath)
}

func readCommand(file *os.File) (string, error) {
	scanner := bufio.NewScanner(file)
	var commandStringBuilder strings.Builder
	for scanner.Scan() {
		commandStringBuilder.WriteString(scanner.Text())
	}

	return commandStringBuilder.String(), scanner.Err()
}

func executeCommand(command string) {
	cmd := buildCommand(command)
	runCommand(cmd)
}

func buildCommand(command string) *exec.Cmd {
	return exec.Command("bash", "-c", command)
}

func runCommand(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logError("Command execution failed with error:", err)
	}
}

func logError(message string, err error) {
	log.Println(message, err)
}
