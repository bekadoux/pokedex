package repl

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bekadoux/pokedex/internal/pokeapi"
)

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	clean := strings.Fields(lowercase)

	return clean
}

func dispatch(cfg *config, input []string) error {
	if len(input) == 0 {
		return errors.New("empty input")
	}

	cmdName := input[0]
	args := input[1:]

	calledCmd, ok := cmdRegistry[cmdName]
	if !ok {
		return fmt.Errorf("unknown command: '%s'", cmdName)
	}

	if len(args) < calledCmd.minArgs || len(args) > calledCmd.maxArgs {
		if calledCmd.minArgs == calledCmd.maxArgs {
			return fmt.Errorf("%s expects %d arguments, got %d", cmdName, calledCmd.minArgs, len(args))
		}
		return fmt.Errorf("%s expects between %d and %d arguments, got %d", cmdName, calledCmd.minArgs, calledCmd.maxArgs, len(args))
	}

	return calledCmd.callback(cfg, args)
}

func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		client: pokeapi.NewClient(10 * time.Second),
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		input := scanner.Text()
		clean := cleanInput(input)

		err := dispatch(cfg, clean)
		if err != nil {
			fmt.Printf("Command error: ")
			fmt.Println(err)
		}
	}
}
