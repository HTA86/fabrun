package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var inputValue string
	flag.StringVar(&inputValue, "var", "", "input value to replace {{input}} placeholder in command")
	flag.Parse()

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
