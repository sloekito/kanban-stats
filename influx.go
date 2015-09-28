package main

import (
	influxdb "github.com/influxdb/influxdb/client"
)

func writePointsToDatabase(client *influxdb.Client, points []influxdb.Point) error {
	bps := influxdb.BatchPoints{
		Points: points,
		Database: "trello",
		RetentionPolicy: "default",
	}
	_, err := client.Write(bps)
	return err
}