package main

import (
	"net/http"
	"net/url"
	"encoding/json"
	"bytes"
	"log"
)


const (
	InfluxHost = "192.168.99.100:8086"
	InfluxDB = "Trello"
	InfluxUser = "root"
	InfluxPass = "root"
)

type InfluxPost [1]struct {
	Name string `json:"name"`
	Columns []string `json:"columns"`
	Points [1][]int `json:"points"`
}

func PublishListsToInflux(seriesName string, lists []List){
	post := InfluxPost{{Name: seriesName}}
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
		
	response, err := http.Post(influxURL.String(), "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	
	if response.StatusCode != 200 {
		log.Fatal("InfluxDB response: " + response.Status)
	}
}