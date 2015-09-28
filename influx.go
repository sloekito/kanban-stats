package main

import (
	influxdb "github.com/influxdb/influxdb/client"
//	"kanban-stats/trello"
//	"log"
)

func writeStatsToDatabase(c *influxdb.Client, board Board) error {
	points := make([]influxdb.Point, len(board.GetColumns())*2)
	teams := []string{"rfid_nordstrom"}
	cardTypes := []string{"feature", "defect"}
	i := 0 
	for _, column := range board.GetColumns() {
		for _, team := range teams{
			for _, cardType := range cardTypes{
				points[i] = influxdb.Point{
					Measurement: "count_cards",
					Tags: map[string]string {
						"board": board.GetID(),
						"list": column.Id,
						"team": team,
						"type": cardType,
					},
					Fields: map[string]interface{}{
						"value": column.CountCardsByType(cardType),
					},
				}
			}
		}		
	}
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