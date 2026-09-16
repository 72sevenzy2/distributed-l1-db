package db

import (
	"encoding/binary"
	"errors"
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

var InvalidInputErr = errors.New("invalid input.")
var KeyAlreadyExistsErr = errors.New("key already exists.")

func (v *DB) SetInt(key string, value uint32, ttl time.Duration) error {
	if value == 0 {
		return InvalidInputErr
	}

	buf := make([]byte, 4) // uint32's has a fixed byte size of 4

	binary.LittleEndian.PutUint32(buf, value)

	v.dataLock.Lock() // excluive write lock.
	defer v.dataLock.Unlock()
	if _, ok := v.data[key]; !ok {
		v.data[key] = Entity{
			Value: "int", // will then be compared in GetAllInt() func to make sure its int
			Data:  buf,
		}
	} else {
		return KeyAlreadyExistsErr
	}

	go v.DeleteAfterTTL(ttl, key)
	return nil
}

func (v *DB) GetInt(key string) (uint32, bool) {
	v.dataLock.RLock() // read lock.
	defer v.dataLock.RUnlock()
	data, ok := v.data[key]         // data will be of type of the Entity struct
	if !ok || len(data.Data) != 4 { // check if key exists and Data has exactly 4 byte (type uint32 is fixed at 4 bytes)
		return 0, false
	}

	var value uint32         // uint32 response (will be typecated to int when displaying to user)
	if data.Value == "int" { // check if int
		value = binary.LittleEndian.Uint32(data.Data)
	}

	// value := binary.LittleEndian.Uint32(data.Data)
	return value, true
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

// string serialization and handlers

func (v *DB) SetString(key string, value string, ttl time.Duration) error {
	if value == "" {
		return InvalidInputErr
	}

	v.dataLock.Lock()
	defer v.dataLock.Unlock()
	if _, ok := v.data[key]; !ok {
		v.data[key] = Entity{
			Value: "string",
			Data:  StringToByte(value), // refers to the underlying []byte representation of value, without []byte conversions with copying
		}
	} else {
		return KeyAlreadyExistsErr
	}

	go v.DeleteAfterTTL(ttl, key)
	return nil
}

// get method for string

func (v *DB) GetString(key string) (string, bool) {
	v.dataLock.RLock()
	defer v.dataLock.RUnlock()
	val, ok := v.data[key] // val is if type Entity struct

	if !ok || val.Value != "string" { // check if string
		return "", false
	}

	// returns string representation of val.Data without a string conversion.
	return ByteToString(val.Data), true
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
