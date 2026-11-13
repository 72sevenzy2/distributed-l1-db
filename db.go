package db

import (
	"sync/atomic"
	"time"
)

// metadata for value types
type Entity struct {
	Data  []byte
	Value string
}

type l1_database map[string]Entity

type DB struct {
	// l1 map with mutex-free reads/writes.
	data atomic.Pointer[l1_database]
}

func NewDB() *DB {
	db := make(l1_database)
	db_instance := &DB{}
	db_instance.data.Store(&db)

	return db_instance
}

func (v *DB) DeleteAfterTTL(ttl time.Duration, key string) {
	time.Sleep(ttl)
	v.Del(key)
}

func (v *DB) SetTTL[T any](key string, value T, ttl time.Duration, encoder Encoder[T]) error {
	data, err := encoder.Encode(value)
	if err != nil {
		return err
	}

	for {
		entries := v.data.Load()
		_, ok := (*entries)[key]
		if !ok {
			// copying on write for data race prevention.
			next := make(l1_database, len(*entries)+1)
			for k, val := range *entries {
				next[k] = val
			}
			newEntity := Entity{
				Value: encoder.TypeID(),
				Data:  data,
			}
			next[key] = newEntity

			// check the map we loaded against the map we're appending.
			if v.data.CompareAndSwap(entries, &next) {
				go v.DeleteAfterTTL(ttl, key)
				return nil
			}
		} else {
			return KeyAlreadyExistsErr
		}
	}
}

// Edit mutates an existing entry in the database.
func (v *DB) Edit[T any](key string, value T, encoder Encoder[T]) bool {
	data, err := encoder.Encode(value)
	if err != nil {
		return false
	}

	for {
		entries := v.data.Load()
		_, ok := (*entries)[key]
		if ok {
			// copying on write for data race prevention.
			newDB := make(l1_database, len(*entries)+1)
			for k, val := range *entries {
				newDB[k] = val
			}
			newEntity := Entity{
				Value: encoder.TypeID(),
				Data:  data,
			}
			newDB[key] = newEntity

			if v.data.CompareAndSwap(entries, &newDB) {
				return true
			}
		} else {
			return false
		}
	}
}

func (v *DB) Set[T any](key string, value T, encoder Encoder[T]) error {
	data, err := encoder.Encode(value)
	if err != nil {
		return err
	}

	for {
		entries := v.data.Load()
		_, ok := (*entries)[key]
		if !ok {
			// copying on write for data race prevention.
			next := make(l1_database, len(*entries)+1)
			for k, val := range *entries {
				next[k] = val
			}
			newEntity := Entity{
				Value: encoder.TypeID(),
				Data:  data,
			}
			next[key] = newEntity

			if v.data.CompareAndSwap(entries, &next) {
				return nil
			}
		} else {
			return KeyAlreadyExistsErr
		}
	}
}

func (v *DB) Get[T any](key string, decoder Encoder[T]) (T, bool) {
	entries := v.data.Load()
	data, ok := (*entries)[key]
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

	entries := v.data.Load()
	for k, entity := range *entries {
		if entity.Value != decoder.TypeID() {
			continue
		}
		data, err := decoder.Decode(entity.Data)
		if err != nil {
			continue
		}
		result[k] = data
	}

	if len(result) == 0 { // no existing entries
		return nil, false
	}

	return result, true
}

func (v *DB) Del(key string) {
	delete(*v.data.Load(), key)
}
