package main

type PokeResponse struct {
	Count    int          `json:"count"`
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []PokeResult `json:"results"`
}

type PokeResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
