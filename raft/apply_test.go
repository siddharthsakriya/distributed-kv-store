package raft

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"
)

// tests adapted from mit 6.824 impl

type applyTracker struct {
	mu           sync.Mutex
	applied      []map[int][]byte
	nextExpected []int
	t            *testing.T
}

func newApplyTracker(t *testing.T, n int) *applyTracker {
	at := &applyTracker{
		t:            t,
		applied:      make([]map[int][]byte, n),
		nextExpected: make([]int, n),
	}
	for i := range at.applied {
		at.applied[i] = make(map[int][]byte)
		at.nextExpected[i] = 1
	}
	return at
}

func (at *applyTracker) record(i int, msg ApplyMsg) {
	at.mu.Lock()
	defer at.mu.Unlock()

	if msg.Index != at.nextExpected[i] {
		at.t.Fatalf("node %d applied index %d out of order (expected %d)",
			i, msg.Index, at.nextExpected[i])
	}
	at.nextExpected[i]++

	for j := range at.applied {
		if prev, ok := at.applied[j][msg.Index]; ok && !bytes.Equal(prev, msg.Command) {
			at.t.Fatalf("apply mismatch at index %d: node %d=%q node %d=%q",
				msg.Index, j, prev, i, msg.Command)
		}
	}
	at.applied[i][msg.Index] = msg.Command
}

// nCommitted: how many nodes have applied this index, and the agreed command.
func (at *applyTracker) nCommitted(index int) (int, []byte) {
	at.mu.Lock()
	defer at.mu.Unlock()
	count := 0
	var cmd []byte
	for i := range at.applied {
		if c, ok := at.applied[i][index]; ok {
			count++
			cmd = c
		}
	}
	return count, cmd
}

func makeApplyCluster(t *testing.T, n int) ([]*Node, *FakeTransport, *applyTracker) {
	ft := NewFakeTransport()
	tracker := newApplyTracker(t, n)

	ids := make([]string, n)
	for i := range ids {
		ids[i] = fmt.Sprintf("n%d", i)
	}

	nodes := make([]*Node, n)
	for i, id := range ids {
		var peers []string
		for _, other := range ids {
			if other != id {
				peers = append(peers, other)
			}
		}
		applyCh := make(chan ApplyMsg)
		node := NewNode(Config{ID: id, Peers: peers, Transport: ft.Endpoint(id), ApplyCh: applyCh})
		nodes[i] = node
		ft.Register(id, node)
		ft.Connect(id)

		go func(idx int, ch chan ApplyMsg) {
			for msg := range ch {
				tracker.record(idx, msg)
			}
		}(i, applyCh)
	}

	for _, node := range nodes {
		go node.Start()
	}
	return nodes, ft, tracker
}

func one(t *testing.T, nodes []*Node, tracker *applyTracker, cmd []byte, expectedServers int) int {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		// find a leader and submit
		index := -1
		for _, node := range nodes {
			if res := node.Submit(cmd); res.IsLeader {
				index = res.Index
				break
			}
		}
		if index != -1 {
			until := time.Now().Add(2 * time.Second)
			for time.Now().Before(until) {
				if count, applied := tracker.nCommitted(index); count >= expectedServers && bytes.Equal(applied, cmd) {
					return index
				}
				time.Sleep(20 * time.Millisecond)
			}
		} else {
			time.Sleep(50 * time.Millisecond) // no leader yet, retry
		}
	}
	t.Fatalf("one(%q): never applied on %d servers", cmd, expectedServers)
	return -1
}

func TestBasicApply(t *testing.T) {
	nodes, _, tracker := makeApplyCluster(t, 3)
	waitForOneLeader(t, nodes, 3*time.Second)

	one(t, nodes, tracker, []byte("cmd1"), len(nodes))
	one(t, nodes, tracker, []byte("cmd2"), len(nodes))
	one(t, nodes, tracker, []byte("cmd3"), len(nodes))
}
