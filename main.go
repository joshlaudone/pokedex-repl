package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	pokeapi "github.com/joshlaudone/pokedex-repl/internal/pokeAPI"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config, []string) error
}

func constructCliCommandMap() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 locations",
			callback:    pokeapi.GetNextLocations,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations",
			callback:    pokeapi.GetPrevLocations,
		},
		"explore": {
			name:        "explore",
			description: "Display the pokemon at a location",
			callback:    pokeapi.GetPokemonAtLocation,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch the specified Pokemon",
			callback:    pokeapi.TryToCatchPokemon,
		},
		"inspect": {
			name:        "inspect",
			description: "View info about the specified Pokemon",
			callback:    pokeapi.InspectPokemon,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View all captured pokemon",
			callback:    commandPokedex,
		},
	}
}

const helpSpaces = 10

func commandHelp(cfg *pokeapi.Config, params []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	cliCommandMap := constructCliCommandMap()
	for _, command := range cliCommandMap {
		spaces := strings.Repeat(" ", helpSpaces-len(command.name))
		fmt.Printf("\t%s:%s%s\n", command.name, spaces, command.description)
	}

	fmt.Println()
	return nil
}

func commandExit(cfg *pokeapi.Config, params []string) error {
	fmt.Println("closed the pokedex")
	os.Exit(0)
	return nil
}

func commandPokedex(cfg *pokeapi.Config, params []string) error {
	fmt.Println("Your Pokedex:")
	for pokemonName := range cfg.Pokedex {
		fmt.Printf("  - %s\n", pokemonName)
	}

	return nil
}

func main() {
	cliCommandMap := constructCliCommandMap()

	cfg := pokeapi.InitConfig()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex -> ")
		scanner.Scan()

		words := strings.Fields(scanner.Text())

		// Call help if the user enters an empty command
		if len(words) == 0 {
			words = append(words, "help")
		}

		command, ok := cliCommandMap[words[0]]
		if !ok {
			fmt.Printf("Invalid command: %s\n", scanner.Text())
			continue
		}

		err := command.callback(cfg, words[1:])
		if err != nil {
			fmt.Println(err)
		}
	}
}
