package db

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"time"
)

func NewNode(id, addr string, role Role, replicas []string) *Node {
	return &Node{
		ID:       id,
		Addr:     addr,
		NodeRole: role,
		Replicas: replicas,
		Peers:    make(map[string]*Peer),
		DB:       NewDB(),
	}
}

func (n *Node) Start() error {
	l, err := net.Listen("tcp", n.Addr)
	if err != nil {
		return err
	}

	defer l.Close()

	go n.Heartbeatloop()
	go n.FailureDetectionloop()

	fmt.Println("node listening on addr:", n.Addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go n.HandleConnection(conn)
	}
}

// node-to-node heartbeat logic.
func (n *Node) Heartbeatloop() {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for range t.C {
		n.lock.RLock()
		role := n.NodeRole
		n.lock.RUnlock()

		if role == Leader {
			go n.Sendheartbeat()
		}
	}
}

// detects failed heartbeats from nodes.
func (n *Node) FailureDetectionloop() {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for range t.C {
		n.lock.Lock()
		defer n.lock.Unlock()

		//heartbeart timeout until its considered dead.
		timeout := 3 * time.Second
		now := time.Now()

		for id, peer := range n.Peers {
			if now.Sub(peer.LastSeen) > timeout {
				if peer.Alive {
					slog.Warn("PEER_DEAD", "node", id)
				}
				peer.Alive = false
			}
		}
	}
}

func (n *Node) Sendheartbeat() {
	for _, replica := range n.Replicas {
		conn, err := net.DialTimeout("tcp", replica, time.Millisecond*500)
		if err != nil {
			slog.Error("ERR", "heartbeat_err", err)
			continue
		}

		msg := &Command{
			Type:   "heartbeat",
			NodeID: n.ID,
		}

		enc := json.NewEncoder(conn)
		if err := enc.Encode(msg); err != nil {
			slog.Error("ERR", "encoding_err", err)
			conn.Close()
			continue
		}

		var ack Command
		dec := json.NewDecoder(conn)
		if err := dec.Decode(&ack); err != nil {
			slog.Error("ERR", "decoding_err", err)
			conn.Close()
			return
		}
		conn.Close()

		if ack.Type != "heartbeat_ack" {
			slog.Error("ERR", "invalid_heartbeat_type", ack.Type)
			return
		}
		slog.Info("HEARTBEAT_ACK", "received_ack", ack.NodeID)
	}
}

func (n *Node) HandleHeartbeat(id string, conn net.Conn) {
	n.lock.Lock()
	n.Peers[id] = &Peer{
		LastSeen: time.Now(),
		Alive:    true,
	}
	n.lock.Unlock()

	msg := &Command{
		Type:   "heartbeat_ack",
		NodeID: n.ID,
	}

	enc := json.NewEncoder(conn)
	if err := enc.Encode(msg); err != nil {
		slog.Error("ERR", "encoding_err", err)
		return
	}
}
