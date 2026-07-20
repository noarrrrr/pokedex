package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		dirtyRes := scanner.Text()
		res := cleanInput(dirtyRes)
		cmd := res[0]
		cmds := getCmds()
		cmdstruct, ok := cmds[cmd]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		cmdstruct.callback()
	}
}
