package db

import (
	"fmt"
	"testing"
	"time"
)

func TestDbInt(t *testing.T) {
	b := NewDB()
	b.SetTTL("name", 100000, time.Minute*10, Uint32Encoder{})
	b.Del("name")

	b.SetTTL("test1", 100000000, time.Minute*10, Uint32Encoder{})
	b.SetTTL("test2", 1434242324, time.Minute*10, Uint32Encoder{})

	all, _ := b.GetAll(Uint32Encoder{})

	// print all
	for k, v := range all {
		fmt.Println(k, v)
	}

	val, ok := b.Get("name", Uint32Encoder{})
	if !ok {
		fmt.Println("not found", val)
	}
	fmt.Println(val)
}
