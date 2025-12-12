package tea

import "testing"

func TestMenu(t *testing.T) {
	var m menu

	m.addGroup("action 1", 1)
	m.addGroup("items", 3)
	m.addGroup("action 2", 1)
	m.addGroup("action 3", 1)

	testCases := []struct {
		givenCursor int
		wantGroup   string
		wantPos     int
	}{
		{
			givenCursor: 0,
			wantGroup:   "action 1",
			wantPos:     0,
		},
		{
			givenCursor: 1,
			wantGroup:   "items",
			wantPos:     0,
		},
		{
			givenCursor: 2,
			wantGroup:   "items",
			wantPos:     1,
		},
		{
			givenCursor: 3,
			wantGroup:   "items",
			wantPos:     2,
		},
		{
			givenCursor: 4,
			wantGroup:   "action 2",
			wantPos:     0,
		},
		{
			givenCursor: 5,
			wantGroup:   "action 3",
			wantPos:     0,
		},
	}

	for _, tt := range testCases {
		t.Run("", func(t *testing.T) {
			g, p := m.getGroup(tt.givenCursor)
			if g != tt.wantGroup || p != tt.wantPos {
				t.Fail()
			}
		})
	}

	t.Run("", func(t *testing.T) {
		if m.maxCursor() != 5 {
			t.Fail()
		}
	})
}
