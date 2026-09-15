package grpctransport

import (
	"context"

	raftpb "github.com/siddharthsakriya/distributed-kv-store/proto/gen/raft/v1"
	"github.com/siddharthsakriya/distributed-kv-store/raft"
)

type Server struct {
	raftpb.UnimplementedRaftServiceServer
	handler raft.RPCHandler
}

var _ raftpb.RaftServiceServer = (*Server)(nil)

func NewServer(handler raft.RPCHandler) *Server {
	return &Server{handler: handler}
}

func (s *Server) RequestVote(ctx context.Context, req *raftpb.RequestVoteRequest) (*raftpb.RequestVoteResponse, error) {
	response := pbRequestVoteResponse(s.handler.HandleRequestVote(raftRequestVote(req)))
	return response, nil
}

func (s *Server) AppendEntries(ctx context.Context, req *raftpb.AppendEntriesRequest) (*raftpb.AppendEntriesResponse, error) {
	response := pbAppendEntriesResponse(s.handler.HandleAppendEntries(raftAppendEntries(req)))
	return response, nil
}
