package pokeapi

import (
	"encoding/json"
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

func (c *Client) GetLocationAreas(pageURL string) (LocationAreaResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != "" {
		url = pageURL
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

	var locResponse LocationAreaResponse
	err = json.Unmarshal(data, &locResponse)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error during Unmarshal: %w", err)
	}

	return locResponse, nil
}
