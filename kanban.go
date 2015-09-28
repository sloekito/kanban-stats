package main

import (
	"kanban-stats/trello"
)

func GetBoardFromTrello(client trello.Client, boardID string) (board Board){
	return Board {
		Id: boardID,
		Columns: TrelloListsToKanbanColumns(client.GetLists(boardID)),
		client: client,
	}
}

type Board struct {
	Id string
	Columns []Column
	client trello.Client
}

func TrelloListsToKanbanColumns(lists []trello.List) (columns []Column){
	columns = make([]Column, len(lists))
	for i, list := range lists {
		columns[i] = Column(list)
	}
	return
}

type Column trello.List

func CountCardsByType(column Column, cardType string) int{
	return 0
}