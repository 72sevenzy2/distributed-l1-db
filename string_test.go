package db

import (
	"fmt"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	d := NewDB()

	err := d.SetTTL("test1", "h", time.Minute*10, StringEncoder{})
	err2 := d.SetTTL("test2", "hi2", time.Minute*10, StringEncoder{})

	if err != nil || err2 != nil {
		fmt.Println(err.Error(), err2.Error())
	}

	val, ok := d.Get("test1", StringEncoder{})
	val2, ok2 := d.Get("test2", StringEncoder{})
	if !ok || !ok2 {
		fmt.Println("data doesnt exist")
	} else {
		fmt.Println(val, val2)
	}

	res, ex := d.GetAll(StringEncoder{})
	if !ex {
		fmt.Println("no data exists")
	}

	for k, v := range res {
		fmt.Println(k, v)
	}
}
