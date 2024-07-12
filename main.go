package main

import (
	"flag"
	"fmt"
)

var (
	listFlag    bool
	versionFlag bool
	version     = "0.5.0"
)

func init() {
	flag.BoolVar(&listFlag, "list", false, "")
	flag.BoolVar(&listFlag, "l", false, "List all available commands")
	flag.BoolVar(&versionFlag, "version", false, "")
	flag.BoolVar(&versionFlag, "v", false, "Show the program's version")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if versionFlag {
		printVersion()
		return
	}

	if listFlag {
		ListCommands()
		return
	}

	if len(flag.Args()) < 1 {
		flag.Usage()
		return
	}

	commandID := flag.Arg(0)
	commandPath := GetCommandFilePath(commandID)

	cmd, err := ReadCommandFromFile(commandPath)
	if err != nil {
		LogError("Error reading command:", err)
		return
	}

	ExecuteCommand(cmd)
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: fabrun [options] <command_name>\n")
	fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  -l, --list     List all available commands\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  -v, --version  Show the program's version\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  -h, --help     Show this help message\n")
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
