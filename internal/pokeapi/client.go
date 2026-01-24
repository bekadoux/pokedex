package pokeapi

import (
	"github.com/bekadoux/pokedex/internal/pokecache"
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
	cache      pokecache.Cache
	Pokedex    Pokedex
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache:   pokecache.NewCache(cacheInterval),
		Pokedex: NewPokedex(),
	}
}
