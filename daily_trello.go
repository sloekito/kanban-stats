package main

import (
	"fmt"
	"log"
	"flag"
	"daily_trello/trello"
	"daily_trello/influx"
)


func main(){
	log.Print("daily_trello: start")
	
	trelloKey := flag.String("trellokey", "", "Trello application key")
	trelloToken :=flag.String("trellotoken", "", "Trello access token")
	trelloBoardID := flag.String("boardid", "", "Trello board ID")
	influxHost := flag.String("influxhost", "", "Influx host:post")
	influxDB := flag.String("influxdb", "", "Influx datbase name")
	influxUser := flag.String("influxuser", "", "Influx username")
	influxPass := flag.String("influxpass", "", "Influx password")
	printOnly := flag.Bool("print-only", false, "Print information rather than publish to Influx" )
	flag.Parse()
	
	trello := trello.Client{
		Key: *trelloKey,
		Token: *trelloToken,
	}
	lists := trello.GetLists(*trelloBoardID)

	if *printOnly == true{
		for _, list := range lists {
			fmt.Printf("%s(%s): %d\n", list.Name, list.Id, len(list.Cards))
		}
	} else {
		influx := influx.Database{
			InfluxHost: *influxHost,
			InfluxDB: *influxDB,
			InfluxUser: *influxUser,
			InfluxPass: *influxPass,
		}
		influx.PublishListsToInflux(*trelloBoardID, lists)
	}	
	
	log.Print("daily_trello: success")	
}