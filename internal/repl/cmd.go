package repl

import (
	"fmt"
	"github.com/bekadoux/pokedex/internal/pokeapi"
	"os"
)

type cmd struct {
	name        string
	description string
	minArgs     int
	maxArgs     int
	callback    func(*config, []string) error
}

type config struct {
	client              pokeapi.Client
	NextLocationURL     string
	PreviousLocationURL string
}

var cmdRegistry map[string]cmd

func init() {
	cmdRegistry = map[string]cmd{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			minArgs:     0,
			maxArgs:     0,
			callback: func(cfg *config, args []string) error {
				return cmdExit()
			},
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			minArgs:     0,
			maxArgs:     0,
			callback: func(cfg *config, args []string) error {
				return cmdHelp()
			},
		},
		"map": {
			name:        "map",
			description: "Displays 20 location areas in the Pokemon world (each subsequent call displays the next 20 locations)",
			minArgs:     0,
			maxArgs:     0,
			callback: func(cfg *config, args []string) error {
				return cmdMapFwd(cfg)
			},
		},
		"mapb": {
			name:        "mapb",
			description: "Return to previous map page",
			minArgs:     0,
			maxArgs:     0,
			callback: func(cfg *config, args []string) error {
				return cmdMapBack(cfg)
			},
		},
		"explore": {
			name:        "explore",
			description: "Explore location area",
			minArgs:     1,
			maxArgs:     1,
			callback: func(cfg *config, args []string) error {
				return cmdExplore(cfg, args)
			},
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

func cmdMapFwd(cfg *config) error {
	locResponse, err := cfg.client.GetLocationAreas(cfg.NextLocationURL)
	if err != nil {
		return fmt.Errorf("error getting location areas: %w", err)
	}

	cfg.NextLocationURL = locResponse.Next
	cfg.PreviousLocationURL = locResponse.Previous

	for _, area := range locResponse.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func cmdMapBack(cfg *config) error {
	if len(cfg.PreviousLocationURL) == 0 {
		fmt.Println("You're on the first page.")
		return nil
	}

	response, err := cfg.client.GetLocationAreas(cfg.PreviousLocationURL)
	if err != nil {
		return fmt.Errorf("error getting location areas: %w", err)
	}

	cfg.NextLocationURL = response.Next
	cfg.PreviousLocationURL = response.Previous

	for _, area := range response.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func cmdExplore(cfg *config, args []string) error {
	location := args[0]
	response, err := cfg.client.GetLocationDetails(args[0])
	if err != nil {
		return fmt.Errorf("error getting location areas: %w", err)
	}

	fmt.Printf("Exploring %s...\n", location)
	if len(response.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")
	}
	for _, encounter := range response.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil
}
