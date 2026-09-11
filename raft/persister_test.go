package raft

import (
	"reflect"
	"testing"
)

func TestPersister(t *testing.T) {
	t.Run("round-trip, node should be built correctly from the same persister", func(t *testing.T) {
		persister := &FakePersister{}

		n0 := NewNode(Config{ID: "n0", Peers: []string{"n1", "n2", "n3", "n4"}, Persister: persister})
		n0.currentTerm = 1
		n0.votedFor = "n1"
		n0.log = []LogEntry{
			{Index: 1, Term: 1},
		}

		n0.persist()

		n1 := NewNode(Config{ID: "n1", Peers: []string{"n0", "n2", "n3", "n4"}, Persister: persister})

		if n1.currentTerm != n0.currentTerm || n1.votedFor != n0.votedFor || !reflect.DeepEqual(n1.log, n0.log) {
			t.Fatalf("n1 and n0 are not matching :0, the load and save methods not working")
		}
	})

	t.Run("no double vote after crash", func(t *testing.T) {
		persister := &FakePersister{}
		n0 := NewNode(Config{ID: "n0", Peers: []string{"n1", "n2", "n3", "n4"}, Persister: persister})
		n0.currentTerm = 5
		n0.votedFor = "n0"
		n0.log = []LogEntry{
			{Index: 1, Term: 1},
			{Index: 2, Term: 2},
			{Index: 3, Term: 4},
		}
		n0.persist()

		// pretending n0 has come down and restarted
		n0Restarted := NewNode(Config{ID: "n0", Peers: []string{"n0", "n2", "n3", "n4"}, Persister: persister})

		response := n0Restarted.HandleRequestVote(
			&RequestVoteArgs{
				Term:         5,
				CandidateID:  "n2",
				LastLogIndex: 3,
				LastLogTerm:  4,
			},
		)

		if n0Restarted.votedFor != "n0" || response.VoteGranted {
			t.Fatalf("n0 has accepted the vote when it shouldn't have after restarting, as we have already pulled the vote from persister")
		}
	})
}
