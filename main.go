package main

import (
	"fmt"
	"log"
	"flag"
	"kanban-stats/trello"
	influxdb "kanban-stats/internal/github.com/influxdb/influxdb/client"
)

type flags struct {
	trelloKey, trelloToken, trelloBoardID, influxHost, influxDB, influxUser, influxPassword string
	verbose, dryRun bool
}

func getCommandLineFlags() (flags flags) {
	flag.StringVar(&flags.trelloKey, "trellokey", "", "Trello application key")
	flag.StringVar(&flags.trelloToken, "trellotoken", "", "Trello access token")
	flag.StringVar(&flags.trelloBoardID, "boardid", "", "Trello board ID")
	flag.StringVar(&flags.influxHost, "influxhost", "", "Influx host:post")
	flag.StringVar(&flags.influxDB, "influxdb", "", "Influx datbase name")
	flag.StringVar(&flags.influxUser, "influxuser", "", "Influx username")
	flag.StringVar(&flags.influxPassword, "influxpass", "", "Influx password")
	flag.BoolVar(&flags.verbose, "v", false, "Print verbose information")
	flag.BoolVar(&flags.dryRun, "d", false, "Dry run does not output to database")
	flag.Parse()
	return
}

func main() {
	log.Print(ApplicationName, ": start")
	defer log.Print(ApplicationName, ": end")
	
	flags := getCommandLineFlags()

	trello := trello.Client{
		Key:   flags.trelloKey,
		Token: flags.trelloToken,
	}
	lists := trello.GetLists(flags.trelloBoardID)

	if flags.verbose == true {
		fmt.Printf("Board ID: %v\n", flags.trelloBoardID)
		for _, list := range lists {
			fmt.Printf("%s(%s): %d\n", list.Name, list.Id, len(list.Cards))
		}
	}
	
	if !flags.dryRun {
		influxdbClient, err := influxdb.NewClient(&influxdb.ClientConfig{
			Host: flags.influxHost,
			Username: flags.influxUser,
			Password: flags.influxPassword,
			Database: flags.influxDB,
		})
		if err != nil {
			log.Fatal(err)
		}
	
		writeListsToDatabase(influxdbClient, lists, flags.trelloBoardID)
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