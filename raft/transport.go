package raft

type Transport interface {
	SendRequestVote(peerID string, args *RequestVoteArgs) (*RequestVoteReply, error)
	SendAppendEntries(peerID string, args *AppendEntriesArgs) (*AppendEntriesReply, error)
}

type RPCHandler interface {
	HandleRequestVote(args *RequestVoteArgs) *RequestVoteReply
	HandleAppendEntries(args *AppendEntriesArgs) *AppendEntriesReply
}
