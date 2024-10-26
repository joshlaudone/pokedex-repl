package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	NextLocationArea *string
	PrevLocationArea *string
}

func InitConfig() *Config {
	nextLoc := "https://pokeapi.co/api/v2/location-area/"
	return &Config{
		NextLocationArea: &nextLoc,
		PrevLocationArea: nil,
	}
}

func GetNextLocations(cfg *Config) error {
	if cfg.NextLocationArea == nil {
		return fmt.Errorf("no more locations to show")
	}

	err := getLocationArea(cfg, cfg.NextLocationArea)
	return err
}

func GetPrevLocations(cfg *Config) error {
	if cfg.PrevLocationArea == nil {
		return fmt.Errorf("no previous locations to show")
	}

	err := getLocationArea(cfg, cfg.PrevLocationArea)
	return err
}

func getLocationArea(cfg *Config, url *string) error {
	resp, err := http.Get(*url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var locations LocationArea
	if err := json.Unmarshal(body, &locations); err != nil {
		return err
	}

	for _, result := range locations.Results {
		fmt.Println(result.Name)
	}

	cfg.NextLocationArea = locations.Next
	cfg.PrevLocationArea = locations.Previous

	return nil
}
