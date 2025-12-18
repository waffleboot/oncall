package menu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MenuTestSuite struct {
	suite.Suite
	menu *Model
}

func TestMenu(t *testing.T) {
	suite.Run(t, &MenuTestSuite{})
}

func (s *MenuTestSuite) SetupTest() {
	s.menu = New(func(group string, pos int) string {
		switch group {
		case "action 1", "action 2", "action 3":
			return group
		case "items":
			return fmt.Sprintf("item %d", pos+1)
		}
		return ""
	})

	s.menu.AddGroup("action 1")
	s.menu.AddGroupWithItems("items", 3)
	s.menu.AddGroup("action 2")
	s.menu.AddDelimiter()
	s.menu.AddGroup("action 3")
}

func (s *MenuTestSuite) TestMoveCursorDown() {
	testCases := []struct {
		wantGroup string
		wantPos   int
	}{
		{wantGroup: "action 1", wantPos: 0},
		{wantGroup: "items", wantPos: 0},
		{wantGroup: "items", wantPos: 1},
		{wantGroup: "items", wantPos: 2},
		{wantGroup: "action 2", wantPos: 0},
		{wantGroup: "action 3", wantPos: 0},
		{wantGroup: "action 1", wantPos: 0},
	}

	for _, tt := range testCases {
		g, p := s.menu.GetGroup()
		s.Equal(tt.wantGroup, g)
		s.Equal(tt.wantPos, p)
		s.menu.MoveCursorDown()
	}
}

func (s *MenuTestSuite) TestMoveCursorUp() {
	testCases := []struct {
		wantGroup string
		wantPos   int
	}{
		{wantGroup: "action 1", wantPos: 0},
		{wantGroup: "action 3", wantPos: 0},
		{wantGroup: "action 2", wantPos: 0},
		{wantGroup: "items", wantPos: 2},
		{wantGroup: "items", wantPos: 1},
		{wantGroup: "items", wantPos: 0},
		{wantGroup: "action 1", wantPos: 0},
		{wantGroup: "action 3", wantPos: 0},
	}

	for _, tt := range testCases {
		g, p := s.menu.GetGroup()
		s.Equal(tt.wantGroup, g)
		s.Equal(tt.wantPos, p)
		s.menu.MoveCursorUp()
	}
}

//func (s *MenuTestSuite) TestMaxCursor() {
//	s.Equal(5, s.m.maxCursor())
//}

func (s *MenuTestSuite) TestGenerateMenu() {
	s.Run("0", func() {
		s.Equal(`> action 1
  item 1
  item 2
  item 3
  action 2

  action 3
`, s.menu.GenerateMenu())
	})

	s.menu.JumpToGroup("action 2")

	s.Run("4", func() {
		s.Equal(`  action 1
  item 1
  item 2
  item 3
> action 2

  action 3
`, s.menu.GenerateMenu())
	})

	s.menu.JumpToGroup("action 3")

	s.Run("5", func() {
		s.Equal(`  action 1
  item 1
  item 2
  item 3
  action 2

> action 3
`, s.menu.GenerateMenu())
	})

	s.menu.JumpToGroup("unknown")

	g, p := s.menu.GetGroup()

	s.Equal("action 1", g)
	s.Equal(0, p)
}
