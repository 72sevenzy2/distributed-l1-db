package db

import "encoding/binary"

// generic encoder/decoder.
type Encoder[T any] interface {
	TypeID() string
	Encode(T) ([]byte, error)
	Decode([]byte) (T, error)
}

type String struct{}
type Integer struct{}

// TypeID() returns the underlying generic type of an encoder/decoder.
func (String) TypeID() string {
	return "string"
}

func (Integer) TypeID() string {
	return "uint32"
}

func (String) Encode(v string) ([]byte, error) {
	return StringToByte(v), nil
}

func (String) Decode(v []byte) (string, error) {
	return ByteToString(v), nil
}

func (Integer) Encode(v uint32) ([]byte, error) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, v)
	return buf, nil
}

func (Integer) Decode(v []byte) (uint32, error) {
	if len(v) != 4 {
		return 0, InvalidInputErr
	}

	return binary.LittleEndian.Uint32(v), nil
}
