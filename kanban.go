package main

import (
	"kanban-stats/trello"
)

func GetBoardFromTrello(client trello.Client, boardID string) (board Board){
	return Board {
		Id: boardID,
		Columns: client.GetLists(boardID),
		client: client,
	}
}

type Board struct {
	Id string
	Columns []trello.List
	client trello.Client
}

type Column trello.List

