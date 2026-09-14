package db

import (
	"sync"
	"time"
)

type Role uint8 // small integer which dictates a nodes role.

// Represents a nodes role, either an follower or leader.
const (
	Leader Role = iota
	Follower
)

// Command defines the details of the commands that will be processed by follower nodes after the leader node.
type Command struct {
	Type  string
	Key   string
	Value any

	TTL time.Duration
}

// a node represents an database instance.
type Node struct {
	lock sync.RWMutex

	ID       string
	Addr     string
	NodeRole Role

	// Replicas represent a string array of existing nodes addresses.
	Replicas  []string
	HeartBeat map[string]time.Time

	DB *DB
}

type HeartBeatPayload struct {
	Type   string
	NodeID string
}
