package raft

import "testing"

func TestIsLogUpToDate(t *testing.T) {
	tests := []struct {
		name      string
		candTerm  int
		candIndex int
		myTerm    int
		myIndex   int
		want      bool
	}{
		{"higher term wins even with shorter log", 3, 1, 2, 9, true},
		{"lower term loses even with longer log", 2, 9, 3, 1, false},
		{"equal term, candidate longer", 4, 5, 4, 3, true},
		{"equal term, candidate shorter", 4, 3, 4, 5, false},
		{"equal term, equal length", 4, 5, 4, 5, true},
		{"both logs empty", 0, 0, 0, 0, true},
		{"candidate empty, receiver has entries", 0, 0, 4, 2, false},
		{"receiver empty, candidate has entries", 4, 2, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isLogUpToDate(tt.candTerm, tt.candIndex, tt.myTerm, tt.myIndex)
			if got != tt.want {
				t.Errorf("isLogUpToDate(candTerm=%d candIndex=%d myTerm=%d myIndex=%d) = %v, want %v",
					tt.candTerm, tt.candIndex, tt.myTerm, tt.myIndex, got, tt.want)
			}
		})
	}
}
