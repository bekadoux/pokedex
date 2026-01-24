package repl

import (
	"errors"
	"fmt"
	"os"

	"github.com/bekadoux/pokedex/internal/pokeapi"
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
			description: "Exit the Pokédex",
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
			description: "Displays 20 location areas in the Pokémon world (each subsequent call displays the next 20 locations)",
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
		"catch": {
			name:        "catch",
			description: "Catch a Pokémon!",
			minArgs:     1,
			maxArgs:     1,
			callback: func(cfg *config, args []string) error {
				return cmdCatch(cfg, args)
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

func cmdCatch(cfg *config, args []string) error {
	name := args[0]
	response, err := cfg.client.GetPokemonDetails(name)
	if err != nil {
		return fmt.Errorf("error getting pokemon details: %w", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	// Max possible base_experience, set arbitrarily for simplicity
	maxBaseExp := 350
	catchSuccess := pokeapi.AttemptCatchPokemon(response.BaseExperience, maxBaseExp, 0.05, 0.95)
	if catchSuccess {
		pokemon := response.ToPokemon()
		err = cfg.client.Pokedex.AddPokemon(pokemon)
		if errors.Is(err, pokeapi.ErrAddDuplicatePokemon) {
			fmt.Printf("You already have a %s!\n", name)
		} else if err != nil {
			return fmt.Errorf("error adding pokemon to pokedex: %s", err)
		} else {
			fmt.Printf("%s was caught!\n", name)
		}
	} else {
		fmt.Printf("%s escaped!\n", name)
	}

	return nil
}
