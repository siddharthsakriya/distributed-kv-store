package grpc

import (
	raftpb "github.com/siddharthsakriya/distributed-kv-store/proto/gen/raft/v1"
	"github.com/siddharthsakriya/distributed-kv-store/raft"
)

func pbRequestVote(a *raft.RequestVoteArgs) *raftpb.RequestVoteRequest {
	return &raftpb.RequestVoteRequest{
		Term:         int64(a.Term),
		CandidateId:  a.CandidateID,
		LastLogIndex: int64(a.LastLogIndex),
		LastLogTerm:  int64(a.LastLogTerm),
	}
}

func raftRequestVote(r *raftpb.RequestVoteRequest) *raft.RequestVoteArgs {
	return &raft.RequestVoteArgs{
		Term:         int(r.Term),
		CandidateID:  r.CandidateId,
		LastLogIndex: int(r.LastLogIndex),
		LastLogTerm:  int(r.LastLogTerm),
	}
}

func pbRequestVoteResponse(a *raft.RequestVoteReply) *raftpb.RequestVoteResponse {
	return &raftpb.RequestVoteResponse{
		Term:        int64(a.Term),
		VoteGranted: a.VoteGranted,
	}
}

func raftRequestVoteResponse(r *raftpb.RequestVoteResponse) *raft.RequestVoteReply {
	return &raft.RequestVoteReply{
		Term:        int(r.Term),
		VoteGranted: r.VoteGranted,
	}
}

func pbAppendEntries(a *raft.AppendEntriesArgs) *raftpb.AppendEntriesRequest {}

func raftAppendEntries(r *raftpb.AppendEntriesRequest) *raft.AppendEntriesArgs {}

func pbAppendEntriesResponse(a *raft.AppendEntriesReply) *raftpb.AppendEntriesResponse {}

func raftAppendEntriesResponse(r *raftpb.AppendEntriesResponse) *raft.AppendEntriesReply {}

func pbLogEntry(a *raft.LogEntry) *raftpb.LogEntry {}

func raftLogEntry(r *raftpb.LogEntry) *raft.LogEntry {}
