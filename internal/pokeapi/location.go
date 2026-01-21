package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type LocationDetailsResponse struct {
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	GameIndex            int                   `json:"game_index"`
	ID                   int                   `json:"id"`
	Location             Location              `json:"location"`
	Name                 string                `json:"name"`
	Names                []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

func (c *Client) GetLocationAreas(pageURL string) (LocationAreaResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != "" {
		url = pageURL
	}

	if data, ok := c.cache.Get(url); ok {
		var locResponse LocationAreaResponse
		err := json.Unmarshal(data, &locResponse)
		if err != nil {
			return LocationAreaResponse{}, fmt.Errorf("error during Unmarshal: %w", err)
		}

		return locResponse, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("network error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return LocationAreaResponse{}, fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error reading response: %w", err)
	}
	c.cache.Add(url, data)

	var locResponse LocationAreaResponse
	err = json.Unmarshal(data, &locResponse)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error during Unmarshal: %w", err)
	}

	return locResponse, nil
}

func (c *Client) GetLocationDetails(location string) (LocationDetailsResponse, error) {
	url := baseURL + "/location-area/" + location
	if location == "" {
		return LocationDetailsResponse{}, errors.New("no location provided")
	}

	if data, ok := c.cache.Get(url); ok {
		var locDetailsResponse LocationDetailsResponse
		err := json.Unmarshal(data, &locDetailsResponse)
		if err != nil {
			return LocationDetailsResponse{}, fmt.Errorf("error during Unmarshal: %w", err)
		}

		return locDetailsResponse, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationDetailsResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return LocationDetailsResponse{}, fmt.Errorf("network error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return LocationDetailsResponse{}, fmt.Errorf("location not found")
	}
	if res.StatusCode > 299 {
		return LocationDetailsResponse{}, fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationDetailsResponse{}, fmt.Errorf("error reading response: %w", err)
	}
	c.cache.Add(url, data)

	var locDetailsResponse LocationDetailsResponse
	err = json.Unmarshal(data, &locDetailsResponse)
	if err != nil {
		return LocationDetailsResponse{}, fmt.Errorf("error during Unmarshal: %w", err)
	}

	return locDetailsResponse, nil
}
