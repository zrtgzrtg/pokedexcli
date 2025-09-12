package main

import (
	"fmt"
	"io"
	"net/http"
)

func explore(cptr *config, args []string) error {

	location := args[0]

	url := "https://pokeapi.co/api/v2/location-area/" + location + "/"
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

	fmt.Println(string(stringResp))

	return nil
}
