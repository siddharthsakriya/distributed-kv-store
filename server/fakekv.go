package server

import (
	"encoding/json"
)

type FakeKV struct {
	Store map[string][]byte
}

type Command struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	Val    []byte `json:"val,omitempty"`
}

func (fakeKV *FakeKV) Apply(cmd []byte) []byte {
	var c Command
	err := json.Unmarshal(cmd, &c)
	if err != nil {
		return []byte("bad command")
	}
	return fakeKV.ApplyCommand(c.Action, c.Key, c.Val)
}

func (fakeKV *FakeKV) ApplyCommand(action string, key string, val []byte) []byte {
	switch action {
	case "PUT":
		fakeKV.Store[key] = []byte(val)
		return []byte("OK")
	case "GET":
		value, ok := fakeKV.Store[key]
		if !ok {
			return nil
		}
		return value
	case "DELETE":
		delete(fakeKV.Store, key)
		return []byte("OK")
	default:
		return []byte("action not allowed")
	}
}
