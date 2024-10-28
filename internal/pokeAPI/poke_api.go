package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/joshlaudone/pokedex-repl/internal/pokecache"
)

const locationAreaUrl = "https://pokeapi.co/api/v2/location-area/"
const pokeDataUrl = "https://pokeapi.co/api/v2/pokemon/"
const maxBaseXP = 350.0

func InitConfig() *Config {
	nextLoc := locationAreaUrl
	cache := pokecache.New(1 * time.Minute)
	return &Config{
		NextLocationArea: &nextLoc,
		PrevLocationArea: nil,
		Cache:            cache,
		Pokedex:          make(map[string]PokemonData),
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

func TryToCatchPokemon(cfg *Config, params []string) error {
	if len(params) < 1 {
		return fmt.Errorf("must pass in the name of the pokemon to catch")
	}

	pokemonToCatch := params[0]
	url := pokeDataUrl + pokemonToCatch

	data, err := cachedGet(cfg, url)
	if err != nil {
		return err
	}

	var pokeData PokemonData
	if err := json.Unmarshal(data, &pokeData); err != nil {
		return err
	}

	baseXP := float64(pokeData.BaseExperience)
	catch_chance := 1 - baseXP/maxBaseXP
	if catch_chance <= 0.0 {
		catch_chance = 0.01
	}

	fmt.Printf("Throwing a Pokeball at %s... \n", pokemonToCatch)

	if rand.Float64() > catch_chance {
		fmt.Printf("%s escaped!\n", pokemonToCatch)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemonToCatch)
	cfg.Pokedex[pokemonToCatch] = pokeData

	return nil
}

func InspectPokemon(cfg *Config, params []string) error {
	if len(params) < 1 {
		return fmt.Errorf("pass in the name of the pokemon to inspect")
	}

	inspectingPokemon := params[0]
	pokeData, found := cfg.Pokedex[inspectingPokemon]
	if !found {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name: %s\n", pokeData.Name)
	fmt.Printf("Height: %d\n", pokeData.Height)
	fmt.Printf("Weight: %d\n", pokeData.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokeData.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pokeType := range pokeData.Types {
		fmt.Printf("  - %s\n", pokeType.Type.Name)
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
