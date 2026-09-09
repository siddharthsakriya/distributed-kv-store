# Distributed KV Store

A distributed, fault-tolerant key-value store built on Raft...

[![CI](https://github.com/siddharthsakriya/distributed-kv-store/actions/workflows/ci.yml/badge.svg)](https://github.com/siddharthsakriya/distributed-kv-store/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/siddharthsakriya/distributed-kv-store)
[![Go Reference](https://pkg.go.dev/badge/github.com/siddharthsakriya/distributed-kv-store.svg)](https://pkg.go.dev/github.com/siddharthsakriya/distributed-kv-store)
![Last commit](https://img.shields.io/github/last-commit/siddharthsakriya/distributed-kv-store)

## KV Store
WIP

## Raft

An implementation of the [Raft consensus algorithm](https://raft.github.io/raft.pdf) in Go. Raft keeps a replicated log consistent across a cluster of nodes so that they agree on the same ordered sequence of commands, even as nodes fail.

The `raft` package is a standalone consensus core: it replicates an opaque log of commands and hands committed entries back to the application. It knows nothing about what those commands mean, which is only the concern of the machine sitting on top of it.

### Status

Core consensus is implemented and tested under the race detector:

- [x] **Leader election** — randomized timeouts, `RequestVote`, term-based step-down
- [x] **Log replication** — `AppendEntries`, per-follower `nextIndex`/`matchIndex`, consistency check and backoff, conflict truncation
- [x] **Commit advancement** — majority `matchIndex` scan with the current-term safety guard (paper §5.4.2 / Figure 8)
- [x] **Apply loop** — committed entries delivered in order, exactly once, on an apply channel
- [x] **Lifecycle** — `Start`/`Stop` with clean goroutine shutdown
- [x] Persistence crash-recovery and Figure 8 safety tests (logic is in place, dedicated tests pending)
- [ ] Additional failure-path tests (partition, kill-leader-mid-flight, log repair)
- [ ] A real network transport (only an in-memory `FakeTransport` exists, used for tests)

### Design

The package talks to the outside world through three interfaces, so it stays decoupled from transport, storage, and the state machine:

| Interface / type | Role |
|------------------|------|
| `Transport`      | Sends `RequestVote` / `AppendEntries` RPCs to peers |
| `Persister`      | Saves and loads persistent state (term, vote, log) |
| `ApplyCh`        | Channel the node pushes committed `ApplyMsg`'s onto |

A node is driven by two touchpoints:

- **`Submit(command)`** — the application hands in a command; the leader appends it to its log and returns an `{Index, Term, IsLeader}` receipt. It does **not** block for commit.
- **`ApplyCh`** — the node delivers committed commands (in index order) as `ApplyMsg{Index, Command}`. The application watches this to learn what has committed.

### Usage

```go
cfg := raft.Config{
    ID:        "n0",
    Peers:     []string{"n1", "n2"}, // other node ids
    Transport: myTransport,          // implements raft.Transport
    Persister: myPersister,          // implements raft.Persister (may be nil)
    ApplyCh:   applyCh,              // chan raft.ApplyMsg (may be nil)
}

node := raft.NewNode(cfg)
node.Start()
defer node.Stop()

// on the leader:
res := node.Submit([]byte("some command"))
if res.IsLeader {
    // watch applyCh for res.Index to know it committed
}

// drain committed entries:
for msg := range applyCh {
    // apply msg.Command to your state machine
}
```

Log indices are 1-based (index 0 means "before the log"), matching the paper.

### Package layout

| File | Contents |
|------|----------|
| `state.go`       | `Node`, `Config`, `LogEntry`, `NewNode` |
| `election.go`    | election timer, `RequestVote` handling, log helpers |
| `replication.go` | `Submit`, replication loop, `AppendEntries` handling |
| `apply.go`       | apply loop, `ApplyMsg` |
| `lifecycle.go`   | `Start` / `Stop` |
| `rpc.go`         | RPC argument / reply structs |
| `transport.go`   | `Transport` / `RPCHandler` interfaces |
| `persister.go`   | `Persister` interface |

### Development

```bash
make test    # run tests
make race    # run tests with the race detector
make vet     # go vet
make build   # compile
```

CI (build, vet, gofmt check, and `go test -race`) runs on every push and pull request.

### Reference

- [In Search of an Understandable Consensus Algorithm (Extended Version)](https://raft.github.io/raft.pdf) — Ongaro & Ousterhout
