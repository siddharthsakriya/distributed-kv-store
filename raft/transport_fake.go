package raft

import (
	"sync"
)

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
