package server

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/siddharthsakriya/distributed-kv-store/raft"
)

var ErrNotLeader = errors.New("not leader")
var ErrTimeout = errors.New("apply timed out")

type KVServer struct {
	node    *raft.Node
	sm      StateMachine
	applyCh chan raft.ApplyMsg
	mu      sync.Mutex
	waiters map[int]chan []byte
}

func NewKVServer(node *raft.Node, sm StateMachine, applyChan chan raft.ApplyMsg) *KVServer {
	return &KVServer{
		node:    node,
		sm:      sm,
		applyCh: applyChan,
		waiters: make(map[int]chan []byte),
	}
}

func (s *KVServer) RunApplyLoop() {
	for msg := range s.applyCh {
		res := s.sm.Apply(msg.Command)
		log.Printf("applied idx=%d cmd=%s", msg.Index, msg.Command)
		s.mu.Lock()
		ch, ok := s.waiters[msg.Index]
		if ok {
			ch <- res
		}
		delete(s.waiters, msg.Index)
		s.mu.Unlock()
	}
}

func (s *KVServer) Submit(cmd []byte) ([]byte, error) {
	res := s.node.Submit(cmd)
	if !res.IsLeader {
		return nil, ErrNotLeader
	}
	ch := make(chan []byte, 1)
	s.mu.Lock()
	s.waiters[res.Index] = ch
	s.mu.Unlock()

	select {
	case result := <-ch:
		return result, nil
	case <-time.After(2 * time.Second):
		s.mu.Lock()
		delete(s.waiters, res.Index)
		s.mu.Unlock()
		return nil, ErrTimeout
	}
}
