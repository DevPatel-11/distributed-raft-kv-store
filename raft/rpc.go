package raft

import "fmt"

// HandleRequestVote handles a RequestVote RPC from a candidate
func (n *RaftNode) HandleRequestVote(req *RequestVoteRequest) *RequestVoteResponse {
	resp := &RequestVoteResponse{
		Term:        n.CurrentTerm,
		VoteGranted: false,
	}

	// If request term is older than current term, reject
	if req.Term < n.CurrentTerm {
		return resp
	}

	// If request term is newer, update current term
	if req.Term > n.CurrentTerm {
		n.CurrentTerm = req.Term
		n.VotedFor = -1
		if n.State == Leader {
			n.State = Follower
		}
	}

	// Check if we can vote for this candidate
	lastLogIndex := int64(len(n.Log) - 1)
	var lastLogTerm int64
	if lastLogIndex >= 0 {
		lastLogTerm = n.Log[lastLogIndex].Term
	}

	if (n.VotedFor == -1 || n.VotedFor == req.CandidateID) &&
		req.LastLogTerm >= lastLogTerm &&
		req.LastLogIndex >= lastLogIndex {
		n.VotedFor = req.CandidateID
		resp.VoteGranted = true
		fmt.Printf("[Node %d] Granted vote to candidate %d for term %d\n",
			n.NodeID, req.CandidateID, req.Term)
	}

	resp.Term = n.CurrentTerm
	return resp
}

// HandleAppendEntries handles an AppendEntries RPC from a leader
func (n *RaftNode) HandleAppendEntries(req *AppendEntriesRequest) *AppendEntriesResponse {
	resp := &AppendEntriesResponse{
		Term:    n.CurrentTerm,
		Success: false,
	}

	// If request term is older than current term, reject
	if req.Term < n.CurrentTerm {
		return resp
	}

	// If request term is newer, update current term
	if req.Term > n.CurrentTerm {
		n.CurrentTerm = req.Term
		n.VotedFor = -1
	}

	// If we're not a follower, step down
	if n.State != Follower {
		n.State = Follower
	}

	// Check if previous log entry exists and has matching term
	if req.PrevLogIndex < 0 ||
		req.PrevLogIndex >= int64(len(n.Log)) ||
		(req.PrevLogIndex >= 0 && n.Log[req.PrevLogIndex].Term != req.PrevLogTerm) {
		return resp
	}

	// Append new entries
	for _, entry := range req.Entries {
		if int64(len(n.Log)) > entry.Index {
			if n.Log[entry.Index].Term != entry.Term {
				// Delete conflicting entries
				n.Log = n.Log[:entry.Index]
			}
		}
		if int64(len(n.Log)) <= entry.Index {
			n.Log = append(n.Log, entry)
		}
	}

	// Update commit index
	if req.LeaderCommit > n.CommitIndex {
		if req.LeaderCommit < int64(len(n.Log)) {
			n.CommitIndex = req.LeaderCommit
		} else {
			n.CommitIndex = int64(len(n.Log) - 1)
		}
	}

	resp.Term = n.CurrentTerm
	resp.Success = true
	return resp
}
