package raft

func (n *Node) HandleAppendEntries(args *AppendEntriesArgs) *AppendEntriesReply {
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

	return
}
