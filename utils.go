package db

import (
	"fmt"
	"net"
	"strconv"

	"unsafe"
)

func StringToByte(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func ByteToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// utility functions for server/main.go

func (c *Command) Set(node *Node, conn net.Conn) bool {
	//parsing values type.
	val, err := c.Value.(uint32)

	if err {
		if c.TTL != 0 {
			err2 := node.SetInt_TTL(c.Key, val, c.TTL)
			if err2 != nil {
				conn.Write(StringToByte(err2.Error() + "\n"))
				return false
			}
			conn.Write(StringToByte("successful\n" +
				".\n",
			))
			return true
		}

		err2 := node.SetInt(c.Key, val)
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
		if c.TTL != 0 {
			err2 := node.SetStr_TTL(c.Key, val2, c.TTL)
			if err2 != nil {
				conn.Write(StringToByte(err2.Error() + "\n"))
				return false
			}
			conn.Write(StringToByte("successful\n" +
				".\n",
			))
			return true
		}
		err2 := node.SetStr(c.Key, val2)
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
	val, ok := n.DB.Get(c.Key, Uint32Encoder{})
	if !ok {
		val2, ok2 := n.DB.Get(c.Key, StringEncoder{})
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
	vals, ok := n.DB.GetAll(Uint32Encoder{})   // returns map[string]uint32, bool
	vals2, ok2 := n.DB.GetAll(StringEncoder{}) // returns map[string]string, bool

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
	if _, ok := n.DB.Get(c.Key, Uint32Encoder{}); ok {
		n.DB.Del(c.Key)
		conn.Write(StringToByte("successfully deleted key\n" +
			".\n",
		))
		return true
	}
	if _, ok := n.DB.Get(c.Key, StringEncoder{}); ok {
		n.DB.Del(c.Key)
		conn.Write(StringToByte("successfuly deleted key\n" +
			".\n",
		))
		return true
	}
	return false
}

func (c Command) GetValueType() (string, bool) {
	switch c.Value.(type) {
	case string:
		return "string", true

	case uint32:
		return "uint32", true
	default:
		return "", false
	}
}
