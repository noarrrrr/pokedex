package main

import (
	"fmt"
	"math/rand"
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
		"catch": {
			name:        "catch",
			description: "Attempt to catch a given pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "List information about a pokemon in your dex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List pokemon in your pokedex",
			callback:    commandDex,
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

func commandCatch(args []string) error {
	name := parseName(args)
	if name == "" {
		return nil
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + name

	bytes, err := getResBytes(url)
	if err != nil {
		fmt.Println("pokemon not found")
		return err
	}

	pokemon := unmarshalPokemon(bytes)

	if pokemon.BaseExperience == 0 {
		fmt.Println("pokemon not found")
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	catchChance := int(5000 / pokemon.BaseExperience)
	fmt.Printf("%v%% chance to catch\n", catchChance)

	if catchChance >= 100 {
		pokemonCaught(pokemon)
	} else {
		if catchChance > rand.Intn(99) {
			pokemonCaught(pokemon)
		} else {
			fmt.Printf("%v escaped!\n", pokemon.Name)
			return nil
		}
	}

	return nil
}

func commandInspect(args []string) error {
	name := parseName(args)
	if name == "" {
		return nil
	}

	p, ok := dex[name]

	if !ok {
		fmt.Println("You do not have this pokemon")
		return nil
	}
	fmt.Printf(`Name: %v
Height: %v
Weight: %v
Stats:
`, p.Name, p.Height, p.Weight)

	for _, stat := range p.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}
	return nil
}

func commandDex(args []string) error {
	if len(dex) == 0 {
		fmt.Println("You don't have any pokemon")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, p := range dex {
		fmt.Printf("  -%s\n", p.Name)
	}
	return nil
}
