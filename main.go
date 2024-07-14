package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: fabrun <subcommand>")
		os.Exit(1)
	}

	commandName := os.Args[1]
	command, err := readCommandFile(commandName)
	if err != nil {
		fmt.Printf("Error reading command from file: %v\n", err)
		os.Exit(1)
	}

	if err := executeCommand(command); err != nil {
		os.Exit(1)
	}

}
