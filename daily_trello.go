package main

import (
	"fmt"
	"log"
	"flag"
	"daily_trello/trello"
	"daily_trello/influx"
)


func main(){
	log.Print("---daily_trello: start")
	
	trelloKey := flag.String("trellokey", "", "Trello application key")
	trelloToken :=flag.String("trellotoken", "", "Trello access token")
	trelloBoardID := flag.String("boardid", "", "Trello board ID")
	influxHost := flag.String("influxhost", "", "Influx host:post")
	influxDB := flag.String("influxdb", "", "Influx datbase name")
	influxUser := flag.String("influxuser", "", "Influx username")
	influxPass := flag.String("influxpass", "", "Influx password")
	flag.Parse()
	
	trello := trello.Client{
		Key: *trelloKey,
		Token: *trelloToken,
	}
	lists := trello.GetLists(*trelloBoardID)
	fmt.Println(trelloBoardID)
	for _, list := range lists {
		fmt.Printf("%s(%s): %d\n", list.Name, list.Id, len(list.Cards))
	}

	influx := influx.Client{
		InfluxHost: *influxHost,
		InfluxDB: *influxDB,
		InfluxUser: *influxUser,
		InfluxPass: *influxPass,
	}
	influx.PublishListsToInflux(*trelloBoardID, lists)
	
	
	log.Print("---daily_trello: success")	
}