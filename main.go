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
)

var conf config

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
}

func checkAndCallReg(name string, cptr *config) {
	if val, ok := cliCommands[name]; !ok {
		fmt.Println("Unknown command")
	} else {
		val.callback(cptr)
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
	conf = config{
		Next:     *parsedUrl,
		Previous: url.URL{},
	}
	for true {

		fmt.Print("Pokedex > ")

		if scanner.Scan() {

			line := scanner.Text()

			lineClean := cleanInput(line)

			checkAndCallReg(lineClean[0], &conf)

		}
	}

}

func commandExit(cptr *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(cptr *config) error {

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, command := range cliCommands {
		output := fmt.Sprintf("%s: %s", command.name, command.description)
		fmt.Println(output)
	}

	return nil
}

func commandMap(cptr *config) error {
	req, err := http.NewRequest("GET", cptr.Next.String(), nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	stringResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var jsonResponse PokeResponse
	err = json.Unmarshal(stringResp, &jsonResponse)
	if err != nil {
		return err
	}
	parsedNext, err := url.Parse(jsonResponse.Next)
	if err != nil {
		return err
	}
	parsedPrevious, err := url.Parse(jsonResponse.Previous)
	if err != nil {
		return err
	}

	conf = config{
		Next:     *parsedNext,
		Previous: *parsedPrevious,
	}

	for _, area := range jsonResponse.Results {
		fmt.Println(area.Name, area.Url)
	}

	return nil
}
func commandMapb(cptr *config) error {
	req, err := http.NewRequest("GET", cptr.Previous.String(), nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	stringResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var jsonResponse PokeResponse
	err = json.Unmarshal(stringResp, &jsonResponse)
	if err != nil {
		return err
	}
	parsedNext, err := url.Parse(jsonResponse.Next)
	if err != nil {
		return err
	}
	parsedPrevious, err := url.Parse(jsonResponse.Previous)
	if err != nil {
		return err
	}

	conf = config{
		Next:     *parsedNext,
		Previous: *parsedPrevious,
	}

	for _, area := range jsonResponse.Results {
		fmt.Println(area.Name, area.Url)
	}

	return nil

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
