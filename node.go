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
		ID:        id,
		Addr:      addr,
		NodeRole:  role,
		Replicas:  replicas,
		HeartBeat: make(map[string]time.Time),
		DB:        NewDB(),
	}
}

func (n *Node) Start() error {
	l, err := net.Listen("tcp", n.Addr)
	if err != nil {
		return err
	}

	defer l.Close()

	go n.Heartbeatloop()

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
	for range t.C {
		go n.Sendheartbeat()
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
		conn.Close()
		if err := enc.Encode(msg); err != nil {
			slog.Error("ERR", "encoding_err", err)
			continue
		}
	}
}

func (n *Node) HandleHeartbeat(id string, t time.Time) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.HeartBeat[id] = t
}
