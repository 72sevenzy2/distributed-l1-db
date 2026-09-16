package db

import (
	"fmt"
	"net"
	"strconv"
)

// utility functions for server/main.go

func (c *Command) Set(node *Node, conn net.Conn) bool {
	//if len(parts) < 4 || len(parts) > 4 {
	//conn.Write(StringToByte("invalid SET format:\n" +
	//	"SET <KeyName> <value> <TTL Expiration (in minutes, eg: 5)>\n",
	//	))
	//	return false
	//	}

	//n, err := strconv.Atoi(parts[3]) // parse ttl expiration
	//if err != nil {
	//	conn.Write(StringToByte("please include a valid TTL, (eg: 10, will be in minutes)\n"))
	//return false
	//	}

	//f, err := strconv.ParseUint(parts[2], 10, 32) // returns uint64, err.

	//parsing values type.
	val, err := c.Value.(uint32)

	if err { // its a int.
		// prevent f from overflowing if number entered is too big
		//if val > math.MaxUint32 {
		//	conn.Write(StringToByte("please include a number value within range of unsigned int32.\n"))
		///	return false
		//	}

		err2 := node.SetInt(c.Key, val, c.TTL)
		if err2 != nil {
			conn.Write(StringToByte(err2.Error() + "\n"))
			return false
		}
		conn.Write(StringToByte("successful\n" +
			".\n",
		))
		return true
	}

	val2, err := c.Value.(string)
	if err {
		// its a string if unable to parse to uint.
		err2 := node.SetStr(c.Key, val2, c.TTL)
		if err2 != nil {
			conn.Write(StringToByte(err2.Error() + "\n"))
			return false
		}
		conn.Write(StringToByte("successful\n" +
			".\n",
		))
		return true
	}
	conn.Write(StringToByte("unknown value type." + "\n"))
	return false
}

func (c *Command) Get(n *Node, conn net.Conn) bool {
	//if len(parts) < 2 || len(parts) > 2 {
	//conn.Write(StringToByte("invalid GET format:\n" +
	//	"GET <KeyName>\n",
	//	))
	//	return false
	//	}

	val, ok := n.DB.GetInt(c.Key)
	if !ok {
		val2, ok2 := n.DB.GetString(c.Key)
		if !ok2 {
			conn.Write(StringToByte("data does not exist\n"))
			return false
		}
		conn.Write(StringToByte(val2 + "\n" +
			".\n",
		))
		return true
	}
	conn.Write(StringToByte(strconv.FormatUint(uint64(val), 10) + "\n" +
		".\n",
	)) // convert uint32 to readable format
	return true
}

func Fetch(n *Node, conn net.Conn) {
	vals, ok := n.DB.GetAllInt()      // returns map[string]uint32, bool
	vals2, ok2 := n.DB.GetAllString() // returns map[string]string, bool

	if !ok && !ok2 {
		conn.Write(StringToByte("no available key-values" + ".\n"))
		return
	}

	for k, v := range vals {
		fmt.Fprint(conn, k+" ", int(v), ".\n")
	}
	for k, v := range vals2 {
		fmt.Fprint(conn, k+" ", v, ".\n")
	}
}

func (c *Command) Del(n *Node, conn net.Conn) bool {
	//if len(parts) < 2 || len(parts) > 2 {
	//	conn.Write(StringToByte("usage: DEL <KeyName>\n"))
	//return false
	//	}
	if _, ok := n.DB.GetInt(c.Key); ok {
		n.DB.Del(c.Key)
		conn.Write(StringToByte("successfully deleted key\n" +
			".\n",
		))
		return true
	}
	if _, ok := n.DB.GetString(c.Key); ok {
		n.DB.Del(c.Key)
		conn.Write(StringToByte("successfuly deleted key\n" +
			".\n",
		))
		return true
	}
	return false
}
