package raft

import "time"

// NodeState represents the state of a Raft node
type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

func (s NodeState) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unknown"
	}
}

// LogEntry represents a single entry in the Raft log
type LogEntry struct {
	Term  int64
	Index int64
	Data  []byte
}

// RaftNode represents a single node in the Raft cluster
type RaftNode struct {
	// Persistent state on all servers
	CurrentTerm int64
	VotedFor    int
	Log         []*LogEntry

	// Volatile state on all servers
	CommitIndex int64
	LastApplied int64
	State       NodeState
	NodeID      int

	// Volatile state on leaders
	NextIndex  map[int]int64
	MatchIndex map[int]int64

	// Timers
	electionTimer   *time.Timer
	heartbeatTicker *time.Ticker

	// Cluster configuration
	ClusterNodes map[int]string // nodeID -> address

	// Channel for applying commands to state machine
	ApplyChan chan *LogEntry
}

// RPCRequest and Response types
type RequestVoteRequest struct {
	Term         int64
	CandidateID  int
	LastLogIndex int64
	LastLogTerm  int64
}

type RequestVoteResponse struct {
	Term        int64
	VoteGranted bool
}

type AppendEntriesRequest struct {
	Term         int64
	LeaderID     int
	PrevLogIndex int64
	PrevLogTerm  int64
	Entries      []*LogEntry
	LeaderCommit int64
}

type AppendEntriesResponse struct {
	Term    int64
	Success bool
}
