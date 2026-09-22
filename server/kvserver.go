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
var ErrLostLeadership = errors.New("lost leadership, entry overwritten")

type applyResult struct {
	value []byte
	err   error
}

type waiter struct {
	term int
	ch   chan applyResult
}

type KVServer struct {
	node    *raft.Node
	sm      StateMachine
	applyCh chan raft.ApplyMsg
	mu      sync.Mutex
	waiters map[int]*waiter
}

func NewKVServer(node *raft.Node, sm StateMachine, applyChan chan raft.ApplyMsg) *KVServer {
	return &KVServer{
		node:    node,
		sm:      sm,
		applyCh: applyChan,
		waiters: make(map[int]*waiter),
	}
}

func (s *KVServer) RunApplyLoop() {
	for msg := range s.applyCh {
		res := s.sm.Apply(msg.Command)
		log.Printf("applied idx=%d cmd=%s", msg.Index, msg.Command)
		s.mu.Lock()
		w, ok := s.waiters[msg.Index]
		if ok {
			if w.term != msg.Term {
				w.ch <- applyResult{
					err: ErrLostLeadership,
				}
			} else {
				w.ch <- applyResult{
					value: res,
				}
			}
		}
		delete(s.waiters, msg.Index)
		s.mu.Unlock()
	}
}

func (s *KVServer) Submit(cmd []byte) ([]byte, error) {
	s.mu.Lock()
	res := s.node.Submit(cmd)
	if !res.IsLeader {
		s.mu.Unlock()
		return nil, ErrNotLeader
	}
	w := &waiter{term: res.Term, ch: make(chan applyResult, 1)}
	s.waiters[res.Index] = w
	s.mu.Unlock()
	select {
	case result := <-w.ch:
		return result.value, result.err
	case <-time.After(2 * time.Second):
		s.mu.Lock()
		delete(s.waiters, res.Index)
		s.mu.Unlock()
		return nil, ErrTimeout
	}
}
