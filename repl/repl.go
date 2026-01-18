package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	clean := strings.Fields(lowercase)

	return clean
}

func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
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

		if err := calledCmd.callback(); err != nil {
			err = fmt.Errorf("command failed: %w", err)
			fmt.Println(err)
			continue
		}
	}
}
