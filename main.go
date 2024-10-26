package main

import (
	"bufio"
	"fmt"
	"os"

	pokeapi "github.com/joshlaudone/pokedex-repl/internal/pokeAPI"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config) error
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
	}
}

func commandHelp(cfg *pokeapi.Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	cliCommandMap := constructCliCommandMap()
	for _, command := range cliCommandMap {
		fmt.Printf("\t%s: %s\n", command.name, command.description)
	}

	fmt.Println()
	return nil
}

func commandExit(cfg *pokeapi.Config) error {
	fmt.Println("closed the pokedex")
	os.Exit(0)
	return nil
}

func main() {
	cliCommandMap := constructCliCommandMap()

	cfg := pokeapi.InitConfig()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex -> ")
		scanner.Scan()

		command, ok := cliCommandMap[scanner.Text()]
		if !ok {
			fmt.Printf("Invalid command: %s\n", scanner.Text())
			continue
		}

		err := command.callback(cfg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
