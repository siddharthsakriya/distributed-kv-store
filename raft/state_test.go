package raft

import "testing"

func TestNewNode(t *testing.T) {
	t.Run("fresh node, no persister", func(t *testing.T) {
		n := NewNode(Config{
			ID:        "node1",
			Peers:     []string{"node2", "node3"},
			Persister: nil,
		})

		if n.role != Follower {
			t.Errorf("role = %v, want Follower", n.role)
		}
		if n.currentTerm != 0 {
			t.Errorf("currentTerm = %d, want 0", n.currentTerm)
		}
		if n.votedFor != "" {
			t.Errorf("votedFor = %q, want empty string", n.votedFor)
		}
		if len(n.log) != 0 {
			t.Errorf("len(log) = %d, want 0", len(n.log))
		}
	})

	t.Run("restore from persister", func(t *testing.T) {
		saved := &FakePersister{
			term:     5,
			votedFor: "node2",
			log: []LogEntry{
				{Index: 1, Term: 1},
				{Index: 2, Term: 4},
			},
		}

		n := NewNode(Config{
			ID:        "node1",
			Peers:     []string{"node2", "node3"},
			Persister: saved,
		})

		if n.currentTerm != 5 {
			t.Errorf("currentTerm = %d, want 5", n.currentTerm)
		}
		if n.votedFor != "node2" {
			t.Errorf("votedFor = %q, want node2", n.votedFor)
		}
		if len(n.log) != 2 {
			t.Fatalf("len(log) = %d, want 2", len(n.log))
		}
		if n.log[1].Index != 2 || n.log[1].Term != 4 {
			t.Errorf("log[1] = %+v, want {Index:2 Term:4}", n.log[1])
		}

		if n.role != Follower {
			t.Errorf("role = %v, want Follower", n.role)
		}
		if n.commitIndex != 0 {
			t.Errorf("commitIndex = %d, want 0", n.commitIndex)
		}
		if n.lastApplied != 0 {
			t.Errorf("lastApplied = %d, want 0", n.lastApplied)
		}
	})
}
