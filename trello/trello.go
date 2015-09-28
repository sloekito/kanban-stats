package trello

import (
	"net/url"
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
)

type List struct {
	Id string
	Name string
	Cards []interface{}
}

type Client interface {
	GetLists(boardID string) (lists []List)
}

type NetworkClient struct {
	Key string
	Token string
}

func (c NetworkClient) GetLists(boardID string) (lists []List) {
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

func (c NetworkClient) PrintLabels(boardID string) {
	query := url.Values{
		"key": {c.Key}, 
		"token": {c.Token},
	}
	boardUrl := url.URL{
		Scheme: "https",
		Host: "api.trello.com",
		Path: "1/board/" + boardID + "/labels",
		RawQuery: query.Encode(),
	}
	
	log.Println("GET ", boardUrl.String())
	response, err := http.Get(boardUrl.String())
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatal("Trello response: " + response.Status)
	}
	
	contents, err := ioutil.ReadAll(response.Body)
	log.Println(string(contents))
}