package raft

import "testing"

type stubHandler struct {
	voteReply   *RequestVoteReply
	appendReply *AppendEntriesReply
}

func (s *stubHandler) HandleRequestVote(args *RequestVoteArgs) *RequestVoteReply {
	return s.voteReply
}

func (s *stubHandler) HandleAppendEntries(args *AppendEntriesArgs) *AppendEntriesReply {
	return s.appendReply
}

func TestFakeTransport_RequestVoteRoutesToCorrectPeer(t *testing.T) {
	ft := NewFakeTransport()
	caller := ft.Endpoint("caller")

	handlerA := &stubHandler{voteReply: &RequestVoteReply{Term: 1, VoteGranted: true}}
	handlerB := &stubHandler{voteReply: &RequestVoteReply{Term: 2, VoteGranted: false}}
	ft.Register("A", handlerA)
	ft.Register("B", handlerB)

	replyA, err := caller.SendRequestVote("A", &RequestVoteArgs{})
	if err != nil {
		t.Fatalf("send to A: unexpected error: %v", err)
	}
	if replyA != handlerA.voteReply {
		t.Errorf("send to A: got %+v, want handler A's reply %+v", replyA, handlerA.voteReply)
	}

	replyB, err := caller.SendRequestVote("B", &RequestVoteArgs{})
	if err != nil {
		t.Fatalf("send to B: unexpected error: %v", err)
	}
	if replyB != handlerB.voteReply {
		t.Errorf("send to B: got %+v, want handler B's reply %+v", replyB, handlerB.voteReply)
	}
}

func TestFakeTransport_AppendEntriesRoutesToCorrectPeer(t *testing.T) {
	ft := NewFakeTransport()
	caller := ft.Endpoint("caller")
	handlerA := &stubHandler{appendReply: &AppendEntriesReply{Term: 1, Success: true}}
	ft.Register("A", handlerA)

	reply, err := caller.SendAppendEntries("A", &AppendEntriesArgs{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != handlerA.appendReply {
		t.Errorf("got %+v, want %+v", reply, handlerA.appendReply)
	}
}

func TestFakeTransport_DisconnectBlocksDelivery(t *testing.T) {
	ft := NewFakeTransport()
	caller := ft.Endpoint("caller")
	handlerA := &stubHandler{voteReply: &RequestVoteReply{Term: 1, VoteGranted: true}}
	ft.Register("A", handlerA)

	ft.Disconnect("A")

	if _, err := caller.SendRequestVote("A", &RequestVoteArgs{}); err == nil {
		t.Fatal("expected error sending to disconnected peer, got nil")
	}
}

func TestFakeTransport_ConnectRestoresDelivery(t *testing.T) {
	ft := NewFakeTransport()
	caller := ft.Endpoint("caller")
	handlerA := &stubHandler{voteReply: &RequestVoteReply{Term: 1, VoteGranted: true}}
	ft.Register("A", handlerA)

	ft.Disconnect("A")
	ft.Connect("A")

	reply, err := caller.SendRequestVote("A", &RequestVoteArgs{})
	if err != nil {
		t.Fatalf("unexpected error after reconnect: %v", err)
	}
	if reply != handlerA.voteReply {
		t.Errorf("got %+v, want %+v", reply, handlerA.voteReply)
	}
}

func TestFakeTransport_UnregisteredPeerErrors(t *testing.T) {
	ft := NewFakeTransport()
	caller := ft.Endpoint("caller")

	if _, err := caller.SendRequestVote("ghost", &RequestVoteArgs{}); err == nil {
		t.Fatal("expected error sending to unregistered peer, got nil")
	}
}
