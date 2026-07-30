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
		args := res[1:]
		cmds := getCmds()
		cmdstruct, ok := cmds[cmd]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		cmdstruct.callback(args)
		err := scanner.Err()
		if err != nil {
			fmt.Println(err)
		}
	}

}
