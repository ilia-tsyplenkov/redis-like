package server

import (
	"errors"
)

var typeMismatchErr = errors.New("types mismatch")
var noItemErr = errors.New("no such item")

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
		return typeMismatchErr
	}
}

func (d *Data) SGet() (string, error) {
	switch x := d.Value.(type) {
	case nil:
		return "", noItemErr
	case string:
		return x, nil
	default:
		return "", typeMismatchErr
	}
}

func (d *Data) LGet() ([]string, error) {
	switch x := d.Value.(type) {
	case nil:
		return nil, noItemErr
	case []string:
		return x, nil
	default:
		return nil, typeMismatchErr
	}
}

func (d *Data) LSet(list []string) error {
	switch d.Value.(type) {
	case nil, []string:
		d.Value = list
		return nil
	default:
		return typeMismatchErr
	}
}

func (d *Data) HSet(dict map[string]string) error {
	switch d.Value.(type) {
	case nil, map[string]string:
		d.Value = dict
		return nil
	default:
		return typeMismatchErr
	}
}

func (d *Data) HGet() (map[string]string, error) {
	switch x := d.Value.(type) {
	case nil:
		return nil, noItemErr
	case map[string]string:
		return x, nil
	default:
		return nil, typeMismatchErr
	}
}
