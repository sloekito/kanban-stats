package main

import (
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

func GetMeasurementPoints(board Board) (points influxdb.BatchPoints) {
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

/*
func (board trelloBoard) GetMeasurementPoints() (points []influxdb.Point) {
	points = make([]influxdb.Point, len(board.GetColumns())*2)
	//	teams := []string{"rfid_nordstrom"}
	cardTypes := []string{"feature", "defect"}
	i := 0
	for _, column := range board.GetColumns() {
		//		for _, team := range teams{
		for _, cardType := range cardTypes {
			points[i] = influxdb.Point{
				Measurement: "count_cards",
				Tags: map[string]string{
					"board": board.GetID(),
					"list":  column.GetID(),
					//"team": team,
					"type": cardType,
				},
				Fields: map[string]interface{}{
					"value": column.CountCardsByType(cardType),
				},
			}
			i++
		}
		//		}
	}
	return
}

func writePointsToDatabase(client influxdb.Client, points []influxdb.Point) error {
	bps, _ := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:        "trello",
		RetentionPolicy: "default",
	})

	for _, point := range points {
		bps.AddPoint(&point)
	}

	err := client.Write(bps)
	return err
}
*/
