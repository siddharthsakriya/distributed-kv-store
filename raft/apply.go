package raft

type ApplyMsg struct {
	Index   int
	Command []byte
}
