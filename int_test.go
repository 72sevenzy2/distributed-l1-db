package db

import (
	"fmt"
	"testing"
	"time"
)

func TestDbInt(t *testing.T) {
	b := NewDB()
	b.SetTTL("name", 100000, time.Minute*10, Integer{})
	b.Del("name")

	b.SetTTL("test1", 100000000, time.Minute*10, Integer{})
	b.SetTTL("test2", 1434242324, time.Minute*10, Integer{})

	all, _ := b.GetAll(Integer{})

	// print all
	for k, v := range all {
		fmt.Println(k, v)
	}

	val, ok := b.Get("name", Integer{})
	if !ok {
		fmt.Println("not found", val)
	}
	fmt.Println(val)
}
