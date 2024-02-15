package handlers

import (
	"errors"
	"sync"
)

type keyValueStore map[any]any

var store keyValueStore
var once sync.Once

var lock sync.RWMutex

func getKeyValueStore() keyValueStore {

	once.Do(func() {
		store = make(keyValueStore)
	})

	return store
}

func Write[K comparable, V any](key K, value V) (*V, error) {
	lock.Lock()
	defer lock.Unlock()
	s := getKeyValueStore()
	v, ok := s[key]
	s[key] = value
	if ok {
		typedV, okType := v.(V)
		if okType {
			return &typedV, nil
		}
		return nil, errors.New("WrongType")
	} else {
		return nil, errors.New("NotFound")
	}
}

func Read[V any, K comparable](key K) (*V, error) {
	lock.RLock()
	defer lock.RUnlock()
	s := getKeyValueStore()
	if v, ok := s[key]; ok {
		typedV, okType := v.(V)
		if okType {
			return &typedV, nil
		}
		return nil, errors.New("WrongType")
	} else {
		return nil, errors.New("NotFound")
	}
}
