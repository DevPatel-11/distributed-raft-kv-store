package raft

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// NewRaftNode creates a new Raft node
func NewRaftNode(nodeID int, clusterNodes map[int]string) *RaftNode {
	node := &RaftNode{
		NodeID:      nodeID,
		CurrentTerm: 0,
		VotedFor:    -1,
		Log:         make([]*LogEntry, 0),
		State:       Follower,
		CommitIndex: 0,
		LastApplied: 0,
		NextIndex:   make(map[int]int64),
		MatchIndex:  make(map[int]int64),
		ClusterNodes: clusterNodes,
		ApplyChan:   make(chan *LogEntry, 100),
	}
	return node
}

// Start starts the Raft node
func (n *RaftNode) Start() {
	go n.eventLoop()
}

// eventLoop is the main event loop for the Raft node
func (n *RaftNode) eventLoop() {
	n.resetElectionTimer()

	for {
		select {
		case <-n.electionTimer.C:
			n.handleElectionTimeout()
		}
	}
}

// resetElectionTimer resets the election timer with a random timeout
func (n *RaftNode) resetElectionTimer() {
	if n.electionTimer != nil {
		n.electionTimer.Stop()
	}
	timeout := time.Duration(150+rand.Intn(150)) * time.Millisecond
	n.electionTimer = time.AfterFunc(timeout, func() {
		n.handleElectionTimeout()
	})
}

// handleElectionTimeout handles election timeout
func (n *RaftNode) handleElectionTimeout() {
	switch n.State {
	case Follower, Candidate:
		n.startElection()
	case Leader:
		n.resetElectionTimer()
	}
}

// startElection starts a new election
func (n *RaftNode) startElection() {
	n.CurrentTerm++
	n.State = Candidate
	n.VotedFor = n.NodeID

	fmt.Printf("[Node %d] Starting election for term %d\n", n.NodeID, n.CurrentTerm)

	votesNeeded := (len(n.ClusterNodes) / 2) + 1
	votesReceived := 1 // vote for self
	var mu sync.Mutex
	var wg sync.WaitGroup

	for nodeID := range n.ClusterNodes {
		if nodeID == n.NodeID {
			continue
		}

		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// In real implementation, make RPC call here
			// For now, simulate vote
			mu.Lock()
			votesReceived++
			mu.Unlock()
		}(nodeID)
	}

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	if votesReceived >= votesNeeded {
		n.becomeLeader()
	} else {
		n.resetElectionTimer()
	}
}

// becomeLeader transitions the node to leader state
func (n *RaftNode) becomeLeader() {
	n.State = Leader
	fmt.Printf("[Node %d] Became leader for term %d\n", n.NodeID, n.CurrentTerm)

	// Initialize leader state
	for nodeID := range n.ClusterNodes {
		n.NextIndex[nodeID] = int64(len(n.Log))
		n.MatchIndex[nodeID] = 0
	}
}

// AppendLog appends a new entry to the log
func (n *RaftNode) AppendLog(data []byte) {
	entry := &LogEntry{
		Term:  n.CurrentTerm,
		Index: int64(len(n.Log)),
		Data:  data,
	}
	n.Log = append(n.Log, entry)
}

// GetState returns the current state of the node
func (n *RaftNode) GetState() (term int64, isLeader bool) {
	return n.CurrentTerm, n.State == Leader
}
