package pebble_driver

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cockroachdb/pebble"
)

func TestNewBadgerBackend(t *testing.T) {
	withDB(t, func(db *pebble.DB) {
		NewPebbleBackend(db, true)
	})
}

func TestBackend_Get(t *testing.T) {
	withDB(t, func(db *pebble.DB) {
		err := db.Set([]byte("key"), []byte("value"), nil)
		if err != nil {
			t.Fatal(err)
		}

		backend := NewPebbleBackend(db, true)
		val, err := backend.Get("key")
		if err != nil {
			t.Fatal(err)
		}
		if val != "value" {
			t.Error("Expected value to be 'value'")
		}
	})
}

func TestBackend_Set(t *testing.T) {
	withDB(t, func(db *pebble.DB) {
		backend := NewPebbleBackend(db, true)
		err := backend.Set("key", "value")
		if err != nil {
			t.Fatal(err)
		}
		byt, closer, err := db.Get([]byte("key"))
		if err != nil {
			t.Fatal(err)
		}
		defer closer.Close()
		if string(byt) != "value" {
			t.Error("Expected value to be 'value'")
		}
	})
}

func TestBackend_GetBulk(t *testing.T) {
	kv := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	withDB(t, func(db *pebble.DB) {
		for k, v := range kv {
			err := db.Set([]byte(k), []byte(v), nil)
			if err != nil {
				t.Fatal(err)
			}
		}
		backend := NewPebbleBackend(db, true)
		val, err := backend.GetBulk([]string{"key1", "key2"})
		if err != nil {
			t.Fatal(err)
		}
		if val["key1"] != kv["key1"] || val["key2"] != kv["key2"] {
			t.Error("Expected values do not match")
		}
	})
}

func TestBackend_SetBulk(t *testing.T) {
	kv := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	withDB(t, func(db *pebble.DB) {
		backend := NewPebbleBackend(db, true)
		err := backend.SetBulk(kv)
		if err != nil {
			t.Fatal(err)
		}

		for k, v := range kv {
			byt, closer, err := db.Get([]byte(k))
			if err != nil {
				t.Fatal(err)
			}
			if string(byt) != v {
				t.Error("invalid expected value")
			}
			if err = closer.Close(); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestBackend_GetPrefix(t *testing.T) {
	kv := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	withDB(t, func(db *pebble.DB) {
		for k, v := range kv {
			err := db.Set([]byte(k), []byte(v), nil)
			if err != nil {
				t.Fatal(err)
			}
		}
		backend := NewPebbleBackend(db, true)
		val, err := backend.GetPrefix("key")
		if err != nil {
			t.Fatal(err)
		}
		if len(val) != 2 || val["key1"] != kv["key1"] || val["key2"] != kv["key2"] {
			t.Error("Expected values do not match")
		}
	})
}

func TestBackend_Delete(t *testing.T) {
	withDB(t, func(db *pebble.DB) {
		err := db.Set([]byte("key"), []byte("value"), nil)
		if err != nil {
			t.Fatal(err)
		}
		backend := NewPebbleBackend(db, true)
		err = backend.Delete("key")
		if err != nil {
			t.Fatal(err)
		}
		_, closer, err := db.Get([]byte("key"))
		if err != pebble.ErrNotFound {
			t.Error("Expected key to be deleted")
			if err == nil {
				closer.Close()
			}
		}

	})
}

func TestBackend_List(t *testing.T) {
	kv := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	withDB(t, func(db *pebble.DB) {
		for k, v := range kv {
			err := db.Set([]byte(k), []byte(v), nil)
			if err != nil {
				t.Fatal(err)
			}
		}

		backend := NewPebbleBackend(db, true)
		val, err := backend.List("key")
		if err != nil {
			t.Fatal(err)
		}
		if len(val) != 2 {
			t.Error("Expected 2 keys")
		}
	})
}

func withDB(t *testing.T, test func(db *pebble.DB)) {
	tmp, err := ioutil.TempDir("", "testdb-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	db, err := pebble.Open(tmp, &pebble.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	test(db)
}
