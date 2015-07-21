package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"encoding/json"
	"bytes"
)

type List struct {
	Id string
	Name string
	Cards []interface{}
}

func main(){
	log.Print("---daily_trello: start")	
	lists := GetTrelloLists("534efa63a3a33edc034ac3d1")
	PublishListsToInflux(lists)
	
	for _, list := range lists {
		fmt.Printf("%s: %d\n", list.Name, len(list.Cards))
	}
	
	log.Print("---daily_trello: success")	
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

const (
	TrelloBoardName = "DevProcess"	
)

const (
	InfluxHost = "192.168.99.100:8086"
	InfluxDB = "Trello"
	InfluxUser = "root"
	InfluxPass = "root"
)

type InfluxPost struct {
	Name string `json:"name"`
	Columns []string `json:"columns"`
	Points [1][]int `json:"points"`
}

func PublishListsToInflux(lists []List){
	post := [1]InfluxPost{
		{Name: TrelloBoardName,},
	}
	columns := make([]string, len(lists))
	point := make([]int, len(lists))
	points := [1][]int{point}
	for i, list := range lists {
		columns[i] = list.Id + "_count_open_cards"
		points[0][i] = len(list.Cards)
	}
	
	post[0].Columns = columns
	post[0].Points = points
	
	body, err := json.Marshal(post)
	if err != nil {
		log.Fatal(err)
	}
	
	query := url.Values{
		"u": {InfluxUser}, 
		"p": {InfluxPass},
	}
	influxURL := url.URL{
		Scheme: "http",
		Host: InfluxHost,
		Path: "db/" + InfluxDB + "/series",
		RawQuery: query.Encode(),
	}
		
	response, err := http.Post(influxURL.String(), "application/x-www-form-urlencoded", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	
	if response.StatusCode != 200 {
		log.Fatal("InfluxDB response: " + response.Status)
	}
}