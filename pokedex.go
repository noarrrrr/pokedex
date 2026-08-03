package main

import (
	"fmt"
)

var totalPokemon = 1351

type pokedex map[int]pokemonStruct

var dex = make(pokedex)

func pokemonCaught(p pokemonStruct) {
	if _, ok := dex[p.ID]; ok {
		fmt.Printf("You already have %s\n", p.Name)
		return
	}
	fmt.Printf("%v was caught!\n", p.Name)
	dex[p.ID] = p
	fmt.Printf("You now have %d/%d pokemon\n", len(dex), totalPokemon)
}
