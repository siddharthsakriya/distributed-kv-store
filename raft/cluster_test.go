package raft

import (
	"fmt"
	"testing"
	"time"
)

// build fake transport, and create fake cluster of nodes
func makeCluster(n int) ([]*Node, *FakeTransport) {
	ft := NewFakeTransport()

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
		node := NewNode(Config{ID: id, Peers: peers, Transport: ft.Endpoint(id)})
		nodes[i] = node
		ft.Register(id, node)
		ft.Connect(id)
	}
	return nodes, ft
}

func countLeaders(nodes []*Node) (int, *Node) {
	count := 0
	var leader *Node
	for _, node := range nodes {
		node.mu.Lock()
		if node.role == Leader {
			count++
			leader = node
		}
		node.mu.Unlock()
	}
	return count, leader
}

func waitForOneLeader(t *testing.T, nodes []*Node, timeout time.Duration) *Node {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if count, leader := countLeaders(nodes); count == 1 {
			return leader
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("no single leader elected within %v", timeout)
	return nil
}

func TestElectsOneLeader(t *testing.T) {
	nodes, _ := makeCluster(3)
	for _, node := range nodes {
		go node.runElectionTimer()
	}
	waitForOneLeader(t, nodes, 3*time.Second)
}

func TestReelectsAfterLeaderDisconnect(t *testing.T) {
	nodes, ft := makeCluster(3)
	for _, node := range nodes {
		go node.runElectionTimer()
	}
	leader := waitForOneLeader(t, nodes, 3*time.Second)

	ft.Disconnect(leader.id)

	var connected []*Node
	for _, node := range nodes {
		if node != leader {
			connected = append(connected, node)
		}
	}
	waitForOneLeader(t, connected, 3*time.Second)
}

func TestFakeTransport_DisconnectedSenderCannotSend(t *testing.T) {
	ft := NewFakeTransport()
	handlerA := &stubHandler{voteReply: &RequestVoteReply{Term: 1, VoteGranted: true}}
	ft.Register("A", handlerA)

	sender := ft.Endpoint("S")
	ft.Disconnect("S")

	if _, err := sender.SendRequestVote("A", &RequestVoteArgs{}); err == nil {
		t.Fatal("expected error: disconnected sender should not be able to send")
	}
}
