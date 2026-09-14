package db

import (
	"fmt"
	"net"
)

func NewNode(id, addr string, role Role, replicas []string) *Node {
	return &Node{
		ID:       id,
		Addr:     addr,
		NodeRole: role,
		Replicas: replicas,
		DB:       NewDB(),
	}
}

func (n *Node) Start() error {
	l, err := net.Listen("tcp", n.Addr)
	if err != nil {
		return err
	}

	defer l.Close()

	fmt.Println("node listening on addr:", n.Addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go n.HandleConnection(conn)
	}
}


