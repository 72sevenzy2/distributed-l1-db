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
var InvalidTypeErr = errors.New("invalid value type.")

// non-TTL based abstractions over db methods.

func (n *Node) Set[T any](key string, value T) error {
	if n.GetRole() != Leader {
		return InvalidRoleErr
	}

	cmd := Command[any]{
		Type:  "SET",
		Key:   key,
		Value: value,
	}

	val, ok := any(value).(string)
	if ok {
		// applying command locally before other nodes
		err := n.DB.Set(key, val, String{})
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
	val2, ok2 := any(value).(uint32)
	if ok2 {
		// applying command locally before other nodes
		err := n.DB.Set(key, val2, Integer{})
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

	return InvalidTypeErr
}

// TTL-based abstractions over db methods.
func (n *Node) Set_TTL[T any](key string, value T, TTL time.Duration) error {
	if n.GetRole() != Leader {
		return InvalidRoleErr
	}

	cmd := Command[any]{
		Type:  "SET",
		Key:   key,
		TTL:   TTL,
		Value: value,
	}

	val, ok := any(value).(string)
	if ok {
		// applying command locally before other nodes
		err := n.DB.SetTTL(key, val, cmd.TTL, String{})
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
	val2, ok2 := any(value).(uint32)
	if ok2 {
		// applying command locally before other nodes
		err := n.DB.SetTTL(key, val2, cmd.TTL, Integer{})
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

	return InvalidTypeErr
}
