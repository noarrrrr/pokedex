package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

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
	callback    func() error
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
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func help() error {
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
var mapNext string = "https://pokeapi.co/api/v2/location-area"

func commandMap() error {
	if mapNext == "" {
		fmt.Println("No more locations to display")
		return nil
	}
	res, err := http.Get(mapNext)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer res.Body.Close()

	var data struct {
		Count    int
		Next     *string
		Previous *string
		Results  []struct {
			Name string
			Url  string
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	json.Unmarshal(body, &data)
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

func commandMapB() error {
	if mapPrevious == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		mapNext = mapPrevious
		commandMap()
	}
	return nil
}
