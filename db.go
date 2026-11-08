package db

import (
	"encoding/binary"
	"sync"
	"time"

	"unsafe"
)

// metadata for value types
type Entity struct {
	Data  []byte
	Value string
}

// core logic
type DB struct {
	dataLock sync.RWMutex // allows for high concurrent reads without potential for-write lock contention.
	data     map[string]Entity
}

func NewDB() *DB { // initialise a new map to hold data
	return &DB{
		data: make(map[string]Entity),
	}
}

func StringToByte(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func ByteToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// ttl deletion worker
func (v *DB) DeleteAfterTTL(ttl time.Duration, key string) {
	time.Sleep(ttl)
	v.Del(key)
}

func (v *DB) Set[T any](key string, value T, ttl time.Duration, encoder Encoder[T]) error {
	data, err := encoder.Encode(value)
	if err != nil {
		return err
	}

	v.dataLock.Lock() // excluive write lock.
	defer v.dataLock.Unlock()
	if _, ok := v.data[key]; !ok {
		v.data[key] = Entity{
			Value: "int", // will then be compared in GetAllInt() func to make sure its int
			Data:  data,
		}
	} else {
		return KeyAlreadyExistsErr
	}

	go v.DeleteAfterTTL(ttl, key)
	return nil
}

func (v *DB) Get[T any](key string, decoder Encoder[T]) (T, bool) {
	v.dataLock.RLock()      // read lock.
	data, ok := v.data[key] // data will be of type of the Entity struct
	v.dataLock.RUnlock()

	if !ok { // check if key exists and Data has exactly 4 byte (type uint32 is fixed at 4 bytes)
		var zero T
		return zero, false
	}

	val, err := decoder.Decode(data.Data)
	if err != nil {
		var zero T
		return zero, false
	}

	// value := binary.LittleEndian.Uint32(data.Data)
	return val, true
}

// func to retrieve all values at once
func (v *DB) GetAllInt() (map[string]uint32, bool) {
	result := make(map[string]uint32)

	// futur: keep track of all data types stored in db and print as so without looping through each.
	v.dataLock.RLock()
	defer v.dataLock.RUnlock()
	for k, v := range v.data {
		if v.Value == "int" { // make sure type is int before serialization
			result[k] = binary.LittleEndian.Uint32(v.Data)
		}
	}

	if len(result) == 0 { // no existing data case
		return nil, false
	}

	return result, true
}

// display all string value data from db

func (v *DB) GetAllString() (map[string]string, bool) {
	results := make(map[string]string)

	v.dataLock.RLock()
	defer v.dataLock.RUnlock()
	for k, val := range v.data {
		if val.Value == "string" {
			results[k] = ByteToString(val.Data)
		}
	}

	if len(results) == 0 { // handling no existing data case
		return nil, false
	}

	return results, true
}

// this will stay fixed as key will always be type string
func (v *DB) Del(key string) {
	v.dataLock.Lock()
	defer v.dataLock.Unlock()
	delete(v.data, key) // built in delete() func to delete a particular key being held in the db map.
}
