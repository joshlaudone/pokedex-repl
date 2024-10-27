package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/joshlaudone/pokedex-repl/internal/pokecache"
)

const locationAreaUrl = "https://pokeapi.co/api/v2/location-area/"

func InitConfig() *Config {
	nextLoc := locationAreaUrl
	cache := pokecache.New(1 * time.Minute)
	return &Config{
		NextLocationArea: &nextLoc,
		PrevLocationArea: nil,
		Cache:            cache,
	}
}

func GetNextLocations(cfg *Config, params []string) error {
	if cfg.NextLocationArea == nil {
		return fmt.Errorf("no more locations to show")
	}

	err := getLocationArea(cfg, *cfg.NextLocationArea)
	return err
}

func GetPrevLocations(cfg *Config, params []string) error {
	if cfg.PrevLocationArea == nil {
		return fmt.Errorf("no previous locations to show")
	}

	err := getLocationArea(cfg, *cfg.PrevLocationArea)
	return err
}

func getLocationArea(cfg *Config, url string) error {
	data, err := cachedGet(cfg, url)
	if err != nil {
		return err
	}

	var locations LocationArea
	if err := json.Unmarshal(data, &locations); err != nil {
		return err
	}

	for _, result := range locations.Results {
		fmt.Println(result.Name)
	}

	cfg.NextLocationArea = locations.Next
	cfg.PrevLocationArea = locations.Previous

	return nil
}

func GetPokemonAtLocation(cfg *Config, params []string) error {
	if len(params) < 1 {
		return fmt.Errorf("pass in the location you would like to explore")
	}

	location := params[0]
	url := locationAreaUrl + location

	data, err := cachedGet(cfg, url)
	if err != nil {
		return err
	}

	var area ExploredArea
	if err := json.Unmarshal(data, &area); err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", location)
	fmt.Println("Found Pokemon:")
	for _, encounter := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func cachedGet(cfg *Config, url string) ([]byte, error) {
	data, found := cfg.Cache.Get(url)
	if !found {
		resp, err := http.Get(url)
		if err != nil {
			return []byte{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode > 299 {
			return []byte{}, fmt.Errorf("received %d http status code", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, err
		}

		cfg.Cache.Add(url, body)

		data = body
	}

	return data, nil
}
