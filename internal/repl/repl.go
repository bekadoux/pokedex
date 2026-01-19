package repl

import (
	"bufio"
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

		if len(clean) == 0 {
			continue
		}

		cmdName := clean[0]
		calledCmd, ok := cmdRegistry[cmdName]
		if !ok {
			fmt.Printf("Unknown command: '%s'\n", clean[0])
			continue
		}

		if err := calledCmd.callback(cfg); err != nil {
			err = fmt.Errorf("command failed: %w", err)
			fmt.Println(err)
			continue
		}
	}
}
