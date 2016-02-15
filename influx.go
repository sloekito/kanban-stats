package main

import (
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

func BuildPointBatch(board Board) (points influxdb.BatchPoints) {
	points, _ = influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:        "kanban",
		RetentionPolicy: "default",
	})

	cardTypes := []string{"feature", "defect"}
	for _, column := range board.GetColumns() {
		//		for _, team := range teams{
		for _, cardType := range cardTypes {
			tags := map[string]string{
				"board": board.GetID(),
				"list":  column.GetID(),
				//"team": team,
				"type": cardType,
			}
			fields := map[string]interface{}{"value": column.CountCardsByType(cardType)}
			point, _ := influxdb.NewPoint("count_cards", tags, fields, time.Now())
			points.AddPoint(point)
		}
		//		}
	}
	return
}

func writePointsToDatabase(client influxdb.Client, batchPoints influxdb.BatchPoints) error {
	err := client.Write(batchPoints)
	return err
}
