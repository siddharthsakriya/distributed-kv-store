package raft

import (
	"sync"
	"time"
)

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

func (r Role) String() string {
	switch r {
	case Follower:
		return "follower"
	case Candidate:
		return "candidate"
	case Leader:
		return "leader"
	default:
		return "unknown"
	}
}

type LogEntry struct {
	Index   int
	Term    int
	Command []byte
}

type Node struct {
	mu sync.Mutex

	// identity + wiring
	id              string
	peers           []string
	transport       Transport
	persister       Persister
	lastHeard       time.Time
	electionTimeout time.Duration

	// persistent
	currentTerm int
	votedFor    string
	log         []LogEntry

	// volatile (all servers)
	role        Role
	commitIndex int
	lastApplied int

	// volatile (leader only)
	nextIndex  map[string]int
	matchIndex map[string]int
}

type Config struct {
	ID        string
	Peers     []string
	Transport Transport
	Persister Persister
}

func NewNode(cfg Config) *Node {
	var currentTerm int
	var votedFor string
	var log []LogEntry

	if cfg.Persister != nil {
		currentTerm, votedFor, log = cfg.Persister.Load()
	}

	return &Node{
		id:          cfg.ID,
		peers:       cfg.Peers,
		transport:   cfg.Transport,
		persister:   cfg.Persister,
		currentTerm: currentTerm,
		votedFor:    votedFor,
		log:         log,
	}
}
