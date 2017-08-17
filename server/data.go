package server

import (
	"errors"
)

var typeMismatchErr = errors.New("ERROR: types mismatch")
var noItemErr = errors.New("ERROR: no such item")

type data struct {
	ttl   int64       // time to live
	value interface{} // field for particular data
}

func (d *data) TTL() int64     { return d.ttl }
func (d *data) SetTTL(t int64) { d.ttl = t }

// SSet sets string s to d instance of data.
// It returns error if d contains a value
// with another type.
func (d *data) SSet(s string) error {
	switch d.value.(type) {
	case nil, string:
		d.value = s
		return nil
	default:
		return typeMismatchErr
	}
}

// SGet gets string from d. It returns
// error if d contains value with another type.
func (d *data) SGet() (string, error) {
	switch x := d.value.(type) {
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
func (d *data) LGet() ([]string, error) {
	switch x := d.value.(type) {
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
func (d *data) LSet(list []string) error {
	switch d.value.(type) {
	case nil, []string:
		d.value = list
		return nil
	default:
		return typeMismatchErr
	}
}

// HSet sets map to d. It returns error
// if d contains another type.
func (d *data) HSet(dict map[string]string) error {
	switch d.value.(type) {
	case nil, map[string]string:
		d.value = dict
		return nil
	default:
		return typeMismatchErr
	}
}

// HGet gets a map from d. It returns error
// if d contains another type
func (d *data) HGet() (map[string]string, error) {
	switch x := d.value.(type) {
	case nil:
		return nil, noItemErr
	case map[string]string:
		return x, nil
	default:
		return nil, typeMismatchErr
	}
}
