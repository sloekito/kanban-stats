package main

import "kanban-stats/trello"

func GetBoardFromTrello(client trello.Client, boardID string) (board Board) {
	trelloLists := client.GetLists(boardID)
	lists := make([]List, len(trelloLists))
	for i, list := range trelloLists {
		lists[i] = List(list)
	}
	return trelloBoard{
		id:     boardID,
		lists:  lists,
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
	id     string
	lists  []List
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
			case defectLabelID:
				foundDefect = true
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
