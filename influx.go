package main

import (
	influxdb "github.com/influxdb/influxdb/client"
//	"kanban-stats/trello"
//	"log"
)

func writePointsToDatabase(c *influxdb.Client, points []influxdb.Point) error {
	
	return nil
}

// Tagsys Label: 54641fc074d650d56757a692
// Defect Label: 54641fc074d650d56757a68e


/*
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
}*/