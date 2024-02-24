package handlers

import (
	"com/mapify/structs"
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

func Populate() {
	contacts := [...]structs.Contact{
		{Id: "1", First: "John", Last: "Smith", Phone: "123-456-7890", Email: "John@example.com"},
		{Id: "2", First: "Dana", Last: "Crandith", Phone: "123-456-7890", Email: "dana@example.com"},
		{Id: "3", First: "Edith", Last: "Neutvaar", Phone: "123-456-7890", Email: "edith@example.com"}}

	Write("contacts", contacts[:])
}

func Write[K comparable, V any](key K, value V) *V {
	lock.Lock()
	defer lock.Unlock()
	s := getKeyValueStore()
	v, ok := s[key]
	s[key] = value
	if ok {
		typedV, okType := v.(V)
		if okType {
			return &typedV
		}
		return nil
	} else {
		return nil
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
