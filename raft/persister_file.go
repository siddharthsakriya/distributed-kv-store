package raft

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

var _ Persister = (*FilePersister)(nil)

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
	err = os.Rename(tmpPath, fp.path)
	if err != nil {
		return err
	}

	dir := filepath.Dir(fp.path)
	d, err := os.Open(dir)
	if err != nil {
		return err
	}

	syncErr := d.Sync()

	err = d.Close()
	if err != nil {
		return err
	}

	return syncErr
}

func (fp *FilePersister) Load() (term int, votedFor string, log []LogEntry) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	data, err := os.ReadFile(fp.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, "", nil
		}
		panic(fmt.Sprintf("load %s: read: %v", fp.path, err))
	}
	var st persistedState
	err = json.Unmarshal(data, &st)
	if err != nil {
		panic(fmt.Sprintf("load %s: corrupt state: %v", fp.path, err))
	}
	return st.Term, st.VotedFor, st.Log
}
