package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	pokecache "github.com/zrtgzrtg/pokedexcli/internal"
)

func commandExplore(cptr *config, args []string) error {

	location := args[0]

	if cache == nil {
		cache = pokecache.NewCache(30 * time.Second)
	}

	fmt.Printf("Exploring %v...\n", location)
	fmt.Println("Found Pokemon:")

	url := "https://pokeapi.co/api/v2/location-area/" + location + "/"
	var stringResp []byte
	var printRes string

	val, ok := cache.Get(url)
	if ok {
		printRes = string(val)
		fmt.Println("cache hit")
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		stringResp, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var jsonResponse PokeLocationArea

		err = json.Unmarshal(stringResp, &jsonResponse)
		if err != nil {
			return err
		}

		foundPokemon := []string{}

		for _, encounter := range jsonResponse.PokemonEncounters {
			foundPokemon = append(foundPokemon, fmt.Sprintf("- %v\n", encounter.Pokemon.Name))
		}

		cacheString := strings.Join(foundPokemon, "")
		cache.Add(url, []byte(cacheString))
		printRes = cacheString
	}

	fmt.Print(printRes)

	return nil
}

func commandCatch(cptr *config, args []string) error {

	pokemonName := args[0]
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName + "/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	stringResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var pokemon Pokemon
	err = json.Unmarshal(stringResp, &pokemon)
	if err != nil {
		return err
	}
	chance := 100 - (float64(pokemon.BaseExperience) * 0.15)

	rand.Seed(time.Now().UnixNano())

	draw := float64(rand.Intn(101))

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)
	if chance >= draw {
		if pokedex == nil {
			pokedex = &Pokedex{
				caught: map[string]Pokemon{},
			}
		}
		pokedex.caught[pokemonName] = pokemon
		fmt.Printf("%v was caught!\n", pokemonName)
	} else {
		fmt.Printf("%v escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(cptr *config, args []string) error {
	pokemonName := args[0]

	pokemon, ok := pokedex.caught[pokemonName]
	if !ok {
		fmt.Println("This Pokemon is not in your Pokedex! Catch it first to inspect")
		return nil
	}
	printList := []string{}

	printList = append(printList, fmt.Sprintf("Name: %v\n", pokemon.Name))
	printList = append(printList, fmt.Sprintf("Height %v\n", pokemon.Height))
	printList = append(printList, fmt.Sprintf("Weight: %v\n", pokemon.Weight))
	printList = append(printList, "Stats:\n")
	for _, stat := range pokemon.Stats {
		printList = append(printList, fmt.Sprintf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat))
	}
	printList = append(printList, "Types:\n")

	for _, pokeType := range pokemon.Types {
		printList = append(printList, fmt.Sprintf("  - %v\n", pokeType.Type.Name))
	}
	printRes := strings.Join(printList, "")
	fmt.Print(printRes)

	return nil
}
