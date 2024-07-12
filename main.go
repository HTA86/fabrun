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
		fmt.Println("fabrun version", version)
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

	commandName := flag.Arg(0)
	commandPath := filepath.Join(os.Getenv("HOME"), ".config", "fabrun", "commands", commandName, "command.md")

	command, err := readCommand(commandPath)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	runCommand(command)
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

func listCommands() {
	commandsDir := filepath.Join(os.Getenv("HOME"), ".config", "fabrun", "commands")
	files, err := os.ReadDir(commandsDir)
	if err != nil {
		log.Println("Error reading commands directory:", err)
		return
	}

	fmt.Println("Available commands:")
	for _, file := range files {
		if file.IsDir() {
			fmt.Println(" -", file.Name())
		}
	}
}

func readCommand(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var commandString strings.Builder
	for scanner.Scan() {
		commandString.WriteString(scanner.Text())
	}

	return commandString.String(), scanner.Err()
}

func runCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println("Command execution failed with error:", err)
	}
}
