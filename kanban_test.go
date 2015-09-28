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
		nonDevLabel := trello.Label{ ID: "54641fc074d650d56757a68d"}
		otherLabel := trello.Label{ ID: "1234567890" }
		cards := []trello.Card{
			{Labels: []trello.Label{defectLabel, otherLabel}},
			{Labels: []trello.Label{defectLabel, nonDevLabel}},
			{Labels: []trello.Label{otherLabel}},
			{Labels: []trello.Label{}},
			{Labels: []trello.Label{nonDevLabel}},
		}
		list := List{Cards: cards}
		
		Convey("Counting cards by type defect returns the number of cards that have the defect label applied and not the non-dev label", func(){
			So(list.CountCardsByType("defect"), ShouldEqual, 1)
		})
		
		Convey("Counting cards by type feature returns the number of cards that do not have the defect label appllied and do not have the non-dev label", func(){
			So(list.CountCardsByType("feature"), ShouldEqual, 2)
		})
	})
}

func Test_TrelloBoard_GetMeasurementPoints(t *testing.T){
	Convey("Given a board full of cards", t, func(){
		defectLabel := trello.Label{ ID: "54641fc074d650d56757a68e" }
		nonDevLabel := trello.Label{ ID: "54641fc074d650d56757a68d"}
		otherLabel := trello.Label{ ID: "1234567890" }
		cards1 := []trello.Card{
			{Labels: []trello.Label{defectLabel, otherLabel}},
			{Labels: []trello.Label{defectLabel}},
			{Labels: []trello.Label{otherLabel}},
			{Labels: []trello.Label{}},
			{Labels: []trello.Label{}},
		}
		cards2 := []trello.Card{
			{Labels: []trello.Label{nonDevLabel}},
		}
		list1 := List{ID: "list1", Cards: cards1}
		list2 := List{ID: "list2", Cards: cards2}
		board := trelloBoard{id: "board1", lists: []List{list1, list2}}
		
		Convey("It returns a structure that influxdb can use", func(){
			result := board.GetMeasurementPoints()
			
			So(len(result), ShouldEqual, 4)
			So(result[0].Measurement, ShouldEqual, "count_cards")
			So(result[0].Tags["board"], ShouldEqual, "board1")
			So(result[0].Tags["list"], ShouldEqual, "list1")
			So(result[0].Tags["type"], ShouldEqual, "feature")
			So(result[0].Fields["value"], ShouldEqual, 3)
			So(result[1].Tags["type"], ShouldEqual, "defect")
			So(result[1].Tags["list"], ShouldEqual, "list1")
			So(result[1].Fields["value"], ShouldEqual, 2)
			So(result[2].Tags["type"], ShouldEqual, "feature")
			So(result[2].Tags["list"], ShouldEqual, "list2")
			So(result[2].Fields["value"], ShouldEqual, 0)
			So(result[3].Fields["value"], ShouldEqual, 0)
			
		})
	})
}