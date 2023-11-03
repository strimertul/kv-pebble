package pebble_driver

import (
	"github.com/cockroachdb/pebble"
	kv "github.com/strimertul/kilovolt/v10"
)

type Driver struct {
	db   *pebble.DB
	sync bool
}

func NewPebbleBackend(db *pebble.DB, sync bool) Driver {
	return Driver{db, sync}
}

func (b Driver) Get(key string) (string, error) {
	out, closer, err := b.db.Get([]byte(key))
	if err != nil {
		if err == pebble.ErrNotFound {
			return "", kv.ErrorKeyNotFound
		}

		return "", err
	}

	return string(out), closer.Close()
}

func (b Driver) GetBulk(keys []string) (map[string]string, error) {
	out := make(map[string]string)
	for _, key := range keys {
		val, err := b.Get(key)
		if err != nil {
			if err == kv.ErrorKeyNotFound {
				out[key] = ""
				continue
			}
			return nil, err
		}
		out[key] = val
	}

	return out, nil
}

func (b Driver) GetPrefix(prefix string) (map[string]string, error) {
	iter, err := b.db.NewIter(prefixIterOptions([]byte(prefix)))
	if err != nil {
		return nil, err
	}
	out := make(map[string]string)
	for iter.First(); iter.Valid(); iter.Next() {
		out[string(iter.Key())] = string(iter.Value())
	}
	return out, iter.Close()
}

func (b Driver) Set(key, value string) error {
	return b.db.Set([]byte(key), []byte(value), &pebble.WriteOptions{Sync: b.sync})
}

func (b Driver) SetBulk(kv map[string]string) error {
	batch := b.db.NewBatch()
	for key, value := range kv {
		if err := batch.Set([]byte(key), []byte(value), &pebble.WriteOptions{Sync: b.sync}); err != nil {
			return err
		}
	}
	return batch.Commit(&pebble.WriteOptions{Sync: b.sync})
}

func (b Driver) Delete(key string) error {
	return b.db.Delete([]byte(key), &pebble.WriteOptions{Sync: b.sync})
}

func (b Driver) List(prefix string) ([]string, error) {
	iter, err := b.db.NewIter(prefixIterOptions([]byte(prefix)))
	if err != nil {
		return nil, err
	}
	out := []string{}
	for iter.First(); iter.Valid(); iter.Next() {
		out = append(out, string(iter.Key()))
	}
	return out, iter.Close()
}

// from
// https://github.com/cockroachdb/pebble/blob/0ba9163b848ca92495c67f048468466b934df97f/iterator_example_test.go#L50
//
// Copyright 2021 The LevelDB-Go and Pebble Authors. All rights reserved. Use
// of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

func keyUpperBound(b []byte) []byte {
	end := make([]byte, len(b))
	copy(end, b)
	for i := len(end) - 1; i >= 0; i-- {
		end[i] = end[i] + 1
		if end[i] != 0 {
			return end[:i+1]
		}
	}
	return nil // no upper-bound
}

func prefixIterOptions(prefix []byte) *pebble.IterOptions {
	return &pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: keyUpperBound(prefix),
	}
}
