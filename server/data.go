package server

import (
	"errors"
)

var typeMismatchErr = errors.New("ERROR: types mismatch")
var noItemErr = errors.New("ERROR: no such item")

type Data struct {
	ttl   int64       // time to live
	Value interface{} // field for particular data
}

func (d *Data) TTL() int64     { return d.ttl }
func (d *Data) SetTTL(t int64) { d.ttl = t }

// SSet sets string s to d instance of Data.
// It returns error if d contains a value
// with another type.
func (d *Data) SSet(s string) error {
	switch d.Value.(type) {
	case nil, string:
		d.Value = s
		return nil
	default:
		return typeMismatchErr
	}
}

// SGet gets string from d. It returns
// error if d contains value with another type.
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

// LGet gets slice from d. It returns
// error if d contains a value with another type.
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

// LSet sets slice list to d. It returns
// if d contains another type.
func (d *Data) LSet(list []string) error {
	switch d.Value.(type) {
	case nil, []string:
		d.Value = list
		return nil
	default:
		return typeMismatchErr
	}
}

// HSet sets map to d. It returns error
// if d contains another type.
func (d *Data) HSet(dict map[string]string) error {
	switch d.Value.(type) {
	case nil, map[string]string:
		d.Value = dict
		return nil
	default:
		return typeMismatchErr
	}
}

// HGet gets a map from d. It returns error
// if d contains another type
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
