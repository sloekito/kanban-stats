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
	query := url.Values{
		"key": {key}, 
		"token": {token},
		"cards": {"all"},
		"card_fields": {"idShort"},
	}
	boardUrl := url.URL{
		Scheme: "https",
		Host: "api.trello.com",
		Path: "1/board/534efa63a3a33edc034ac3d1/lists",
		RawQuery: query.Encode(),
	}
	response, err := http.Get(boardUrl.String())
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	
	decoder := json.NewDecoder(response.Body)
	var lists []List
	err = decoder.Decode(&lists)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, list := range lists {
		fmt.Printf("%s: %d\n", list.Name, len(list.Cards))
	}	
}