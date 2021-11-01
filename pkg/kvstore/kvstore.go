package kvstore

import (
	"fmt"
)

type KVStore struct {
	values map[string]string
}

func New() KVStore {
	kvs := KVStore{
		values: make(map[string]string),
	}
	return kvs
}

func (kvs KVStore) Set(key string, value string) {
	kvs.values[key] = value
}

func (kvs KVStore) Get(key string) (string, error) {
	v, ok := kvs.values[key]
	if !ok {
		return v, fmt.Errorf("the key %s does not exist", key)
	}
	return v, nil
}
