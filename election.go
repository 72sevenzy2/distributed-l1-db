package db

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"time"
)

func (n *Node) ElectionLoop() {
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	for range t.C {
		n.lock.RLock()

		if n.GetRole() == Leader {
			n.lock.RUnlock()
			continue
		}

		now := time.Now()
		leaderDied := true
		for _, peer := range n.Peers {
			if peer.Role == Leader && peer.Alive && now.Sub(peer.LastSeen) < 3*time.Second {
				leaderDied = false // meaning that the leader had not died yet.
				break
			}
		}
		n.lock.RUnlock()

		if leaderDied {
			go n.StartElection()
		}
	}
}

func (n *Node) StartElection() {
	n.lock.Lock()

	// skip if leader
	if n.NodeRole == Leader {
		n.lock.Unlock()
		return
	}

	// update currentTerm
	n.currentTerm++

	n.NodeRole = Candidate

	n.votedFor = n.ID

	term := n.currentTerm
	n.lock.Unlock()

	slog.Info("ELECTION_BEGUN", "node", n.ID, "term", term)
	votes := 1

	for _, replica := range n.Replicas {
		granted, err := n.RequestVote(replica, term)

		if err != nil {
			slog.Error("ERR", "vote_request_err", err, "node", replica)
			continue
		}

		if granted {
			votes++
		}
	}

	cluster := len(n.Replicas) + 1
	// quorum represents the number of votes required for a node to be elected.
	quorum := cluster/2 + 1

	if n.currentTerm != term {
		return
	}

	if n.NodeRole != Candidate {
		return
	}

	if votes >= quorum {
		n.NodeRole = Leader

		slog.Info("LEADER_ELECTED", "node", n.ID, "term", n.currentTerm, "votes", votes)
	}
}

func (n *Node) RequestVote(addr string, term uint64) (bool, error) {
	conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)

	if err != nil {
		return false, err
	}

	defer conn.Close()

	// ask another peer for its vote.
	msg := &Command{
		Type:   "vote_request",
		NodeID: n.ID,
		Term:   term,
	}

	enc := json.NewEncoder(conn)
	if err := enc.Encode(msg); err != nil {
		return false, err
	}

	var resp Command

	dec := json.NewDecoder(conn)
	if err := dec.Decode(&resp); err != nil {
		return false, err
	}

	// peer has a newer recent term
	if resp.Term > term {
		n.lock.Lock()

		if resp.Term > n.currentTerm {
			// if peers term is greater than current nodes term, then the peer is the one starting the Election, thus, this node should be a Follower node.
			n.currentTerm = resp.Term
			n.NodeRole = Follower
			n.votedFor = ""
		}

		n.lock.Unlock()
		return false, nil
	}

	if resp.Type != "vote_response" {
		return false, errors.New("invalid response type")
	}

	return resp.GivenVote, nil
}

// HandleVoteRequest validates whether nodes term is either outdated/up to date, which then allows it to vote.
func (n *Node) HandleVoteRequest(msg Command, conn net.Conn) {
	n.lock.Lock()
	defer n.lock.Unlock()

	allowed := false

	if msg.Term < n.currentTerm {
		// outdated
		allowed = false
	} else {
		if msg.Term > n.currentTerm { // up to date candidate node term.
			n.currentTerm = msg.Term
			n.NodeRole = Follower // assign current node to Follower.
			n.votedFor = ""
		}

		// n.votedFor as fast path.
		if n.votedFor == "" || n.votedFor == msg.NodeID { // or if current node voted for candidate.
			n.votedFor = msg.NodeID
			allowed = true
		}
	}

	resp := Command{
		Type:      "vote_request",
		NodeID:    n.ID,
		Term:      n.currentTerm,
		GivenVote: allowed,
	}

	enc := json.NewEncoder(conn)
	if err := enc.Encode(resp); err != nil {
		slog.Error("ERR", "encoding_err", err)
	}
}
