package main

import(
	"testing"
	"reflect"
	. "github.com/smartystreets/goconvey/convey"
	
	"kanban-stats/mocks"
	"kanban-stats/trello"
)

func TestGetBoardFromTrello(t *testing.T){
	Convey("Given a Trello board", t, func() {
		client := new(mocks.Client)
		boardId := "abc"
		lists := make([]trello.List,5)
		client.On("GetLists", boardId).Return(lists)
		
			
		Convey("When calling GetBoardFromTrello", func() {
			board := GetBoardFromTrello(client, boardId)
			
			Convey("The returned board is populated", func(){
				So(board.GetID(), ShouldEqual, boardId)
				So(reflect.DeepEqual(board.GetColumns(), TrelloListsToKanbanColumns(lists)), ShouldBeTrue)
			})
		})
	})
}

func TestCountCardsByType(t *testing.T){
	Convey("Given a column with 3 cards labeled 'Defect' and 2 cards with no labels", t, func(){
		
	})
}