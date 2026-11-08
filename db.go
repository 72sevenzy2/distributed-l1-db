package db

import (
	"sync"
	"time"
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

func (v *DB) DeleteAfterTTL(ttl time.Duration, key string) {
	time.Sleep(ttl)
	v.Del(key)
}

func (v *DB) Set[T any](key string, value T, ttl time.Duration, encoder Encoder[T]) error {
	data, err := encoder.Encode(value)
	if err != nil {
		return err
	}

	v.dataLock.Lock()
	defer v.dataLock.Unlock()
	if _, ok := v.data[key]; !ok {
		v.data[key] = Entity{
			Value: encoder.TypeID(),
			Data:  data,
		}
	} else {
		return KeyAlreadyExistsErr
	}

	go v.DeleteAfterTTL(ttl, key)
	return nil
}

func (v *DB) Get[T any](key string, decoder Encoder[T]) (T, bool) {
	v.dataLock.RLock()
	data, ok := v.data[key]
	v.dataLock.RUnlock()

	if !ok {
		var zero T
		return zero, false
	}

	val, err := decoder.Decode(data.Data)
	if err != nil {
		var zero T
		return zero, false
	}

	return val, true
}

// GetAll retrieves all the data in the db.
func (v *DB) GetAll[T any](decoder Encoder[T]) (map[string]T, bool) {
	result := make(map[string]T)

	v.dataLock.RLock()
	defer v.dataLock.RUnlock()
	for k, entity := range v.data {
		if entity.Value != decoder.TypeID() {
			continue
		}
		data, err := decoder.Decode(entity.Data)
		if err != nil {
			continue
		}
		result[k] = data
	}

	if len(result) == 0 { // no existing data case
		return nil, false
	}

	return result, true
}

func (v *DB) Del(key string) {
	v.dataLock.Lock()
	defer v.dataLock.Unlock()
	delete(v.data, key) // built in delete() func to delete a particular key being held in the db map.
}
