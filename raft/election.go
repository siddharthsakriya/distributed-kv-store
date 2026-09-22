package raft

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"
)

const (
	// hardcoded upper and lower bounds for now according to what is specified in the paper, maybe could generalise :0
	minTimeout = 150 * time.Millisecond
	maxTimeout = 300 * time.Millisecond
)

func (n *Node) HandleRequestVote(args *RequestVoteArgs) *RequestVoteReply {
	n.mu.Lock()
	defer n.mu.Unlock()

	if args.Term < n.currentTerm {
		return &RequestVoteReply{
			Term:        n.currentTerm,
			VoteGranted: false,
		}
	}

	if args.Term > n.currentTerm {
		n.currentTerm = args.Term
		n.votedFor = ""
		n.role = Follower
	}

	if (n.votedFor != "" && n.votedFor != args.CandidateID) || !isLogUpToDate(args.LastLogTerm, args.LastLogIndex, n.lastLogTerm(), n.lastLogIndex()) {
		n.persist()
		return &RequestVoteReply{
			Term:        n.currentTerm,
			VoteGranted: false,
		}
	}

	n.votedFor = args.CandidateID
	log.Printf("[%s] voted for %s in term %d", n.id, args.CandidateID, n.currentTerm)
	n.resetElectionTimer()
	n.persist()

	return &RequestVoteReply{
		Term:        n.currentTerm,
		VoteGranted: true,
	}
}

func (n *Node) runElectionTimer() {
	for {
		time.Sleep(10 * time.Millisecond)

		n.mu.Lock()
		if n.stopped {
			n.mu.Unlock()
			return
		}

		shouldElect := n.role != Leader && time.Since(n.lastHeard) > n.electionTimeout
		n.mu.Unlock()
		if shouldElect {
			n.startElection()
		}
	}
}

func (n *Node) startElection() {
	n.mu.Lock()
	if n.stopped || n.role == Leader || time.Since(n.lastHeard) <= n.electionTimeout {
		n.mu.Unlock()
		return
	}
	n.currentTerm++
	n.votedFor = n.id
	n.role = Candidate
	log.Printf("[%s] starting election for term %d", n.id, n.currentTerm)
	n.resetElectionTimer()
	n.persist()
	term := n.currentTerm
	lastLogIndex := n.lastLogIndex()
	lastLogTerm := n.lastLogTerm()
	n.mu.Unlock()

	globalVoteTally := 1

	for _, peerID := range n.peers {
		go func(peerID string) {
			args := &RequestVoteArgs{
				Term:         term,
				CandidateID:  n.id,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			}

			reply, err := n.transport.SendRequestVote(peerID, args)

			if err != nil {
				return
			}

			n.mu.Lock()
			// guard 1 : peer term ahead of ours
			if reply.Term > n.currentTerm {
				n.role = Follower
				n.currentTerm = reply.Term
				log.Printf("[%s] stepping down to follower, term %d", n.id, n.currentTerm)
				n.votedFor = ""
				n.persist()
				n.mu.Unlock()
				return
			}

			// guard 2 : u are not candidate anymore
			if n.role != Candidate || n.currentTerm != term {
				n.mu.Unlock()
				return
			}

			if reply.VoteGranted {
				globalVoteTally++
				if n.isMajorityVote(globalVoteTally) {
					n.role = Leader
					log.Printf("[%s] became LEADER for term %d", n.id, n.currentTerm)
					// initialise nextIndex and matchIndex maps
					n.constructLeaderMaps()
					go n.runReplication(n.currentTerm)
				}
			}
			n.mu.Unlock()
		}(peerID)
	}
}

/*** helpers ***/

func isLogUpToDate(candidateLastTerm int, candidateLastIndex int, myLastTerm int, myLastIndex int) bool {
	if candidateLastTerm != myLastTerm {
		return candidateLastTerm > myLastTerm
	}
	return candidateLastIndex >= myLastIndex
}

func (n *Node) constructLeaderMaps() {
	n.nextIndex = make(map[string]int)
	n.matchIndex = make(map[string]int)
	lastLogIndex := n.lastLogIndex()
	for _, initPeerID := range n.peers {
		n.nextIndex[initPeerID] = lastLogIndex + 1
		n.matchIndex[initPeerID] = 0
	}
}

func (n *Node) isMajorityVote(globalVoteTally int) bool {
	total := len(n.peers) + 1
	return globalVoteTally*2 > total
}

func (n *Node) lastLogIndex() int {
	if len(n.log) == 0 {
		return 0
	}
	return n.log[len(n.log)-1].Index
}

func (n *Node) lastLogTerm() int {
	if len(n.log) == 0 {
		return 0
	}
	return n.log[len(n.log)-1].Term
}

func (n *Node) persist() {
	if n.persister == nil {
		return
	}
	err := n.persister.Save(n.currentTerm, n.votedFor, n.log)
	if err != nil {
		panic(fmt.Sprintf("[%s] persist failed at term %d: %v", n.id, n.currentTerm, err))
	}
}

func (n *Node) resetElectionTimer() {
	n.lastHeard = time.Now()
	n.electionTimeout = minTimeout + time.Duration(rand.Int64N(int64(maxTimeout-minTimeout)))
}
