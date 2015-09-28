package main

import(
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	
	"kanban-stats/mocks"
	"kanban-stats/trello"
)

func TestGetBoardFromTrello(t *testing.T){
	Convey("Given a Trello board", t, func() {
		client := new(mocks.Client)
		boardId := "abc"
		lists := []trello.List{{ID: "a"}, {ID: "b"}}
		client.On("GetLists", boardId).Return(lists)
		
			
		Convey("When calling GetBoardFromTrello", func() {
			board := GetBoardFromTrello(client, boardId)
			
			Convey("The returned board is populated", func(){
				So(board.GetID(), ShouldEqual, boardId)
				So(board.GetColumns()[0].GetID(), ShouldEqual, lists[0].ID)
				So(board.GetColumns()[1].GetID(), ShouldEqual, lists[1].ID)
				
			})
		})
	})
}

func TestCountCardsByType(t *testing.T){
	Convey("Given a column with 3 cards labeled 'Defect' and 2 cards with no labels", t, func(){
		
	})
}