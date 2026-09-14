package main

import (
	"flag"
	"strings"

	"github.com/72sevenzy2/in-memory-database"
)

func main() {
	s := flag.String("id", "node-#", "-id <node id>")
	addr := flag.String("addr", "", "-addr <node port>")

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
}
