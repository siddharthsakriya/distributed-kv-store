package cluster

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		yaml := `
nodes:
  - id: n0
    addr: localhost:9000
  - id: n1
    addr: localhost:9001
`
		path := filepath.Join(t.TempDir(), "cluster.yaml")
		if err := os.WriteFile(path, []byte(yaml), 0o600); err != nil {
			t.Fatal(err)
		}

		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}
		if len(cfg.Nodes) != 2 {
			t.Fatalf("expected 2 nodes, got %d", len(cfg.Nodes))
		}
		if cfg.Nodes[0].ID != "n0" || cfg.Nodes[0].Addr != "localhost:9000" {
			t.Fatalf("node 0 parsed wrong: %+v", cfg.Nodes[0])
		}
	})

	t.Run("missing file errors", func(t *testing.T) {
		if _, err := Load("does-not-exist.yaml"); err == nil {
			t.Fatal("expected error for missing file, got nil")
		}
	})

	t.Run("empty config errors", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "empty.yaml")
		os.WriteFile(path, []byte("nodes: []\n"), 0o600)
		if _, err := Load(path); err == nil {
			t.Fatal("expected error for empty node list, got nil")
		}
	})
}

func TestView(t *testing.T) {
	cfg := Config{Nodes: []NodeConfig{
		{ID: "n0", Addr: "localhost:9000"},
		{ID: "n1", Addr: "localhost:9001"},
		{ID: "n2", Addr: "localhost:9002"},
	}}

	t.Run("splits self from peers", func(t *testing.T) {
		myAddr, peers, err := cfg.View("n0")
		if err != nil {
			t.Fatalf("View failed: %v", err)
		}
		if myAddr != "localhost:9000" {
			t.Fatalf("expected own addr localhost:9000, got %q", myAddr)
		}
		want := map[string]string{"n1": "localhost:9001", "n2": "localhost:9002"}
		if !reflect.DeepEqual(peers, want) {
			t.Fatalf("peers wrong:\n got %v\n want %v", peers, want)
		}
	})

	t.Run("unknown id errors", func(t *testing.T) {
		if _, _, err := cfg.View("nope"); err == nil {
			t.Fatal("expected error for unknown id, got nil")
		}
	})
}
