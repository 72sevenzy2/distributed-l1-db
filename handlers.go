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

	_, ok := cmd.Value.(string)
	_, ok2 := cmd.Value.(uint32)
	if !ok && !ok2 {
		slog.Error("ERR", "type_error", "unknown value type.")
		return
	}

	for {
		UCinput := strings.ToUpper(cmd.Type) // normalise to all capital

		switch UCinput {
		case "HEARTBEAT":
			n.HandleHeartbeat(cmd.NodeID, conn, cmd.Term)
		case "SET":
			ok := cmd.Set(n, conn)
			if !ok {
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
