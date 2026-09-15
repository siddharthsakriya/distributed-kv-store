package grpctransport

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

func pbAppendEntries(a *raft.AppendEntriesArgs) *raftpb.AppendEntriesRequest {
	entries := make([]*raftpb.LogEntry, len(a.Entries))
	for i, e := range a.Entries {
		entries[i] = pbLogEntry(e)
	}
	return &raftpb.AppendEntriesRequest{
		Term:         int64(a.Term),
		LeaderId:     a.LeaderID,
		PrevLogIndex: int64(a.PrevLogIndex),
		PrevLogTerm:  int64(a.PrevLogTerm),
		Entries:      entries,
		LeaderCommit: int64(a.LeaderCommit),
	}
}

func raftAppendEntries(r *raftpb.AppendEntriesRequest) *raft.AppendEntriesArgs {
	entries := make([]raft.LogEntry, len(r.Entries))
	for i, e := range r.Entries {
		entries[i] = raftLogEntry(e)
	}
	return &raft.AppendEntriesArgs{
		Term:         int(r.Term),
		LeaderID:     r.LeaderId,
		PrevLogIndex: int(r.PrevLogIndex),
		PrevLogTerm:  int(r.PrevLogTerm),
		Entries:      entries,
		LeaderCommit: int(r.LeaderCommit),
	}
}

func pbAppendEntriesResponse(a *raft.AppendEntriesReply) *raftpb.AppendEntriesResponse {
	return &raftpb.AppendEntriesResponse{
		Term:    int64(a.Term),
		Success: a.Success,
	}
}

func raftAppendEntriesResponse(r *raftpb.AppendEntriesResponse) *raft.AppendEntriesReply {
	return &raft.AppendEntriesReply{
		Term:    int(r.Term),
		Success: r.Success,
	}
}

func pbLogEntry(a raft.LogEntry) *raftpb.LogEntry {
	return &raftpb.LogEntry{
		Index:   int64(a.Index),
		Term:    int64(a.Term),
		Command: a.Command,
	}
}

func raftLogEntry(r *raftpb.LogEntry) raft.LogEntry {
	return raft.LogEntry{
		Index:   int(r.Index),
		Term:    int(r.Term),
		Command: r.Command,
	}
}
