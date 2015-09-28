package main

import (
	"kanban-stats/trello"
	influxdb "github.com/influxdb/influxdb/client"
)

func GetBoardFromTrello(client trello.Client, boardID string) (board Board){
	trelloLists := client.GetLists(boardID)
	lists := make([]List, len(trelloLists))
	for i, list := range trelloLists {
		lists[i] = List(list)
	}
	return trelloBoard {
		id: boardID,
		lists: lists,
		client: client,
	}
}

type Board interface {
	GetID() string
	GetColumns() []Column
	GetMeasurementPoints() []influxdb.Point
}

type Column interface {
	CountCardsByType(string) int
	GetID() string
	GetName() string
	GetCards() []trello.Card
}


type trelloBoard struct {
	id string
	lists []List
	client trello.Client
}

type List trello.List

func (board trelloBoard) GetColumns() []Column {
	columns := make([]Column, len(board.lists))
	for i, list := range board.lists {
		columns[i] = Column(list)
	}
	return columns
}

func (board trelloBoard) GetID() string {
	return board.id
}

func (board trelloBoard) GetMeasurementPoints() (points []influxdb.Point) {
	points = make([]influxdb.Point, len(board.GetColumns())*2)
//	teams := []string{"rfid_nordstrom"}
	cardTypes := []string{"feature", "defect"}
	i := 0 
	for _, column := range board.GetColumns() {
//		for _, team := range teams{
			for _, cardType := range cardTypes{
				points[i] = influxdb.Point{
					Measurement: "count_cards",
					Tags: map[string]string {
						"board": board.GetID(),
						"list": column.GetID(),
						//"team": team,
						"type": cardType,
					},
					Fields: map[string]interface{}{
						"value": column.CountCardsByType(cardType),
					},
				}
				i += 1
			}
//		}		
	}
	return
}


// Tagsys Label: 54641fc074d650d56757a692
// Defect Label: 54641fc074d650d56757a68e
func (list List) CountCardsByType(cardType string) (found int) {
	defectLabelID := "54641fc074d650d56757a68e" //TODO: Move these to configuration
	nonDevLabelID := "54641fc074d650d56757a68d"

	for _, card := range list.Cards {
		var foundDefect, foundNonDev bool
		for _, label := range card.Labels {
			switch label.ID {
			case nonDevLabelID:
				foundNonDev = true
				break
			case defectLabelID: foundDefect = true
			}
		}
		if !foundNonDev && (cardType == "defect" && foundDefect || cardType != "defect" && !foundDefect) {
			found += 1
		}
	}
	return
}

func (list List) GetID() string {
	return list.ID
}

func (list List) GetName() string {
	return list.Name
}

func (list List) GetCards() []trello.Card {
	return list.Cards
}