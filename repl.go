package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/noarrrrr/pokedex/internal/pokecache"
)

var cache = pokecache.CreateCache(time.Minute)

func cleanInput(text string) []string {
	var result []string
	splitted := strings.Split(strings.ToLower(text), " ")
	for i := range splitted {
		if splitted[i] != "" {
			result = append(result, splitted[i])
		}
	}
	return result
}

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

func getCmds() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "List available commands",
			callback:    help,
		},
		"map": {
			name:        "map",
			description: "List pokemon locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "List previous pokemon locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "List pokemon at a given location",
			callback:    commandExplore,
		},
	}
}

func commandExit(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func help(args []string) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:
 `)
	cmds := getCmds()
	for cmd := range cmds {
		cmdstruct := cmds[cmd]
		fmt.Printf("%v: %v\n", cmdstruct.name, cmdstruct.description)
	}
	fmt.Println("")
	return nil
}

var mapPrevious string
var mapNext string = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"

func commandMap(args []string) error {
	if mapNext == "" {
		fmt.Println("No more locations to display")
		return nil
	}

	bytes, err := getResBytes(mapNext)
	if err != nil {
		fmt.Println(err)
		return err
	}

	data := unmarshalMap(bytes)

	if data.Previous == nil {
		mapPrevious = ""
	} else {
		mapPrevious = *data.Previous
	}

	if data.Next == nil {
		mapNext = ""
	} else {
		mapNext = *data.Next
	}

	for location := range data.Results {
		fmt.Println(data.Results[location].Name)
	}

	return nil
}

func commandMapB(args []string) error {
	if mapPrevious == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		mapNext = mapPrevious
		commandMap([]string{})
	}
	return nil
}

func commandExplore(args []string) error {
	if len(args) < 1 {
		fmt.Println("Please add a location to explore")
		return nil
	} else if len(args) > 1 {
		fmt.Println("Too many locations provided")
		return nil
	}
	url := "https://pokeapi.co/api/v2/location-area/" + args[0]

	bytes, err := getResBytes(url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	data := unmarshalExplore(bytes)

	for _, encounter := range data.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}
	return nil
}
