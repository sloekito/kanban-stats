package main

import (
	"fmt"
	"log"
	"flag"
	"net/url"
	
	"kanban-stats/trello"
	influxdb "github.com/influxdb/influxdb/client"
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

	trello := trello.NetworkClient{
		Key:   flags.trelloKey,
		Token: flags.trelloToken,
	}
	
	board := GetBoardFromTrello(trello, flags.trelloBoardID)

	if flags.verbose { printInfo(board) }	
	if !flags.dryRun {
		influxURL, err := url.Parse(fmt.Sprintf("http://%s:%d", flags.influxHost, 8086))
		if err != nil {
			log.Fatal(ApplicationName, ": url.Parse: ", err)
		}
		
		influxdb, err := influxdb.NewClient(influxdb.Config{
			URL: *influxURL,
			Username: flags.influxUser,
			Password: flags.influxPassword,
		})
		if err != nil {
			log.Fatal(ApplicationName, ": influxdb.NewClient: ", err)
		}
	
		err = writeStatsToDatabase(influxdb, board)
		if err != nil {
			log.Fatal(ApplicationName, ": writeStatsToDatabase: ", err)
		}			
		//writeListsToDatabase(influxdbClient, board.Columns, flags.trelloBoardID)
	}
}

func printInfo(board Board) {
	fmt.Printf("Board ID: %v\n", board.GetID())
	for _, column := range board.GetColumns() {
		fmt.Printf("%s(%s): %d\n", column.GetName(), column.GetID(), len(column.GetCards()))
	}
}