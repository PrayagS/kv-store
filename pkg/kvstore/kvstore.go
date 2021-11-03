package kvstore

import (
	"fmt"
)

type KVStore struct {
	Pairs map[string]string
}

func New() KVStore {
	kvs := KVStore{
		Pairs: make(map[string]string),
	}
	return kvs
}

func (kvs KVStore) Set(key string, value string) {
	kvs.Pairs[key] = value
}

func (kvs KVStore) Get(key string) (string, error) {
	v, ok := kvs.Pairs[key]
	if !ok {
		return v, fmt.Errorf("the key %s does not exist", key)
	}
	return v, nil
}
