package raft

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
	n.persist()

	return &RequestVoteReply{
		Term:        n.currentTerm,
		VoteGranted: true,
	}
}

/*** helpers for handlerequest to vote method  ***/

func isLogUpToDate(candidateLastTerm int, candidateLastIndex int, myLastTerm int, myLastIndex int) bool {
	if candidateLastTerm != myLastTerm {
		return candidateLastTerm > myLastTerm
	}
	return candidateLastIndex >= myLastIndex
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
	n.persister.Save(n.currentTerm, n.votedFor, n.log)
}
