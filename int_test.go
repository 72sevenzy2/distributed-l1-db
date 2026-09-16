package db

import (
	"fmt"
	"testing"
	"time"
)

func TestDbInt(t *testing.T) {
	b := NewDB()
	b.SetInt("name", 100000, time.Minute*10)
	b.Del("name")

	b.SetInt("test1", 100000000, time.Minute*10)
	b.SetInt("test2", 1434242324, time.Minute*10)

	all, _ := b.GetAllInt()

	// print all
	for k, v := range all {
		fmt.Println(k, v)
	}

	val, ok := b.GetInt("name")
	if !ok {
		fmt.Println("not found", val)
	}
	fmt.Println(val)
}
