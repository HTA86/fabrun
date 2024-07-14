package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var inputValue string
	var listCommands bool

	flag.StringVar(&inputValue, "var", "", "input value to replace {{input}} placeholder in command")
	flag.BoolVar(&listCommands, "list", false, "list all available commands")
	flag.BoolVar(&listCommands, "l", false, "list all available commands (shorthand)")
	flag.Parse()

	if listCommands {
		commandNames, err := getAllCommandNames()
		if err != nil {
			fmt.Printf("Error retrieving command names: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Available commands:")
		for _, name := range commandNames {
			fmt.Println(name)
		}
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Println("Usage: fabrun <command name> [--var <input>]")
		os.Exit(1)
	}

	commandName := flag.Arg(0)
	command, err := readCommandFile(commandName)
	if err != nil {
		fmt.Printf("Error reading command from file: %v\n", err)
		os.Exit(1)
	}

	commandRequireInput := strings.Contains(command, "{{input}}")

	if commandRequireInput {
		if inputValue == "" {
			fmt.Println("Usage: fabrun <command name> --var <input>")
			os.Exit(1)
		}
		command = strings.ReplaceAll(command, "{{input}}", inputValue)
	} else {
		fmt.Println("Ignoring --var <input> as it is not required for this command. Running without it.")
	}

	if err := executeCommand(command); err != nil {
		os.Exit(1)
	}

}
