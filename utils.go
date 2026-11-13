package db

import (
	"encoding/json"
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

func (c *Command[V]) Set(node *Node, conn net.Conn) bool {
	//parsing values type.
	val, err := any(c.Value).(uint32)

	if err {
		if c.TTL != 0 {
			err2 := node.Set_TTL(c.Key, val, c.TTL)
			if err2 != nil {
				conn.Write(StringToByte(err2.Error() + "\n"))
				return false
			}
			conn.Write(StringToByte("successful\n" +
				".\n",
			))
			return true
		}

		err2 := node.Set(c.Key, val)
		if err2 != nil {
			conn.Write(StringToByte(err2.Error() + "\n"))
			return false
		}
		conn.Write(StringToByte("successful\n" +
			".\n",
		))
		return true
	}

	val2, err := any(c.Value).(string)
	if err {
		if c.TTL != 0 {
			err2 := node.Set_TTL(c.Key, val2, c.TTL)
			if err2 != nil {
				conn.Write(StringToByte(err2.Error() + "\n"))
				return false
			}
			conn.Write(StringToByte("successful\n" +
				".\n",
			))
			return true
		}
		err2 := node.Set(c.Key, val2)
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
func (c *Command[V]) Get(n *Node, conn net.Conn) bool {
	val, ok := n.DB.Get(c.Key, Integer{})
	if !ok {
		val2, ok2 := n.DB.Get(c.Key, String{})
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
	vals, ok := n.DB.GetAll(Integer{})  // returns map[string]uint32, bool
	vals2, ok2 := n.DB.GetAll(String{}) // returns map[string]string, bool

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

func (c *Command[V]) Edit(n *Node, conn net.Conn) bool {
	var ok bool
	switch value := any(c.Value).(type) {
	case string:
		ok = n.DB.Edit(c.Key, value, String{})
	case json.Number:
		parsed, err := strconv.ParseUint(value.String(), 10, 32)
		if err != nil {
			conn.Write(StringToByte("invalid integer value\n.\n"))
			return false
		}

		ok = n.DB.Edit(c.Key, uint32(parsed), Integer{})
	default:
		conn.Write(StringToByte("unknown value type." + "\n"))
		return false
	}

	if !ok {
		conn.Write(StringToByte("data does not exist\n"))
		return false
	}

	conn.Write(StringToByte("successful\n" +
		".\n",
	))
	return true
}

func (c *Command[V]) Del(n *Node, conn net.Conn) bool {
	if _, ok := n.DB.Get(c.Key, Integer{}); ok {
		n.DB.Del(c.Key)
		conn.Write(StringToByte("successfully deleted key\n" +
			".\n",
		))
		return true
	}
	if _, ok := n.DB.Get(c.Key, String{}); ok {
		n.DB.Del(c.Key)
		conn.Write(StringToByte("successfuly deleted key\n" +
			".\n",
		))
		return true
	}
	return false
}

func (c Command[V]) GetValueType() (string, bool) {
	switch any(c.Value).(type) {
	case string:
		return "string", true

	case uint32:
		return "uint32", true
	default:
		return "", false
	}
}
