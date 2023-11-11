# MOVED TO <https://git.sr.ht/~ashkeel/kilovolt-driver-pebble>

# Pebble driver for Kilovolt

Simple Pebble driver for Kilovolt.

Since Pebble has no fixed versioning data, please refer to the table below for tags to use:

| kv-pebble version | Pebble version tag                 | Kilovolt version |
| ----------------- | ---------------------------------- | ---------------- |
| <= v1.1.1         | v0.0.0-20220127212634-b958d9a7760b | v8.0.5           |
| v1.2.0            | v0.0.0-20221116223310-87eccabb90a3 | v9.0.0           |
| v1.2.1            | v0.0.0-20230209222158-0568b5fd3d14 | v9.0.1           |
| v1.2.2            | v0.0.0-20230418161327-101876aa7088 | v10.0.0          |
| v1.2.3            | v0.0.0-20231102162011-844f0582c2eb | v11.0.0          |

## Usage

Usage is literally a function call to wrap an existing Pebble instance in a Kilovolt driver interface and then passing it over.

```go
package example

import (
	"github.com/cockroachdb/pebble"
	kv "github.com/strimertul/kilovolt/v10"
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
