package main

import (
	"fmt"
	"log"
	"flag"
	"daily_trello/trello"
	influxdb "daily_trello/internal/github.com/influxdb/influxdb/client"
)

/*type args struct {
	trelloKey, trelloToken, trelloBoardID, influxHost, influxDB, influxUser, influxPassword string
	printOnly bool
}

func parseArgs() args {
	var _args args
	
}*/

func main() {
	log.Print("daily_trello: start")

	trelloKey := flag.String("trellokey", "", "Trello application key")
	trelloToken := flag.String("trellotoken", "", "Trello access token")
	trelloBoardID := flag.String("boardid", "", "Trello board ID")
	influxHost := flag.String("influxhost", "", "Influx host:post")
	influxDB := flag.String("influxdb", "", "Influx datbase name")
	influxUser := flag.String("influxuser", "", "Influx username")
	influxPassword := flag.String("influxpass", "", "Influx password")
	printOnly := flag.Bool("print-only", false, "Print information rather than publish to Influx")
	flag.Parse()

	trello := trello.Client{
		Key:   *trelloKey,
		Token: *trelloToken,
	}
	lists := trello.GetLists(*trelloBoardID)

	if *printOnly == true {
		for _, list := range lists {
			fmt.Printf("%s(%s): %d\n", list.Name, list.Id, len(list.Cards))
		}
		
	} else {
		influxdbClient, err := influxdb.NewClient(&influxdb.ClientConfig{
			Host: *influxHost,
			Username: *influxUser,
			Password: *influxPassword,
			Database: *influxDB,
		})
		if err != nil {
			log.Fatal(err)
		}
		
		databases, err := influxdbClient.GetDatabaseList()
		databaseFound := false
		for _,db := range databases {
			if db["name"] == *influxDB {
				databaseFound = true
				break
			}
		}
		
		if !databaseFound {
			err = influxdbClient.CreateDatabase(*influxDB)
			if err != nil {
				log.Fatal(err)
			}
		}
		
		point := make([]interface{}, len(lists))
		series := influxdb.Series{	
			Name: *trelloBoardID,
			Columns: make([]string, len(lists)),
			Points: [][]interface{}{point},
		}
		for i, list := range lists {
			series.Columns[i] = list.Id + "_count_open_cards"
			point[i] = len(list.Cards)
		}
		
		err = influxdbClient.WriteSeries([]*influxdb.Series{&series})
		if err != nil {
			log.Fatal(err)
		}

	}

	log.Print("daily_trello: success")
}
