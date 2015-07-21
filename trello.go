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

func GetTrelloLists(boardID string) (lists []List) {
	query := url.Values{
		"key": {key}, 
		"token": {token},
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
	
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&lists)
	if err != nil {
		log.Fatal(err)
	}
	return
}