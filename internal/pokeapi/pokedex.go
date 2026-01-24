package pokeapi

import (
	"errors"
	"fmt"
	"math/rand"
)

var ErrAddDuplicatePokemon = errors.New("pokemon already present in pokedex")
var ErrGetAbsentPokemon = errors.New("pokemon is not in pokedex")

type Pokedex struct {
	CaughtPokemon map[string]Pokemon
}

func NewPokedex() Pokedex {
	return Pokedex{
		CaughtPokemon: make(map[string]Pokemon),
	}
}

func (p *Pokedex) AddPokemon(pokemon Pokemon) error {
	_, ok := p.CaughtPokemon[pokemon.Name]
	if ok {
		return ErrAddDuplicatePokemon
	}

	p.CaughtPokemon[pokemon.Name] = pokemon

	return nil
}

func (p *Pokedex) GetPokemon(name string) (Pokemon, error) {
	pokemon, ok := p.CaughtPokemon[name]
	if !ok {
		return Pokemon{}, ErrGetAbsentPokemon
	}

	return pokemon, nil
}

func AttemptCatchPokemon(baseExp, maxBaseExp int, minChance, maxChance float64) bool {
	chance := catchChance(baseExp, maxBaseExp, minChance, maxChance)
	fmt.Println(chance)
	return rand.Float64() <= chance
}

func catchChance(baseExp, maxBaseExp int, minChance, maxChance float64) float64 {
	if baseExp > maxBaseExp {
		baseExp = maxBaseExp
	}

	ratio := float64(baseExp) / float64(maxBaseExp)
	chance := 1.0 - ratio

	if chance < minChance {
		return minChance
	}
	if chance > maxChance {
		return maxChance
	}
	return chance
}
