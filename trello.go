package main

import (
	"net/url"
	"net/http"
	"log"
	"encoding/json"
)

type List struct {
	Id string
	Name string
	Cards []interface{}
}

type Client struct {
	Key string
	Token string
}

func (c Client) GetTrelloLists(boardID string) (lists []List) {
	query := url.Values{
		"key": {c.Key}, 
		"token": {c.Token},
		"cards": {"open"},
		"card_fields": {"idShort"},
	}
	boardUrl := url.URL{
		Scheme: "https",
		Host: "api.trello.com",
		Path: "1/board/" + boardID + "/lists",
		RawQuery: query.Encode(),
	}
	response, err := http.Get(boardUrl.String())
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatal("Trello response: " + response.Status)
	}
	
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&lists)
	if err != nil {
		log.Fatal(err)
	}
	return
}