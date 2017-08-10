package server

import (
	"fmt"
	"testing"
)

func TestMapSet(t *testing.T) {
	key := "test"
	want := "hello world"
	var data DataMap
	data.Init()
	err := data.Set(key, want)
	if err != nil {
		t.Fatalf("Set(%s, %s) error: %v", key, want, err)
	}
	got := fmt.Sprintf("%v", data.hash[key].Value)
	if got != want {
		t.Fatalf("got %v, want: %v", got, want)
	}
	mapVal := map[string]string{"hello": "world"}
	data.hash[key] = &Data{Value: mapVal}
	err = data.Set(key, want)
	if err == nil {
		t.Errorf("Set for another type returns nil error")
	}
}

func TestMapGet(t *testing.T) {
	var data DataMap
	key := "test"
	want := "hello world"
	data.Init()
	data.hash[key] = &Data{Value: want}
	got, err := data.Get(key)
	if err != nil {
		t.Fatalf("Get(%s) error: %v", key, err)
	}
	if got != want {
		t.Fatalf("Get(%s) = %v, want: %v", key, got, want)
	}
	got, err = data.Get("foo")
	if err == nil {
		t.Fatalf("Get for non-existing key returns nil error")
	}

	listKey := "list"
	data.hash[listKey] = &Data{Value: []string{"hello", "world"}}
	got, err = data.Get(listKey)
	if err == nil {
		t.Errorf("Get for another type returns nil error")
	}

}

func TestMapLSet(t *testing.T) {
	key := "test"
	have := []string{"hello", "world"}
	var data DataMap
	data.Init()
	err := data.LSet(key, have)
	if err != nil {
		t.Fatalf("LSet(%s, %v) error: %v", key, have, err)
	}
	got := fmt.Sprintf("%v", data.hash[key].Value)
	if got != fmt.Sprintf("%v", have) {
		t.Fatalf("got %s, want: %v", got, have)
	}
	mapVal := map[string]string{"hello": "world"}
	data.hash[key] = &Data{Value: mapVal}
	err = data.LSet(key, have)
	if err == nil {
		t.Errorf("LSet for another type returns nil error")
	}
}

func TestMapLGet(t *testing.T) {
	var dm DataMap
	key := "test"
	have := []string{"hello", "world"}
	dm.Init()
	dm.hash[key] = &Data{Value: have}
	got, err := dm.LGet(key)
	if err != nil {
		t.Fatalf("LGet(%s) error: %v", key, err)
	}
	if fmt.Sprintf("%v", have) != fmt.Sprintf("%v", got) {
		t.Fatalf("Get(%s) = %v, want: %v", key, got, have)
	}
	skey := "str"
	dm.hash[skey] = &Data{Value: "string"}
	if _, err = dm.LGet(skey); err == nil {
		t.Errorf("LGet for another type returns nil error")
	}
}

func TestMapHSet(t *testing.T) {
	key := "test"
	have := map[string]string{"hello": "world"}
	var dm DataMap
	dm.Init()
	if err := dm.HSet(key, have); err != nil {
		t.Fatalf("HSet(%s, %v) error: %v", key, have, err)
	}
	got := fmt.Sprintf("%v", dm.hash[key].Value)
	if len(fmt.Sprintf("%v", have)) != len(got) {
		t.Fatalf("got %s, want: %v", got, have)
	}
	lst := []string{"hello", "world"}
	dm.hash[key] = &Data{Value: lst}
	if err := dm.HSet(key, have); err == nil {
		t.Errorf("HSet fot another type returns nil error")
	}

}
