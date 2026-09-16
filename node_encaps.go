package db

import (
	"errors"
	"time"
)

// Encapsulations for node roles.

func (n *Node) GetRole() Role {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.NodeRole
}

func (n *Node) AssignCandidate() {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.NodeRole = Candidate
}

func (n *Node) AssignLeader() {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.NodeRole = Leader
}

func (n *Node) AssignFollower() {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.NodeRole = Follower
}

var InvalidRoleErr = errors.New("node is not a leader.")
var InvalidSetStrTypeErr = errors.New("SetStr requires string value")
var InvalidSetIntTypeErr = errors.New("SetInt requires int value")

// Commands for which each node will replicate.
func (n *Node) SetStr(key, value string, TTL time.Duration) error {
	if n.NodeRole != Leader {
		return InvalidRoleErr
	}

	cmd := Command{
		Type:  "SET",
		Key:   key,
		Value: value,
		TTL:   TTL,
	}

	// type checking cmd.value
	value, ok := cmd.Value.(string)
	if !ok {
		return InvalidSetStrTypeErr
	}

	// applying command locally before other nodes
	err := n.DB.SetString(cmd.Type, value, cmd.TTL)
	if err != nil {
		return err
	}

	for _, r := range n.Replicas { // replicate for each node
		err := n.Replicate(r, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) SetInt(key string, value uint32, TTL time.Duration) error {
	if n.NodeRole != Leader {
		return InvalidRoleErr
	}

	cmd := Command{
		Type:  "SET",
		Key:   key,
		Value: value,
		TTL:   TTL,
	}

	val, ok := cmd.Value.(uint32)
	if !ok {
		return InvalidSetIntTypeErr
	}

	// applying command locally before other nodes
	err := n.DB.SetInt(cmd.Type, val, cmd.TTL)
	if err != nil {
		return err
	}

	for _, r := range n.Replicas { // replicate for each node
		err := n.Replicate(r, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
