package db

import (
	"encoding/json"
	"log/slog"
	"net"
	"strings"
	"time"
)

func (n *Node) HandleConnection(conn net.Conn) {
	defer conn.Close()

	var cmd Command

	d := json.NewDecoder(conn)
	if err := d.Decode(&cmd); err != nil {
		slog.Error("ERR", "decoding_err", err.Error())
		return
	}

	valType, ok := cmd.GetValueType()
	if !ok {
		slog.Error("ERR", "type_error", "unknown value type.")
		return
	}

	for {
		UCinput := strings.ToUpper(cmd.Type) // normalise to all capital

		switch UCinput {
		case "VOTE_REQUEST":
			n.HandleVoteRequest(cmd, conn)
		case "HEARTBEAT":
			n.HandleHeartbeat(cmd.NodeID, conn, cmd.Term, cmd.Role)
		case "SET":
			switch valType {
			case "uint32":
				ok := cmd.Set[uint32](n, conn)
				if !ok {
					conn.Write(StringToByte(".\n"))
				}
			case "string":
				ok := cmd.Set[string](n, conn)
				if !ok {
					conn.Write(StringToByte(".\n"))
				}
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

func (n *Node) Replicate(addr string, cmd Command) error {
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
