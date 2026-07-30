package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getResBytes(url string) ([]byte, error) {
	body, ok := cache.Get(url)

	if !ok {
		fmt.Println("fetching data from server")
		response, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return []byte{}, err
		}
		body, err = io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return []byte{}, err
		}
		cache.Add(url, body)
		defer response.Body.Close()
	} else {
		fmt.Println("successfully pulled data from cache")
	}
	return body, nil
}

func unmarshalMap(bytes []byte) mapStruct {
	var data mapStruct
	json.Unmarshal(bytes, &data)
	return data
}

func unmarshalExplore(bytes []byte) exploreStruct {
	var data exploreStruct
	json.Unmarshal(bytes, &data)
	return data
}

type mapStruct struct {
	Count    int
	Next     *string
	Previous *string
	Results  []struct {
		Name string
		Url  string
	}
}

type exploreStruct struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel       int `json:"min_level"`
				PokemonDetails any `json:"pokemon_details"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
