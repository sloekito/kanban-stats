package mocks

import "kanban-stats/trello"
import "github.com/stretchr/testify/mock"

type Client struct {
	mock.Mock
}

func (m *Client) GetLists(boardID string) []trello.List {
	ret := m.Called(boardID)

	var r0 []trello.List
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]trello.List)
	}

	return r0
}
