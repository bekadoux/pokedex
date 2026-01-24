package pokeapi

import "time"

const (
	baseURL       = "https://pokeapi.co/api/v2"
	cacheInterval = 1 * time.Minute
)

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Name struct {
	Language NamedAPIResource `json:"language"`
	Name     string           `json:"name"`
}
