package db

import (
	"sync"
	"time"
)

type Role uint8 // small integer which dictates a nodes role.

// Represents a nodes role, either an follower or leader, Candidate for new leader node.
const (
	Leader    Role = iota // 0
	Follower              // 1
	Candidate             // 2
)

// Command defines the details of the commands that will be processed by follower nodes after the leader node.
type Command struct {
	Type string
	// NodeID will be used for heartbeat checks.
	NodeID string

	Key   string
	Value any

	TTL time.Duration
}

// Peer represents a nodes heartbeat meta data for failure alert detection.
type Peer struct {
	LastSeen time.Time
	Alive    bool
}

// a node represents an database instance.
type Node struct {
	lock sync.RWMutex

	ID       string
	Addr     string
	NodeRole Role

	// Replicas represent a string array of existing nodes addresses.
	Replicas []string
	Peers    map[string]*Peer

	DB *DB

	// new leader node election configs.
	currentTerm uint64
	voteFor     string
}
