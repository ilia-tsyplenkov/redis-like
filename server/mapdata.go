// package server provides API for building memcache server.
package server

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var keyNotExistErr = errors.New("ERROR: key not exists")
var invalidIndexErr = errors.New("ERROR: invalid list index")
var invalidInnerKeyErr = errors.New("ERROR: invalid inner key")

type DataMap struct {
	DbId string
	mu   sync.RWMutex
	hash map[string]*data
}

// Init initializes hash map in dm.
func (dm *DataMap) Init() {
	dm.hash = make(map[string]*data)
}

// Set sets string in dm by key.
// Returns error if key contains another type.
func (dm *DataMap) Set(key, val string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	if _, ok := dm.hash[key]; !ok {
		dm.hash[key] = new(data)
	}
	return dm.hash[key].SSet(val)
}

// Get gets string from dm by key.
// Returns error if key not exists or
// contains another type.
func (dm *DataMap) Get(key string) (string, error) {
	dm.mu.RLock()
	val, ok := dm.hash[key]
	dm.mu.RUnlock()
	if !ok {
		return "", keyNotExistErr
	}
	return val.SGet()
}

// LSet sets list to dm by key.
// Returns error if key contains another type.
func (dm *DataMap) LSet(key string, val []string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	if _, ok := dm.hash[key]; !ok {
		dm.hash[key] = new(data)
	}
	return dm.hash[key].LSet(val)
}

// LGet gets slice from dm by key.
// Returns error if key not exists or
// contains another type.
func (dm *DataMap) LGet(key string) ([]string, error) {
	dm.mu.RLock()
	val, ok := dm.hash[key]
	dm.mu.RUnlock()
	if !ok {
		return nil, keyNotExistErr
	}
	return val.LGet()
}

// LGetIt gets slice from dm by key and
// get item from slice by index.
// Returns error if key or index is invalid
// and if key contains another type
func (dm *DataMap) LGetIt(key string, index int) (string, error) {
	s, err := dm.LGet(key)
	if err != nil {
		return "", err
	}
	if index < 0 || index >= len(s) {
		return "", invalidIndexErr
	}
	return s[index], nil

}

// LUpdate updates slice in dm by key.
// It changes data by index to value.
// Returns error if key or index is invalid
// and if key contains another type.
func (dm *DataMap) LUpdate(key string, index int, value string) error {
	s, err := dm.LGet(key)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(s) {
		return invalidIndexErr
	}
	s[index] = value
	return nil
}

// HSet sets map to dm by key.
// Returns error if key contains
// another type.
func (dm *DataMap) HSet(key string, val map[string]string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	if _, ok := dm.hash[key]; !ok {
		dm.hash[key] = new(data)
	}
	return dm.hash[key].HSet(val)
}

// HGet gets map from dm by key.
// Returns error if key not exists or
// contains another type.
func (dm *DataMap) HGet(key string) (map[string]string, error) {
	dm.mu.RLock()
	val, ok := dm.hash[key]
	dm.mu.RUnlock()
	if !ok {
		return nil, keyNotExistErr
	}
	return val.HGet()
}

// HGetVal gets map from dm by outerKey and then
// gets a value from map by innerKey.
// Returns error if outerKey not exists or
// outerKey contains another type or
// innerKey not exists.
func (dm *DataMap) HGetVal(outerKey, innerKey string) (string, error) {
	dict, err := dm.HGet(outerKey)
	if err != nil {
		return "", err
	}
	val, ok := dict[innerKey]
	if !ok {
		return "", invalidInnerKeyErr
	}
	return val, nil
}

// HUpdate gets map from dm by outKey and then
// updates inKey value. Returns error if outKey
// not exists.
func (dm *DataMap) HUpdate(outKey, inKey, value string) error {
	dict, err := dm.HGet(outKey)
	if err != nil {
		return err
	}
	dict[inKey] = value
	return nil
}

// Keys gets all keys in dm.
func (dm *DataMap) Keys() []string {
	var keys []string
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	for key := range dm.hash {
		keys = append(keys, key)
	}
	return keys
}

// Remove deletes key from dm.
func (dm *DataMap) Remove(key string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	delete(dm.hash, key)
}

// Expire sets ttl for key in dm as sum of
// current unix time + dur. Returns error
// if key not exists.
func (dm *DataMap) Expire(key string, dur int64) error {
	if dur <= 0 {
		return fmt.Errorf("ERROR: wrong duration for ttl")
	}
	ttl := time.Now().UTC().Unix() + dur
	dm.mu.Lock()
	defer dm.mu.Unlock()
	data, ok := dm.hash[key]
	if !ok {
		return keyNotExistErr
	}
	data.SetTTL(ttl)
	return nil
}

// Expireat sets ttl for key in dm.
// Returns error if key not exists.
func (dm *DataMap) Expireat(key string, ttl int64) error {
	if ttl <= time.Now().UTC().Unix() {
		return fmt.Errorf("ERROR: ttl must be greater than now")
	}
	dm.mu.Lock()
	defer dm.mu.Unlock()
	data, ok := dm.hash[key]
	if !ok {
		return keyNotExistErr
	}
	data.SetTTL(ttl)
	return nil
}

// Persist reset ttl to default for data in dm by key.
// Returns error if key not exists.
func (dm *DataMap) Persist(key string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	data, ok := dm.hash[key]
	if !ok {
		return keyNotExistErr
	}
	data.SetTTL(0)
	return nil

}

// TTL gets ttl from dm by key.
// Returns error if key not exists.
func (dm *DataMap) TTL(key string) (string, error) {
	dm.mu.RLock()
	data, ok := dm.hash[key]
	dm.mu.RUnlock()
	if !ok {
		return "", keyNotExistErr
	}
	ttl := data.TTL()
	if ttl > 0 {
		return fmt.Sprintf("%d", ttl), nil
	}
	return "-1", nil
}
