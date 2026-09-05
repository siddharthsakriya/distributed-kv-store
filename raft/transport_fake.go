package raft

import (
	"fmt"
	"sync"
)

type fakeEndpoint struct {
	id  string
	net *FakeTransport
}

func (ft *FakeTransport) Endpoint(id string) *fakeEndpoint {
	return &fakeEndpoint{id: id, net: ft}
}

type FakeTransport struct {
	handlers map[string]RPCHandler
	down     map[string]bool
	mu       sync.Mutex
}

func NewFakeTransport() *FakeTransport {
	return &FakeTransport{
		handlers: make(map[string]RPCHandler),
		down:     make(map[string]bool),
	}
}

var _ Transport = (*fakeEndpoint)(nil)

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

func (e *fakeEndpoint) SendRequestVote(peerID string, args *RequestVoteArgs) (*RequestVoteReply, error) {
	e.net.mu.Lock()
	handler, ok := e.net.handlers[peerID]
	senderDown := e.net.down[e.id]
	targetDown := e.net.down[peerID]
	e.net.mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("peer %s not registered", peerID)
	}
	if senderDown || targetDown {
		return nil, fmt.Errorf("cannot reach %s from %s", peerID, e.id)
	}
	return handler.HandleRequestVote(args), nil
}

func (e *fakeEndpoint) SendAppendEntries(peerID string, args *AppendEntriesArgs) (*AppendEntriesReply, error) {
	e.net.mu.Lock()
	handler, ok := e.net.handlers[peerID]
	senderDown := e.net.down[e.id]
	targetDown := e.net.down[peerID]
	e.net.mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("peer %s not registered", peerID)
	}
	if senderDown || targetDown {
		return nil, fmt.Errorf("cannot reach %s from %s", peerID, e.id)
	}
	return handler.HandleAppendEntries(args), nil
}
