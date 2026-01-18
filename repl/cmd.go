package repl

import (
	"fmt"
	"os"
)

type cmd struct {
	name        string
	description string
	callback    func() error
}

var cmdRegistry map[string]cmd

func init() {
	cmdRegistry = map[string]cmd{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    cmdExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    cmdHelp,
		},
	}
}

func cmdExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cmdHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	for k, v := range cmdRegistry {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}
