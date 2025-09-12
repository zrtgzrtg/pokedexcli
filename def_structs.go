package main

import (
	"net/url"
)

type config struct {
	Next     url.URL
	Previous url.URL
}
type cliCommand struct {
	name        string
	description string
	callback    func(cptr *config, args []string) error
}
type Pokedex struct {
	caught map[string]Pokemon
}
