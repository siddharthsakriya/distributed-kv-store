package raft

type Persister interface {
	Save(term int, votedFor string, log []LogEntry) error
	Load() (term int, votedFor string, log []LogEntry)
}
