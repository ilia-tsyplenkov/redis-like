package server

import (
	"fmt"
)

type Data struct {
	ttl   int64
	Value interface{}
}

func (d *Data) TTL() int64     { return d.ttl }
func (d *Data) SetTTL(t int64) { d.ttl = t }

func (d *Data) SSet(s string) error {
	switch d.Value.(type) {
	case nil, string:
		d.Value = s
		return nil
	default:
		return fmt.Errorf("types mismatch")
	}
}

func (d *Data) SGet() (string, error) {
	switch x := d.Value.(type) {
	case nil:
		return "", fmt.Errorf("no such item")
	case string:
		return x, nil
	default:
		return "", fmt.Errorf("wrong type")
	}
}

func (d *Data) LGet() ([]string, error) {
	switch x := d.Value.(type) {
	case nil:
		return nil, fmt.Errorf("no such item")
	case []string:
		return x, nil
	default:
		return nil, fmt.Errorf("wrong type")
	}
}

func (d *Data) LSet(list []string) error {
	switch d.Value.(type) {
	case nil, []string:
		d.Value = list
		return nil
	default:
		return fmt.Errorf("types mismatch")
	}
}

func (d *Data) HSet(dict map[string]string) error {
	switch d.Value.(type) {
	case nil, map[string]string:
		d.Value = dict
		return nil
	default:
		return fmt.Errorf("types mismatch")
	}
}

func (d *Data) HGet() (map[string]string, error) {
	switch x := d.Value.(type) {
	case nil:
		return nil, fmt.Errorf("no such item")
	case map[string]string:
		return x, nil
	default:
		return nil, fmt.Errorf("wrong type")
	}
}
