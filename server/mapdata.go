package server

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var keyNotExistErr = errors.New("key not exists")

type DataMap struct {
	mu   sync.RWMutex
	hash map[string]*Data
}

func (d *DataMap) Init() {
	d.hash = make(map[string]*Data)
}

func (d *DataMap) Set(key, val string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.hash[key]; !ok {
		d.hash[key] = new(Data)
	}
	return d.hash[key].SSet(val)
}

func (d *DataMap) Get(key string) (string, error) {
	d.mu.RLock()
	val, ok := d.hash[key]
	d.mu.RUnlock()
	if !ok {
		return "", keyNotExistErr
	}
	return val.SGet()
}

func (d *DataMap) LSet(key string, val []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.hash[key]; !ok {
		d.hash[key] = new(Data)
	}
	return d.hash[key].LSet(val)
}

func (d *DataMap) LGet(key string) ([]string, error) {
	d.mu.RLock()
	val, ok := d.hash[key]
	d.mu.RUnlock()
	if !ok {
		return nil, keyNotExistErr
	}
	return val.LGet()
}

func (d *DataMap) HSet(key string, val map[string]string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.hash[key]; !ok {
		d.hash[key] = new(Data)
	}
	return d.hash[key].HSet(val)
}

func (d *DataMap) HGet(key string) (map[string]string, error) {
	d.mu.RLock()
	val, ok := d.hash[key]
	d.mu.RUnlock()
	if !ok {
		return nil, keyNotExistErr
	}
	return val.HGet()
}

func (dm *DataMap) Keys() []string {
	var keys []string
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	for key := range dm.hash {
		keys = append(keys, key)
	}
	return keys
}

func (dm *DataMap) Remove(key string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	delete(dm.hash, key)
}

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
