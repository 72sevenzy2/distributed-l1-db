package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	db "github.com/72sevenzy2/in-memory-database"
)

func main() {
	// pre-requisites: must start a tcp server before-hand on port 9000.
	conn, err := net.Dial("tcp", ":9000") // connect with tcp server
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	reader := bufio.NewReader(conn)

	for {
		fmt.Print("> ")
		safe := scanner.Scan()
		if !safe {
			err := scanner.Err()
			if err != nil {
				fmt.Println(err)
			}
			break
		}

		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}

		cmd := db.Command[any]{
			Type: strings.ToUpper(parts[0]),
		}

		switch cmd.Type {
		case "SET":
			if len(parts) < 3 || len(parts) > 4 {
				fmt.Println("usage: SET <key> <value> [TTL, e.g. 30s/2m]")
				continue
			}
			cmd.Key = parts[1]

			cmd.Value = parts[2]
			val, err := strconv.ParseUint(parts[2], 10, 32)
			if err == nil {
				cmd.Value = uint32(val)
			}

			if len(parts) == 4 {
				ttl, err := time.ParseDuration(parts[3])
				if err != nil {
					fmt.Println("invalid ttl:", err)
					continue
				}
				cmd.TTL = ttl
			}

		case "DEL":
			if len(parts) != 2 {
				fmt.Println("usage: DEL <key>")
				continue
			}

			cmd.Key = parts[1]

		case "GET":
			if len(parts) != 2 {
				fmt.Println("usage: GET <key>")
				continue
			}
			cmd.Key = parts[1]

		case "EDIT":
			if len(parts) != 3 {
				fmt.Println("usage: EDIT <key> <new value>")
				continue
			}

			cmd.Key = parts[1]
			cmd.Value = parts[2]

			val, err := strconv.ParseUint(parts[2], 10, 32)
			if err == nil {
				cmd.Value = uint32(val)
			}
		}
		if err := json.NewEncoder(conn).Encode(cmd); err != nil {
			fmt.Println(err)
			break
		}

		for {
			resp, err2 := reader.ReadString('\n')
			if resp == ".\n" {
				break
			}
			// exit delimeter
			if resp == ",\n" {
				return
			}
			if err2 != nil {
				if errors.Is(err2, io.EOF) {
					return
				}
				fmt.Println(err2.Error())
				continue
			}

			fmt.Print(resp)
		}
	}

}
