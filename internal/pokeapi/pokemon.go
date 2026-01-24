package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Pokemon struct {
	Types  []PokemonType
	Stats  []PokemonStat
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
}

type PokemonType struct {
	Type NamedAPIResource `json:"type"`
	Slot int              `json:"slot"`
}

type PokemonTypePast struct {
	Generation NamedAPIResource `json:"generation"`
	Types      []PokemonType    `json:"types"`
}

type PokemonAbility struct {
	AbilityInfo NamedAPIResource `json:"ability"`
	Slot        int              `json:"slot"`
	IsHidden    bool             `json:"is_hidden"`
}

type PokemonAbilityPast struct {
	Generation NamedAPIResource `json:"generation"`
	Abilities  []PokemonAbility `json:"abilities"`
}

type PokemonCries struct {
	Latest string `json:"latest"`
	Legacy string `json:"legacy"`
}

type GenerationGameIndex struct {
	Version NamedAPIResource `json:"version"`
	Index   int              `json:"game_index"`
}

type PokemonHeldItem struct {
	Item           NamedAPIResource         `json:"item"`
	VersionDetails []PokemonHeldItemVersion `json:"version_details"`
}

type PokemonHeldItemVersion struct {
	Version NamedAPIResource `json:"version"`
	Rarity  int              `json:"rarity"`
}

type PokemonMove struct {
	MoveInfo            NamedAPIResource     `json:"move"`
	VersionGroupDetails []PokemonMoveVersion `json:"version_group_details"`
}

type PokemonMoveVersion struct {
	MoveLearnMethod NamedAPIResource `json:"move_learn_method"`
	VersionGroup    NamedAPIResource `json:"version_group"`
	LevelLearnedAt  int              `json:"level_learned_at"`
	Order           int              `json:"order"`
}

type PokemonStat struct {
	Stat     NamedAPIResource `json:"stat"`
	Effort   int              `json:"effort"`
	BaseStat int              `json:"base_stat"`
}

// `json:"sprites"` is ignored
type PokemonDetailsResponse struct {
	Species                NamedAPIResource      `json:"species"`
	Cries                  PokemonCries          `json:"cries"`
	Types                  []PokemonType         `json:"types"`
	PastTypes              []PokemonTypePast     `json:"past_types"`
	Moves                  []PokemonMove         `json:"moves"`
	Abilities              []PokemonAbility      `json:"abilities"`
	PastAbilities          []PokemonAbilityPast  `json:"past_abilities"`
	Stats                  []PokemonStat         `json:"stats"`
	HeldItems              []PokemonHeldItem     `json:"held_items"`
	Forms                  []NamedAPIResource    `json:"forms"`
	GameIndices            []GenerationGameIndex `json:"game_indices"`
	Name                   string                `json:"name"`
	LocationAreaEncounters string                `json:"location_area_encounters"`
	ID                     int                   `json:"id"`
	Height                 int                   `json:"height"`
	Weight                 int                   `json:"weight"`
	Order                  int                   `json:"order"`
	BaseExperience         int                   `json:"base_experience"`
	IsDefault              bool                  `json:"is_default"`
}

func (r *PokemonDetailsResponse) ToPokemon() Pokemon {
	return Pokemon{
		Types:  r.Types,
		Stats:  r.Stats,
		Name:   r.Name,
		Height: r.Height,
		Weight: r.Weight,
	}
}

func (c *Client) GetPokemonDetails(pokemonName string) (PokemonDetailsResponse, error) {
	if pokemonName == "" {
		return PokemonDetailsResponse{}, errors.New("no pokemon name provided")
	}

	url := baseURL + "/pokemon/" + pokemonName
	var pokemonResponse PokemonDetailsResponse

	if data, ok := c.cache.Get(url); ok {
		err := json.Unmarshal(data, &pokemonResponse)
		if err != nil {
			return PokemonDetailsResponse{}, fmt.Errorf("error during Unmarshal: %w", err)
		}

		return pokemonResponse, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokemonDetailsResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonDetailsResponse{}, fmt.Errorf("network error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return PokemonDetailsResponse{}, fmt.Errorf("unknown pokemon: '%s'", pokemonName)
	}
	if res.StatusCode > 299 {
		return PokemonDetailsResponse{}, fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokemonDetailsResponse{}, fmt.Errorf("error reading response: %w", err)
	}
	c.cache.Add(url, data)

	err = json.Unmarshal(data, &pokemonResponse)
	if err != nil {
		return PokemonDetailsResponse{}, fmt.Errorf("error during Unmarshal: %w", err)
	}

	return pokemonResponse, nil
}
