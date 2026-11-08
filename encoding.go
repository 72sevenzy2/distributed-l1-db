package db

import "encoding/binary"

//
type Encoder[T any] interface {
	Encode(T) ([]byte, error)
	Decode([]byte) (T, error)
}

type StringEncoder struct{}
type Uint32Encoder struct{}

func (StringEncoder) Encode(v string) ([]byte, error) {
	return StringToByte(v), nil
}

func (StringEncoder) Decode(v []byte) (string, error) {
	return ByteToString(v), nil
}

func (Uint32Encoder) Encode(v uint32) ([]byte, error) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, v)
	return buf, nil
}

func (Uint32Encoder) Decode(v []byte) (uint32, error) {
	if len(v) != 4 {
		return 0, InvalidInputErr
	}

	return binary.LittleEndian.Uint32(v), nil
}
