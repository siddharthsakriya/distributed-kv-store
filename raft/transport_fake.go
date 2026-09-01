package raft

import (
	"fmt"
	"sync"
)

type FakeTransport struct {
	handlers map[string]RPCHandler
	down     map[string]bool
	mu       sync.Mutex
}

var _ Transport = (*FakeTransport)(nil)

func NewFakeTransport() *FakeTransport {
	return &FakeTransport{
		handlers: make(map[string]RPCHandler),
		down:     make(map[string]bool),
	}
}

func (ft *FakeTransport) Register(peerID string, handler RPCHandler) {
	ft.mu.Lock()
	defer ft.mu.Unlock()
	ft.handlers[peerID] = handler
}

func (ft *FakeTransport) Connect(peerID string) {
	ft.mu.Lock()
	defer ft.mu.Unlock()
	ft.down[peerID] = false
}

func (ft *FakeTransport) Disconnect(peerID string) {
	ft.mu.Lock()
	defer ft.mu.Unlock()
	ft.down[peerID] = true
}

func (ft *FakeTransport) SendRequestVote(peerID string, args *RequestVoteArgs) (*RequestVoteReply, error) {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	handlerVal, handlerExists := ft.handlers[peerID]
	if !handlerExists {
		return nil, fmt.Errorf("peer %s not registered", peerID)
	}

	if ft.down[peerID] {
		return nil, fmt.Errorf("peer %s disconnected", peerID)
	}

	result := handlerVal.HandleRequestVote(args)
	return result, nil
}

func (ft *FakeTransport) SendAppendEntries(peerID string, args *AppendEntriesArgs) (*AppendEntriesReply, error) {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	handlerVal, handlerExists := ft.handlers[peerID]
	if !handlerExists {
		return nil, fmt.Errorf("peer %s not registered", peerID)
	}

	if ft.down[peerID] {
		return nil, fmt.Errorf("peer %s disconnected", peerID)
	}

	result := handlerVal.HandleAppendEntries(args)
	return result, nil
}
