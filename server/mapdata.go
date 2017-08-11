package server

import (
	"errors"
	"sync"
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
