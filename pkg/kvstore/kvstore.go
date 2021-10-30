package kvstore

import (
	"fmt"
)

type Value struct {
	Object interface{}
}

type KVStore struct {
	values map[string]Value
}

func New() KVStore {
	kvs := KVStore{
		values: make(map[string]Value),
	}
	return kvs
}

func (kvs KVStore) Set(key string, value interface{}) {
	kvs.values[key] = Value{Object: value}
}

func (kvs KVStore) Get(key string) (Value, error) {
	v, ok := kvs.values[key]
	if !ok {
		return v, fmt.Errorf("the key %s does not exist", key)
	}
	return v, nil
}
