package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	pokecache "github.com/zrtgzrtg/pokedexcli/internal"
)

var conf *config
var cache *pokecache.Cache
var pokedex *Pokedex

var cliCommands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"map": {
		name:        "map",
		description: "prints next 20 locations",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "prints previous 20 locations",
		callback:    commandMapb,
	},
	"explore": {
		name:        "explore",
		description: "print all pokemon for a given location",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "attempts to catch pokemon by throwing a pokeball",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "prints data about a pokemon if available in pokedex",
		callback:    commandInspect,
	},
}

func checkAndCallReg(name string, cptr *config, args []string) {
	if val, ok := cliCommands[name]; !ok {
		fmt.Println("Unknown command")
	} else {
		val.callback(cptr, args)
	}
}
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cliCommands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}

	parsedUrl, _ := url.Parse("https://pokeapi.co/api/v2/location-area/")
	conf = new(config)
	conf.Next = *parsedUrl
	conf.Previous = url.URL{}
	for {

		fmt.Print("Pokedex > ")

		args := []string{}

		if scanner.Scan() {

			line := scanner.Text()

			lineClean := cleanInput(line)

			if len(lineClean) > 1 {
				for i, val := range lineClean {
					if i == 0 {
						continue
					}
					args = append(args, val)

				}
			}

			checkAndCallReg(lineClean[0], conf, args)

		}
	}

}

func commandExit(cptr *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(cptr *config, args []string) error {

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, command := range cliCommands {
		output := fmt.Sprintf("%s: %s", command.name, command.description)
		fmt.Println(output)
	}

	return nil
}

func commandMap(cptr *config, args []string) error {

	// save current call to assign as previous
	currentURL := cptr.Next

	//parse http.Response initialized for caching purpose
	var stringResp []byte

	// cache initializing
	if cache == nil {
		cache = pokecache.NewCache(30 * time.Second)
	}
	// look if this url cached already

	val, ok := cache.Get(currentURL.String())
	if ok {
		fmt.Println("cache hit")
		stringResp = val
	} else {
		req, err := http.NewRequest("GET", currentURL.String(), nil)
		if err != nil {
			return err
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		stringResp, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	jsonResponse, parsedNext, _, err := getJson(stringResp)
	if err != nil {
		return err
	}

	// add Previous to Cache
	cache.Add(currentURL.String(), stringResp)

	conf.Next = parsedNext
	conf.Previous = currentURL

	for _, area := range jsonResponse.Results {
		fmt.Println(area.Name, area.Url)
	}
	for key, _ := range cache.CacheEntries {

		fmt.Println(string(key))
	}

	return nil
}
func commandMapb(cptr *config, args []string) error {

	// save current call to assign as previous
	currentURL := cptr.Previous

	// mapb just copies logic from map. Look at comments from map to understand

	var stringResp []byte
	val, ok := cache.Get(currentURL.String())
	if ok {
		fmt.Println("cache hit")
		stringResp = val
	} else {
		req, err := http.NewRequest("GET", currentURL.String(), nil)
		if err != nil {
			return err
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		stringResp, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	jsonResponse, _, parsedPrevious, err := getJson(stringResp)
	if err != nil {
		return err
	}

	conf.Next = currentURL
	conf.Previous = parsedPrevious

	for _, area := range jsonResponse.Results {
		fmt.Println(area.Name, area.Url)
	}

	return nil

}

// map helper for unmarshalling json

func getJson(stringResp []byte) (resp PokeResponse, prev, next url.URL, err error) {
	var jsonResponse PokeResponse
	err = json.Unmarshal(stringResp, &jsonResponse)
	if err != nil {
		return PokeResponse{}, url.URL{}, url.URL{}, err
	}
	parsedNext, err := url.Parse(jsonResponse.Next)
	if err != nil {
		return PokeResponse{}, url.URL{}, url.URL{}, err
	}
	parsedPrevious, err := url.Parse(jsonResponse.Previous)
	if err != nil {
		return PokeResponse{}, url.URL{}, url.URL{}, err
	}
	return jsonResponse, *parsedNext, *parsedPrevious, nil
}

func cleanInput(text string) []string {
	if text == "" {
		return []string{}
	}

	lower := strings.ToLower(text)
	byteS := []byte(lower)
	const SEP = byte(0)
	longString := []byte{}

	for i, b := range byteS {
		if b == ' ' {
			if len(longString) > 0 && longString[len(longString)-1] != SEP {
				longString = append(longString, SEP)
			}
			continue
		}
		longString = append(longString, b)

		if i+1 < len(byteS) && byteS[i+1] == ' ' {
			longString = append(longString, SEP)
		}
	}

	res := []string{}
	word := []byte{}
	for _, b := range longString {
		if b == SEP {
			if len(word) > 0 {
				res = append(res, string(word))
				word = []byte{}
			}
		} else {
			word = append(word, b)
		}
	}
	if len(word) > 0 {
		res = append(res, string(word))
	}

	return res
}
