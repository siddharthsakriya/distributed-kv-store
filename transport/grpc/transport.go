package grpctransport

import (
	"context"
	"fmt"
	"sync"

	raftpb "github.com/siddharthsakriya/distributed-kv-store/proto/gen/raft/v1"
	"github.com/siddharthsakriya/distributed-kv-store/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Transport struct {
	mu      sync.Mutex
	addrs   map[string]string
	clients map[string]raftpb.RaftServiceClient
	conns   map[string]*grpc.ClientConn
}

var _ raft.Transport = (*Transport)(nil)

func New(addrs map[string]string) *Transport {
	return &Transport{
		addrs:   addrs,
		clients: make(map[string]raftpb.RaftServiceClient),
		conns:   make(map[string]*grpc.ClientConn),
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

func (t *Transport) SendRequestVote(peerID string, args *raft.RequestVoteArgs) (*raft.RequestVoteReply, error) {
	client, err := t.client(peerID)
	if err != nil {
		return nil, err
	}
	response, err := client.RequestVote(context.Background(), pbRequestVote(args))
	if err != nil {
		return nil, err
	}
	return raftRequestVoteResponse(response), nil
}

func (t *Transport) SendAppendEntries(peerID string, args *raft.AppendEntriesArgs) (*raft.AppendEntriesReply, error) {
	client, err := t.client(peerID)
	if err != nil {
		return nil, err
	}
	response, err := client.AppendEntries(context.Background(), pbAppendEntries(args))
	if err != nil {
		return nil, err
	}
	return raftAppendEntriesResponse(response), nil
}

func (t *Transport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, conn := range t.conns {
		_ = conn.Close()
	}
	t.conns = make(map[string]*grpc.ClientConn)
	t.clients = make(map[string]raftpb.RaftServiceClient)
	return nil
}
