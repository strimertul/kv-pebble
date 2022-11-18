# Pebble driver for Kilovolt

Simple Pebble driver for Kilovolt.

## Usage

Usage is literally a function call to wrap an existing Pebble instance in a Kilovolt driver interface and then passing it over.

```go
package example

import (
	"github.com/cockroachdb/pebble"
	kv "github.com/strimertul/kilovolt/v9"
	pebble_driver "github.com/strimertul/kv-pebble"
)

func main() {
	// Initialize your database 
	db, err := pebble.Open("test", &pebble.Options{})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create driver using database instance, second parameter is "sync" (should be kept to true)
	driver := pebble_driver.NewPebbleBackend(db, true)

	// Pass it to Kilovolt
	hub, err := kv.NewHub(driver, kv.HubOptions{}, nil)
	if err != nil {
		panic(err)
	}
	go hub.Run()

	// etc.
}
```