package grpctransport

import (
	"fmt"
	"sync"

	raftpb "github.com/siddharthsakriya/distributed-kv-store/proto/gen/raft/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (t *Transport) client(peerID string) (raftpb.RaftServiceClient, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	client, ok := t.clients[peerID]
	// hit
	if ok {
		return client, nil
	}
	addr, ok := t.addrs[peerID]
	if !ok {
		return nil, fmt.Errorf("unknown peer: %s", peerID)
	}

	//dial and cache it
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := raftpb.NewRaftServiceClient(conn)
	t.conns[peerID] = conn
	t.clients[peerID] = c
	return c, nil
}
