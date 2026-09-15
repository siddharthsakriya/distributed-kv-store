package grpctransport

import (
	"sync"

	raftpb "github.com/siddharthsakriya/distributed-kv-store/proto/gen/raft/v1"
	"google.golang.org/grpc"
)

type Transport struct {
	mu      sync.Mutex
	addrs   map[string]string
	clients map[string]raftpb.RaftServiceClient
	conns   map[string]*grpc.ClientConn
}

func New(addrs map[string]string) *Transport {
	return &Transport{
		addrs:   addrs,
		clients: make(map[string]raftpb.RaftServiceClient),
		conns:   map[string]*grpc.ClientConn{},
	}
}
