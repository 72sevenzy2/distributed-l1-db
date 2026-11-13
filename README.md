<h1 align="center"> k-v distributed l1 db. </h1>
<br>
<ul>
  <li>mutex-free concurrency-safe reads with atomic.Pointer[T].</li>
  <li>optional TTL-based eviction for short-lived data.</li>
  <li>utilises an inter-node based communication system with a leader election mechanism upon leader failure.</li>
  <li>centralised serverside error handling.</li>
  <li>interactive cli mode, which allows storing data as well as data retrieval.</li>
  <li>encodes integer type values to bytes before storing to the in-memory db for optimised performance during node-to-node communication.</li>
</ul>

<br>
<h1 align="center">pre-requisite:</h1>
<h2 align="center">go version 1.27+</h2>
<h1 align="center">instillation</h1>

```
go get github.com/72sevenzy2/in-memory-database
```

<br>
<h1 align="center">getting started:</h1>

```
package main

import (
	"fmt"
	db "github.com/72sevenzy2/in-memory-database"
)

func main() {
	v := db.NewDB()
	// db.Integer{} specifies that its an int-type value.
	err := v.Set("test1", 200, db. db.Integer{}) // setting ints
	if err != nil {
		panic(err)
	}

	// db.String{} specifies that its an string-type value.
	err2 := v.Set("test2", "hello", db.String{}) // setting strings
	if err2 != nil {
		panic(err2)
	}
}
```
<h3 align="center">setting data with TTLs:</h3>

```
package main

import (
	"fmt"
	db "github.com/72sevenzy2/in-memory-database"
)

func main() {
	v := db.NewDB()
	err := v.Set("test1", 200, time.Second*5, db.Integer{})
	if err != nil {
		panic(err)
	}

	err2 := v.Set("test2", "hello", time.Second*10, db.String{})
	if err2 != nil {
		panic(err2)
	}
}
```
<h3 align="center">afterwards, you may retrieve them like so:</h3>

```
package main


import (
	"fmt"
	db "github.com/72sevenzy2/in-memory-database"
)

func main() {
	v := db.NewDB()

	// here, db.Integer{} specifies that your retrieving an int-typed value of the given key.
	val, ok := v.Get("test1", db.Integer{})
	if ok {
		fmt.Println(val) // 200
	}

	// and here, db.String{} specifies that your retrieving an string-typed value of the given key.
	val2, ok2 := v.Get("test2", db.String{})
	if ok2 {
		fmt.Println(val2) // "hello"
	}
	
}
```

<h3 align="center">to retrieve all values of type string/int:</h3>

```
package main

import (
	"fmt"
	db "github.com/72sevenzy2/in-memory-database"
)

func main() {
	v := db.NewDB()

	vals, ok := v.GetAll(db.String{}) // vals is of type map[string]string
	if ok {
		for k, v := range vals {
			fmt.Println(k, v)
		}
	}

	vals2, ok2 := v.GetAll(db.Integer{}) // vals2 is of type map[string]uint32
	if ok2 {
		for k, v := range vals2 {
			fmt.Println(k, v)
		}
	}

}
```
<h3 align="center">you can also edit existing keys:</h3>

```
package main

import (
	"fmt"
	db "github.com/72sevenzy2/in-memory-database"
)

func main() {
	v := db.NewDB()

	err := v.Set("test1", 200, db. db.Integer{}) // setting ints
	if err != nil {
		panic(err)
	}

	err := v.Edit("test1", "new value", db.String{}) // "test1" is now set to the string value "new value"
	if err != nil {
		panic(err)
	}
}
```

<h3 align="center">lastly, to delete keys:</h3>

```
package main

import (
	"fmt"
	db "github.com/72sevenzy2/in-memory-database"
)

func main() {
	v := db.NewDB()

	v.Del("keyName")
}

```

<br>
<h1 align="center">interactive cli tutorial:</h1>
<h2 align="center">firstly, switch to the cmd/server directory:</h2>
<h3 align="center">
	<code> cd cmd/server </code>
</h3>
<h2 align="center">and begin with spinning up a db node:</h2>
<h3 align="center">
  <code> go run . -id [node-id] -addr [port] -replicas [port1,port2,port3] </code>
</h3>
<h4 align="center">where -id represents an unique id for the leader node, eg "node-#1", -addr is the port in which it will listen on, and -replicas is for the ports in which the replica nodes will run on (comma-seperated).</h4>
<br>
<h2 align="center">and afterwards, switch to the cmd/cli directory and start the cli (in another terminal): </h2>
<h3 align="center">
	<code> cd cmd/cli </code>
	<br>
	<br>
	<code> go run . </code>
</h3>
<h2 align="center">follow up by declaring a variable:</h2>
<h3 align="center">
  <code> SET [KeyName] [Value] </code>
</h3>
<h2 align="center">with TTL expiration:</h2>
<h3 align="center">
  <code> SET [KeyName] [Value] [TTL, eg: 30s / 1m]</code>
</h3>
<h2 align="center">and retrieve it like so:</h2>
<h3 align="center">
  <code> GET [KeyName] </code>
</h3>
<h2 align="center">you can also edit variables:</h2>
<h3 align="center">
	<code> EDIT [KeyName] [NewValue] </code>
</h3>

<h2 align="center">and to delete a variable/key:</h2>
<h3 align="center">
  <code> DEL [KeyName] </code>
</h3>

<h2 align="center">lastly, to exit:</h2>
<h3 align="center">
  <code> EXIT </code>
</h3>
