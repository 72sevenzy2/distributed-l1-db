package main

import (
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/72sevenzy2/in-memory-database"
)

func main() {
	s := flag.String("id", "node-#", "-id <node id>")
	addr := flag.String("addr", "", "-addr <node port>")

	// start tcp server
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	replicaFlag := flag.String("replicas", "", "comma-separated replica addresses")
	flag.Parse()

	var replicas []string

	if *replicaFlag != "" {
		replicas = strings.Split(*replicaFlag, ",")
	}

	node1 := db.NewNode(*s, *addr, db.Leader, replicas)
	if err := node1.Start(); err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err.Error()) // err acceping connections
			continue
		}

		go db.HandleConnection(conn, node1)
	}

}
