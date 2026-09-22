package raft

import (
	"encoding/json/v2"
	"os"
	"sync"
)

type FilePersister struct {
	mu   sync.Mutex
	path string
}

type persistedState struct {
	Term     int
	VotedFor string
	Log      []LogEntry
}

func NewFilePersister(path string) *FilePersister {
	return &FilePersister{
		path: path,
	}
}

func (fp *FilePersister) Save(term int, votedFor string, log []LogEntry) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	// encode the info
	data, err := json.Marshal(persistedState{
		Term:     term,
		VotedFor: votedFor,
		Log:      log,
	})

	if err != nil {
		return err
	}

	tmpPath := fp.path + ".tmp"
	f, err := os.Create(tmpPath)

	if err != nil {
		return err
	}

	_, err = f.Write(data)

	if err != nil {
		f.Close()
		return err
	}

	err = f.Sync()

	if err != nil {
		f.Close()
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	// swap old file with new
	return os.Rename(tmpPath, fp.path)
}

func (fp *FilePersister) Load() (term int, votedFor string, log []LogEntry) {

}
