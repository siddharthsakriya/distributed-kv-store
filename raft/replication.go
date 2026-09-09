package raft

import (
	"time"
)

const (
	heartBeatInterval = 50 * time.Millisecond
)

type SubmitResult struct {
	Index    int
	Term     int
	IsLeader bool
}

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

	prevLogTerm := n.entryTerm(args.PrevLogIndex)
	if prevLogTerm != args.PrevLogTerm {
		return &AppendEntriesReply{
			Term:    n.currentTerm,
			Success: false,
		}
	}

	// entries should slot in contiguously due to the term check (entryTerm when idx doesnt exist in log entries, thus send 0)
	for _, entry := range args.Entries {
		idx := entry.Index
		lastLogIndex := n.lastLogIndex()
		if idx > lastLogIndex {
			n.log = append(n.log, entry)
		} else if n.entryTerm(entry.Index) != entry.Term {
			n.log = n.log[:entry.Index-1]
			n.log = append(n.log, entry)
		}
	}

	if args.LeaderCommit > n.commitIndex {
		n.commitIndex = min(args.LeaderCommit, args.PrevLogIndex+len(args.Entries))
		n.applyCond.Signal()
	}

	n.persist()

	return &AppendEntriesReply{
		Term:    n.currentTerm,
		Success: true,
	}
}

func (n *Node) runReplication() {
	for {
		n.mu.Lock()
		role := n.role
		currentTerm := n.currentTerm
		stopped := n.stopped
		n.mu.Unlock()
		if stopped || role != Leader {
			return
		}
		for _, peerID := range n.peers {
			go func(peerID string) {
				n.mu.Lock()
				peerNextIndex := n.nextIndex[peerID]
				prevLogIndex := peerNextIndex - 1
				prevLogTerm := n.entryTerm(prevLogIndex)
				entries := n.entriesFrom(peerNextIndex)
				leaderCommit := n.commitIndex
				n.mu.Unlock()

				args := &AppendEntriesArgs{
					Term:         currentTerm,
					LeaderID:     n.id,
					PrevLogIndex: prevLogIndex,
					PrevLogTerm:  prevLogTerm,
					Entries:      entries,
					LeaderCommit: leaderCommit,
				}

				reply, err := n.transport.SendAppendEntries(peerID, args)

				if err != nil {
					return
				}

				n.mu.Lock()
				defer n.mu.Unlock()

				// handle stale response
				if n.role != Leader || n.currentTerm != currentTerm || n.stopped {
					return
				}

				if reply.Term > n.currentTerm {
					n.role = Follower
					n.currentTerm = reply.Term
					n.votedFor = ""
					n.persist()
					return
				}

				if reply.Success {
					newMatch := prevLogIndex + len(entries)
					if newMatch > n.matchIndex[peerID] {
						n.matchIndex[peerID] = newMatch
						n.nextIndex[peerID] = newMatch + 1
						n.advanceCommitIndex()
					}
				} else {
					if n.nextIndex[peerID] > 1 {
						// try lower index in next iter
						n.nextIndex[peerID]--
					}
				}
			}(peerID)
		}
		time.Sleep(heartBeatInterval)
	}
}

func (n *Node) Submit(command []byte) *SubmitResult {
	n.mu.Lock()
	defer n.mu.Unlock()
	isLeader := n.role == Leader
	if !isLeader {
		return &SubmitResult{
			IsLeader: false,
		}
	}
	index := n.lastLogIndex() + 1
	term := n.currentTerm
	entry := LogEntry{
		Index:   index,
		Term:    term,
		Command: command,
	}
	n.log = append(n.log, entry)
	n.persist()

	return &SubmitResult{
		Index:    index,
		Term:     term,
		IsLeader: true,
	}
}

/*** Helpers ***/
func (n *Node) entryTerm(idx int) int {
	return n.entryAt(idx).Term
}

func (n *Node) entryAt(idx int) LogEntry {
	if idx <= 0 || idx > n.lastLogIndex() {
		return LogEntry{}
	}
	return n.log[idx-1]
}

func (n *Node) entriesFrom(idx int) []LogEntry {
	if idx <= 1 {
		return n.log
	}
	return n.log[idx-1:]
}

func (n *Node) advanceCommitIndex() {

	// top down (trying highest possible)
	for idx := n.lastLogIndex(); idx > n.commitIndex; idx-- {
		// ensure we only commit entries from our term
		if n.entryTerm(idx) != n.currentTerm {
			continue
		}

		count := 1
		for _, peerID := range n.peers {
			if n.matchIndex[peerID] >= idx {
				count++
			}
		}
		if n.isMajorityVote(count) {
			n.commitIndex = idx
			n.applyCond.Signal()
			break
		}
	}
}
