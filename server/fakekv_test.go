package server

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestFakeKV(t *testing.T) {
	t.Run("test basic put", func(t *testing.T) {
		fakeStore := FakeKV{
			Store: make(map[string][]byte),
		}
		payload, _ := json.Marshal(Command{
			Action: "PUT",
			Key:    "x",
			Val:    []byte("1"),
		})
		res := fakeStore.Apply([]byte(payload))
		if !bytes.Equal(res, []byte("OK")) || !bytes.Equal(fakeStore.Store["x"], []byte("1")) {
			t.Fatalf("Basic put test failed")
		}
	})

	t.Run("test basic get", func(t *testing.T) {
		fakeStore := FakeKV{
			Store: make(map[string][]byte),
		}
		fakeStore.Store["x"] = []byte("1")
		payload, _ := json.Marshal(Command{
			Action: "GET",
			Key:    "x",
			Val:    nil,
		})
		res := fakeStore.Apply([]byte(payload))
		if !bytes.Equal(res, []byte("1")) {
			t.Fatalf("Basic get test failed")
		}
	})

	t.Run("test basic delete", func(t *testing.T) {
		fakeStore := FakeKV{
			Store: make(map[string][]byte),
		}
		fakeStore.Store["x"] = []byte("1")
		payload, _ := json.Marshal(Command{
			Action: "DELETE",
			Key:    "x",
			Val:    nil,
		})
		res := fakeStore.Apply([]byte(payload))
		_, ok := fakeStore.Store["x"]
		if !bytes.Equal(res, []byte("OK")) || ok {
			t.Fatalf("Basic delete test failed")
		}
	})

	t.Run("test random ah command which doesnt exist", func(t *testing.T) {
		fakeStore := FakeKV{
			Store: make(map[string][]byte),
		}
		payload, _ := json.Marshal(Command{
			Action: "RANDOM AH COMMAND",
			Key:    "d",
			Val:    nil,
		})
		res := fakeStore.Apply([]byte(payload))
		if !bytes.Equal(res, []byte("action not allowed")) {
			t.Fatalf("Basic delete test failed")
		}
	})
}
