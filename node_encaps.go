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

// non-TTL based abstractions over db methods.

func (n *Node) SetStr(key, value string) error {
	if n.GetRole() != Leader {
		return InvalidRoleErr
	}

	cmd := Command{
		Type:  "SET",
		Key:   key,
		Value: value,
	}
	// // // type checking cmd.value
	value, ok := cmd.Value.(string)
	if !ok {
		return InvalidSetStrTypeErr
	}

	// applying command locally before other nodes
	err := n.DB.Set(key, value, String{})
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

func (n *Node) SetInt(key string, value uint32) error {
	if n.GetRole() != Leader {
		return InvalidRoleErr
	}

	cmd := Command{
		Type:  "SET",
		Key:   key,
		Value: value,
	}

	val, ok := cmd.Value.(uint32)
	if !ok {
		return InvalidSetIntTypeErr
	}

	// applying command locally before other nodes
	err := n.DB.Set(key, val, Integer{})
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

// TTL-based abstractions over db methods.
func (n *Node) SetStr_TTL(key, value string, TTL time.Duration) error {
	if n.GetRole() != Leader {
		return InvalidRoleErr
	}

	cmd := Command{
		Type:  "SET",
		Key:   key,
		Value: value,
		TTL:   TTL,
	}
	//
	// // // type checking cmd.value
	value, ok := cmd.Value.(string)
	if !ok {
		return InvalidSetStrTypeErr
	}

	// applying command locally before other nodes
	err := n.DB.SetTTL(key, value, cmd.TTL, String{})
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

func (n *Node) SetInt_TTL(key string, value uint32, TTL time.Duration) error {
	if n.GetRole() != Leader {
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
	err := n.DB.SetTTL(key, val, cmd.TTL, Integer{})
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
