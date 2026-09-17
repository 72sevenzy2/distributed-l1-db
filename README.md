<h1 align="center"> key-value style l1 db. </h1>
<br>
<ul>
  <li>TTLs for data with limited life-span.</li>
  <li>utilises a node-based communication system with a leader election mechanism if it were to fail.</li>
  <li>persistant serverside error handling.</li>
  <li>interactive cli mode, which stores variable-like data  to then be retrieved later with methods like "GET", "SET", "DEL", and "EXIT" to exit the program.</li>
  <li>serializes values to raw bytes before appending to the in-memory map for optimised performance upon node-to-node communication.</li>
  <li>listens on tcp server for database logic and validation.</li>
</ul>

<h1 align="center">usage:</h1>
<h3 align="center">to start off, you can either set data using SetInt() for key-values with values having type int, or SetString() for key-values with values having type string:</h3>

```
package main

import (
	"fmt"
	"github.com/72sevenzy2/in-memory-database/db"
)

func main() {
	v := db.NewDB()

	err := v.SetInt("test1", 200) // setting ints
	if err != nil {
		panic(err)
	}

	err2 := v.SetString("test2", "hello") // setting strings
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
	"github.com/72sevenzy2/in-memory-database/db"
)

func main() {
	v := db.NewDB()

	val, ok := v.GetInt("test1")
	if ok {
		fmt.Println(val)
	}

	val2, ok2 := v.GetString("test2")
	if ok2 {
		fmt.Println(val2)
	}
	
}
```

<h3 align="center">if you want to retrieve all key-values with values of type strings/ints:</h3>

```
package main

import (
	"fmt"
	"github.com/72sevenzy2/in-memory-database/db"
)

func main() {
	v := db.NewDB()

	vals, ok := v.GetAllString() // vals is of type map[string]string
	if ok {
		for k, v := range vals {
			fmt.Println(k, v)
		}
	}

	vals2, ok2 := v.GetAllInt() // vals2 is of type map[string]uint32
	if ok2 {
		for k, v := range vals2 {
			fmt.Println(k, v)
		}
	}

}
```

<h3 align="center">and finally, to delete keys:</h3>

```
package main

import (
	"fmt"
	"github.com/72sevenzy2/in-memory-database/db"
)

func main() {
	v := db.NewDB()

	v.Del("keyName")

}

```

<br>
<h1 align="center">interactive cli tutorial:</h1>
<h2 align="center">run the following to begin:</h2>
<h3 align="center">
  <code> go run . -id node-id -addr port -replicas port1,port2,port3 </code>
</h3>
<h4 align="center">-id represents an unique id for the leader node, eg "node-#1", -addr is the port in which it will listen on, and -replicas is for the ports in which the replica nodes will run on (comma-seperated).</h4>
<br>
<h2 align="center">follow up by declaring a variable:</h2>
<h3 align="center">
  <code> SET [KeyName] [Value] </code>
</h3>
<h3 align="center">(key names only support either integers or strings as of now).</h3>

<h2 align="center">after setting a variable, you can then retrieve it like so:</h2>
<h3 align="center">
  <code> GET [KeyName] </code>
</h3>
<h3 align="center">and it returns the value of the key given.</h3>

<h2 align="center">to delete a variable/key:</h2>
<h3 align="center">
  <code> DEL [KeyName] </code>
</h3>

<h2 align="center">finally, to exit:</h2>
<h3 align="center">
  <code> EXIT </code>
</h3>
