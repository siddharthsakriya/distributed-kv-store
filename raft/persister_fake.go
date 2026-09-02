package raft

type FakePersister struct {
	term     int
	votedFor string
	log      []LogEntry
}

var _ Persister = (*FakePersister)(nil)

func (fp *FakePersister) Save(term int, votedFor string, log []LogEntry) error {
	fp.term = term
	fp.votedFor = votedFor
	fp.log = log
	return nil
}

func (fp *FakePersister) Load() (term int, votedFor string, log []LogEntry) {
	return fp.term, fp.votedFor, fp.log
}
