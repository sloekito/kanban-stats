package influx

import (
	"net/http"
	"net/url"
	"encoding/json"
	"bytes"
	"log"
	"daily_trello/trello"
	"io/ioutil"
)


type Database struct{
	InfluxHost string
	InfluxDB string
	InfluxUser string
	InfluxPass string
}

type InfluxPost [1]struct {
	Name string `json:"name"`
	Columns []string `json:"columns"`
	Points [1][]int `json:"points"`
}

func (c Database)PublishListsToInflux(seriesName string, lists []trello.List){
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
		"u": {c.InfluxUser}, 
		"p": {c.InfluxPass},
	}
	influxURL := url.URL{
		Scheme: "http",
		Host: c.InfluxHost,
		Path: "db/" + c.InfluxDB + "/series",
		RawQuery: query.Encode(),
	}
		
	response, err := http.Post(influxURL.String(), "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		log.Fatalf("InfluxDB response: %v\n%v", response.Status, string(body))
	}
}