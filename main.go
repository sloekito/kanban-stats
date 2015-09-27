package main

import (
	"fmt"
	"log"
	"flag"
	"trello-stats/trello"
	influxdb "trello-stats/internal/github.com/influxdb/influxdb/client"
)

type args struct {
	trelloKey, trelloToken, trelloBoardID, influxHost, influxDB, influxUser, influxPassword *string
	verbose, dryRun *bool
}

func parseArgs() (_args args) {
	_args = args{
		trelloKey: flag.String("trellokey", "", "Trello application key"),
		trelloToken: flag.String("trellotoken", "", "Trello access token"),
		trelloBoardID: flag.String("boardid", "", "Trello board ID"),
		influxHost: flag.String("influxhost", "", "Influx host:post"),
		influxDB: flag.String("influxdb", "", "Influx datbase name"),
		influxUser: flag.String("influxuser", "", "Influx username"),
		influxPassword: flag.String("influxpass", "", "Influx password"),
		verbose: flag.Bool("v", false, "Print verbose information"),
		dryRun: flag.Bool("d", false, "Dry run does not output to database"),
	}
	flag.Parse()
	return
}

func main() {
	log.Print(ApplicationName, ": start")
	defer log.Print(ApplicationName, ": end")
	
	config := parseArgs()

	trello := trello.Client{
		Key:   *config.trelloKey,
		Token: *config.trelloToken,
	}
	lists := trello.GetLists(*config.trelloBoardID)

	if *config.verbose == true {
		fmt.Printf("Board ID: %v\n", *config.trelloBoardID)
		for _, list := range lists {
			fmt.Printf("%s(%s): %d\n", list.Name, list.Id, len(list.Cards))
		}
	}
	
	if !*config.dryRun {
		influxdbClient, err := influxdb.NewClient(&influxdb.ClientConfig{
			Host: *config.influxHost,
			Username: *config.influxUser,
			Password: *config.influxPassword,
			Database: *config.influxDB,
		})
		if err != nil {
			log.Fatal(err)
		}
	
		writeListsToDatabase(influxdbClient, lists, *config.trelloBoardID)
	}
}

func writeListsToDatabase(client *influxdb.Client, lists []trello.List, seriesName string){
	point := make([]interface{}, len(lists))
	series := influxdb.Series{	
		Name: seriesName,
		Columns: make([]string, len(lists)),
		Points: [][]interface{}{point},
	}
	for i, list := range lists {
		series.Columns[i] = list.Id + "_count_open_cards"
		point[i] = len(list.Cards)
	}
	
	err := client.WriteSeries([]*influxdb.Series{&series})
	if err != nil {
		log.Fatal(err)
	}
}