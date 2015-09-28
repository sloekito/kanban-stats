package main

import (
	"kanban-stats/trello"
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

func (list List) CountCardsByType(cardType string) (found int) {
	var searchingFor string
	switch cardType{ //TODO pull this data from configuration
	case "defect": searchingFor = "54641fc074d650d56757a68e"
	default: searchingFor = ""
	}
	
	for _, card := range list.Cards {
		for _, label := range card.Labels {
			if label.ID == searchingFor {
				found += 1
				break
			}
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