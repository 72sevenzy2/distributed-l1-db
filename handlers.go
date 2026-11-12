package db

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"
)

func (n *Node) HandleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	decoder.UseNumber()
	for {
		var cmd Command[any]
		if err := decoder.Decode(&cmd); err != nil {
			if !errors.Is(err, io.EOF) {
				slog.Error("ERR", "decoding_err", err.Error())
			}
			return
		}

		UCinput := strings.ToUpper(cmd.Type) // normalise to all capital

		switch UCinput {
		case "VOTE_REQUEST":
			n.HandleVoteRequest(cmd, conn)
		case "HEARTBEAT":
			n.HandleHeartbeat(cmd.NodeID, conn, cmd.Term, cmd.Role)
		case "SET":
			switch value := cmd.Value.(type) {
			case json.Number:
				parsed, err := strconv.ParseUint(value.String(), 10, 32)
				if err != nil {
					conn.Write(StringToByte("invalid integer value\n.\n"))
					continue
				}
				cmd.Value = uint32(parsed)
			case string: // preserving type strings here
			default:
				conn.Write(StringToByte("invalid value type\n.\n"))
				continue
			}

			if !cmd.Set(n, conn) {
				conn.Write(StringToByte(".\n"))
			}
		case "GET":
			ok := cmd.Get(n, conn)
			if !ok {
				conn.Write(StringToByte(".\n"))
			}
		case "FETCH":
			Fetch(n, conn)
			conn.Write(StringToByte(".\n"))
		case "DEL":
			ok := cmd.Del(n, conn)
			if !ok {
				conn.Write(StringToByte(".\n"))
			}
		case "HELP":
			conn.Write(StringToByte("General usage: \n" +
				"SET <KeyName> <value> <TTL> \n" +
				"|| GET <KeyName>\n" +
				"|| DEL <KeyName>\n" +
				"\\n\n" +
				"to exit: run <exit>\n" +
				".\n",
			))
		case "EXIT":
			conn.Write(StringToByte("exited program. \n" +
				",\n",
			))
		default:
			conn.Write(StringToByte("invalid command \n" +
				".\n",
			))
		}

	}

}

type ValueType interface {
	string | uint32
}

func (n *Node) Replicate(addr string, cmd Command[any]) error {
	conn, err := net.DialTimeout("tcp", addr, time.Second*3) // timeout safeguards against unresponsive tcp servers.
	if err != nil {
		return err
	}
	defer conn.Close()

	enc := json.NewEncoder(conn)
	if err := enc.Encode(cmd); err != nil {
		return err
	}
	return nil
}
