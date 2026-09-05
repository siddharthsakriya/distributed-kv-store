package raft

import (
	"time"
)

const (
	heartBeatInterval = 50 * time.Millisecond
)

func (n *Node) HandleAppendEntries(args *AppendEntriesArgs) *AppendEntriesReply {
	n.mu.Lock()
	defer n.mu.Unlock()

	if args.Term < n.currentTerm {
		return &AppendEntriesReply{
			Term:    n.currentTerm,
			Success: false,
		}
	}

	if args.Term > n.currentTerm {
		n.currentTerm = args.Term
		n.votedFor = ""
	}

	n.role = Follower
	n.resetElectionTimer()
	n.persist()

	return &AppendEntriesReply{
		Term:    n.currentTerm,
		Success: true,
	}
}

func (n *Node) runHeartbeats() {
	for {
		n.mu.Lock()
		role := n.role
		currentTerm := n.currentTerm
		n.mu.Unlock()
		if role != Leader {
			return
		}
		for _, peerID := range n.peers {
			go func(peerID string) {
				args := &AppendEntriesArgs{
					Term:     currentTerm,
					LeaderID: n.id,
				}
				reply, err := n.transport.SendAppendEntries(peerID, args)

				if err != nil {
					return
				}

				n.mu.Lock()
				if reply.Term > n.currentTerm {
					n.role = Follower
					n.currentTerm = reply.Term
					n.votedFor = ""
					n.persist()
				}
				n.mu.Unlock()
			}(peerID)
		}
		time.Sleep(heartBeatInterval)
	}
}
