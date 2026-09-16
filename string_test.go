package db

import (
	"fmt"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	d := NewDB()

	err := d.SetString("test1", "h", time.Minute*10)
	err2 := d.SetString("test2", "hi2", time.Minute*10)

	if err != nil || err2 != nil {
		fmt.Println(err.Error(), err2.Error())
	}

	val, ok := d.GetString("test1")
	val2, ok2 := d.GetString("test2")
	if !ok || !ok2 {
		fmt.Println("data doesnt exist")
	} else {
		fmt.Println(val, val2)
	}

	res, ex := d.GetAllString()
	if !ex {
		fmt.Println("no data exists")
	}

	for k, v := range res {
		fmt.Println(k, v)
	}
}
