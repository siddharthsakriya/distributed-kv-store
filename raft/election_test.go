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

func TestHandleRequestVote(t *testing.T) {
	t.Run("reject stale term", func(t *testing.T) {
		n := NewNode(Config{ID: "n1"})
		n.currentTerm = 5

		reply := n.HandleRequestVote(&RequestVoteArgs{
			Term:        3,
			CandidateID: "n2",
		})

		if reply.VoteGranted {
			t.Error("VoteGranted = true, want false (candidate term is stale)")
		}
		if reply.Term != 5 {
			t.Errorf("reply.Term = %d, want 5", reply.Term)
		}
		if n.votedFor != "" {
			t.Errorf("votedFor = %q, want empty (must not vote for a stale candidate)", n.votedFor)
		}
	})

	t.Run("grant vote to fresh candidate", func(t *testing.T) {
		n := NewNode(Config{ID: "n1"})
		n.currentTerm = 1

		reply := n.HandleRequestVote(&RequestVoteArgs{
			Term:         1,
			CandidateID:  "n2",
			LastLogIndex: 0,
			LastLogTerm:  0,
		})

		if !reply.VoteGranted {
			t.Error("VoteGranted = false, want true")
		}
		if n.votedFor != "n2" {
			t.Errorf("votedFor = %q, want n2 (vote must be recorded)", n.votedFor)
		}
	})

	t.Run("higher term bumps then grants", func(t *testing.T) {
		n := NewNode(Config{ID: "n1"})
		n.currentTerm = 1
		n.votedFor = "old"
		n.role = Candidate

		reply := n.HandleRequestVote(&RequestVoteArgs{
			Term:        3,
			CandidateID: "n2",
		})

		if !reply.VoteGranted {
			t.Error("VoteGranted = false, want true")
		}
		if n.currentTerm != 3 {
			t.Errorf("currentTerm = %d, want 3 (should adopt higher term)", n.currentTerm)
		}
		if reply.Term != 3 {
			t.Errorf("reply.Term = %d, want 3", reply.Term)
		}
		if n.role != Follower {
			t.Errorf("role = %v, want Follower (higher term must step down)", n.role)
		}
		if n.votedFor != "n2" {
			t.Errorf("votedFor = %q, want n2 (old-term vote cleared, then voted)", n.votedFor)
		}
	})

	t.Run("deny when already voted for another candidate", func(t *testing.T) {
		n := NewNode(Config{ID: "n1"})
		n.currentTerm = 2
		n.votedFor = "nX"

		reply := n.HandleRequestVote(&RequestVoteArgs{
			Term:        2,
			CandidateID: "nY",
		})

		if reply.VoteGranted {
			t.Error("VoteGranted = true, want false (already voted for someone else)")
		}
		if n.votedFor != "nX" {
			t.Errorf("votedFor = %q, want nX (must not change)", n.votedFor)
		}
	})

	t.Run("grant again to same candidate (idempotent)", func(t *testing.T) {
		n := NewNode(Config{ID: "n1"})
		n.currentTerm = 2
		n.votedFor = "n2"

		reply := n.HandleRequestVote(&RequestVoteArgs{
			Term:        2,
			CandidateID: "n2",
		})

		if !reply.VoteGranted {
			t.Error("VoteGranted = false, want true (re-request from same candidate)")
		}
	})

	t.Run("deny when candidate log is behind", func(t *testing.T) {
		n := NewNode(Config{ID: "n1"})
		n.currentTerm = 2
		n.log = []LogEntry{{Index: 1, Term: 2}}

		reply := n.HandleRequestVote(&RequestVoteArgs{
			Term:         2,
			CandidateID:  "n2",
			LastLogIndex: 5,
			LastLogTerm:  1,
		})

		if reply.VoteGranted {
			t.Error("VoteGranted = true, want false (candidate log is stale despite being longer)")
		}
		if n.votedFor != "" {
			t.Errorf("votedFor = %q, want empty", n.votedFor)
		}
	})
}
