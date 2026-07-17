package main

import (
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
