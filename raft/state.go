package raft

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
