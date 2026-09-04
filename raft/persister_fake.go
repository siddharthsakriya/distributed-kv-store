package raft

import (
	"sync"
)

type FakePersister struct {
	mu       sync.Mutex
	term     int
	votedFor string
	log      []LogEntry
}

var _ Persister = (*FakePersister)(nil)

func (fp *FakePersister) Save(term int, votedFor string, log []LogEntry) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	fp.term = term
	fp.votedFor = votedFor
	fp.log = log

	return nil
}

func (fp *FakePersister) Load() (term int, votedFor string, log []LogEntry) {

	fp.mu.Lock()
	defer fp.mu.Unlock()

	return fp.term, fp.votedFor, fp.log
}
