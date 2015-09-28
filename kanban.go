package main

import (
	"kanban-stats/trello"
)

func GetBoardFromTrello(client trello.Client, boardID string) (board Board){
	return trelloBoard {
		id: boardID,
		columns: TrelloListsToKanbanColumns(client.GetLists(boardID)),
		client: client,
	}
}

type Board interface {
	GetID() string
	GetColumns() []Column
}


type trelloBoard struct {
	id string
	columns []Column
	client trello.Client
}

func (board trelloBoard) GetColumns() []Column {
	return board.columns
}

func (board trelloBoard) GetID() string {
	return board.id
}

func TrelloListsToKanbanColumns(lists []trello.List) (columns []Column){
	columns = make([]Column, len(lists))
	for i, list := range lists {
		columns[i] = Column(list)
	}
	return
}

type Column trello.List

func (column Column) CountCardsByType(cardType string) int{
	return 0
}