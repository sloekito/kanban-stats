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

func Test_TrelloBoard_CountCardsByType(t *testing.T){
	Convey("Given a list with cards with a mix of labels", t, func(){
		defectLabel := trello.Label{ ID: "54641fc074d650d56757a68e" }
		otherLabel := trello.Label{ ID: "1234567890" }
		cards := []trello.Card{
			{Labels: []trello.Label{defectLabel, otherLabel}},
			{Labels: []trello.Label{defectLabel}},
			{Labels: []trello.Label{otherLabel}},
			{Labels: []trello.Label{}},
			{Labels: []trello.Label{}},
		}
		list := List{Cards: cards}
		
		Convey("Counting cards by type defect returns the number of cards that have the defect label applied", func(){
			So(list.CountCardsByType("defect"), ShouldEqual, 2)
		})
		
		Convey("Counting cards by type feature returns the number of cards that do not have the defect label appllied", func(){
			So(list.CountCardsByType("feature"), ShouldEqual, 3)
		})
	})
}