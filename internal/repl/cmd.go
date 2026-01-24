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
	pokedex             pokeapi.Pokedex
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
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokémon in your Pokédex",
			minArgs:     1,
			maxArgs:     1,
			callback: func(cfg *config, args []string) error {
				return cmdInspect(cfg, args)
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all Pokémon in your Pokédex",
			minArgs:     0,
			maxArgs:     0,
			callback: func(cfg *config, args []string) error {
				return cmdPokedex(cfg)
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
		err = cfg.pokedex.AddPokemon(pokemon)
		if errors.Is(err, pokeapi.ErrAddDuplicatePokemon) {
			fmt.Printf("You already have a %s!\n", name)
		} else if err != nil {
			return fmt.Errorf("error adding pokemon '%s' to pokedex: %w", name, err)
		} else {
			fmt.Printf("%s was caught!\n", name)
			fmt.Println("You may now inspect it with the inspect command.")
		}
	} else {
		fmt.Printf("%s escaped!\n", name)
	}

	return nil
}

func cmdInspect(cfg *config, args []string) error {
	name := args[0]
	pokemon, err := cfg.pokedex.GetPokemon(name)
	if errors.Is(err, pokeapi.ErrGetAbsentPokemon) {
		fmt.Printf("You haven't caught a %s yet!\n", name)
		return nil
	}
	if err != nil {
		return fmt.Errorf("error inspecting pokemon '%s': %w", name, err)
	}

	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", pokemon.Name, pokemon.Height, pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("\t- %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pType := range pokemon.Types {
		fmt.Printf("\t- %s\n", pType.Type.Name)
	}

	return nil
}

func cmdPokedex(cfg *config) error {
	allPokemon := cfg.pokedex.GetAllPokemon()
	if len(allPokemon) == 0 {
		fmt.Println("You have not caught any Pokémon yet!")
		return nil
	}

	fmt.Println("Your Pokédex:")
	for _, p := range allPokemon {
		fmt.Printf("\t- %s\n", p.Name)
	}

	return nil
}
