package main

import (
	"fmt"
	"strings"
)

var totalPokemon = 1351

type pokedex map[string]pokemonStruct

var dex = make(pokedex)

func parseName(args []string) string {
	var name string
	if len(args) > 1 {
		for i, arg := range args {
			args[i] = strings.Trim(arg, ".:")
		}
		name = strings.Join(args, "-")
	} else if len(args) < 1 {
		fmt.Println("please add the name of the pokemon")
		return ""
	} else {
		name = args[0]
	}
	return name
}

func pokemonCaught(p pokemonStruct) {
	if _, ok := dex[p.Name]; ok {
		fmt.Printf("You already have %s\n", p.Name)
		return
	}
	fmt.Printf("%v was caught!\n", p.Name)
	dex[p.Name] = p
	fmt.Printf("You now have %d/%d pokemon\n", len(dex), totalPokemon)
}
