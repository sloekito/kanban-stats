package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"encoding/json"
)

type List struct {
	IdBoard string
	Name string
	Cards []interface{}
}

func main(){	
	lists := GetTrelloLists("534efa63a3a33edc034ac3d1")
	
	for _, list := range lists {
		fmt.Printf("%s: %d\n", list.Name, len(list.Cards))
	}	
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