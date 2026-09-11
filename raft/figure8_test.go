package raft

import "testing"

func TestFigure8CommitGuard(t *testing.T) {

	t.Run("old-term entry not committed by count", func(t *testing.T) {
		n := NewNode(Config{ID: "n0", Peers: []string{"n1", "n2", "n3", "n4"}})
		n.role = Leader
		n.currentTerm = 2

		n.log = []LogEntry{
			{Index: 1, Term: 1},
			{Index: 2, Term: 1},
		}

		n.matchIndex = map[string]int{
			"n1": 2,
			"n2": 2,
			"n3": 0,
			"n4": 0,
		}

		n.commitIndex = 0

		n.advanceCommitIndex()

		if n.commitIndex != 0 {
			t.Fatalf("commitIndex is at to %d, the commit index should be %d", n.commitIndex, 3)
		}

	})

	t.Run("commits once current-term entry replicated", func(t *testing.T) {
		n := NewNode(Config{ID: "n0", Peers: []string{"n1", "n2", "n3", "n4"}})
		n.role = Leader
		n.currentTerm = 2

		n.log = []LogEntry{
			{Index: 1, Term: 1},
			{Index: 2, Term: 1},
			{Index: 3, Term: 2},
		}

		n.matchIndex = map[string]int{
			"n1": 3,
			"n2": 3,
			"n3": 0,
			"n4": 0,
		}

		n.commitIndex = 0

		n.advanceCommitIndex()

		if n.commitIndex != 3 {
			t.Fatalf("commitIndex is at to %d, the commit index should be %d", n.commitIndex, 3)
		}
	})
}
