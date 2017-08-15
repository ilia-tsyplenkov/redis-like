package server

import (
	"fmt"
	"testing"
	"time"
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
	if err != typeMismatchErr {
		t.Errorf("got '%v', expected '%v' error", err, typeMismatchErr)
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
	if err != keyNotExistErr {
		t.Fatalf("got '%v', expected '%v' error", err, keyNotExistErr)
	}

	listKey := "list"
	data.hash[listKey] = &Data{Value: []string{"hello", "world"}}
	got, err = data.Get(listKey)
	if err != typeMismatchErr {
		t.Errorf("got '%v', expected '%v' error", err, typeMismatchErr)
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
	if err != typeMismatchErr {
		t.Errorf("got '%v', expected '%v' error", err, typeMismatchErr)
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
	if _, err = dm.LGet(skey); err != typeMismatchErr {
		t.Errorf("got '%v', expected '%v' error", err, typeMismatchErr)
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
	if err := dm.HSet(key, have); err != typeMismatchErr {
		t.Errorf("got '%v', expected '%v' error", err, typeMismatchErr)
	}

}

func TestMapExpire(t *testing.T) {
	key := "test"
	var dm DataMap
	dm.Init()
	dm.hash[key] = &Data{}
	if err := dm.Expire(key, 0); err == nil {
		t.Error("ttl duration must be positive")
	}
	if err := dm.Expire("i'm groot", 10); err != keyNotExistErr {
		t.Errorf("got '%v', want '%v' error", err, keyNotExistErr)
	}
	if err := dm.Expire(key, 1); err != nil {
		t.Fatalf("got '%v', expected 'nil' error", err)
	}
}

func TestMapExpireat(t *testing.T) {
	key := "test"
	var dm DataMap
	dm.Init()
	dm.hash[key] = &Data{}
	if err := dm.Expireat(key, time.Now().UTC().Unix()); err == nil {
		t.Error("ttl should be in the future")
	}
	if err := dm.Expireat("i'm groot", time.Now().UTC().Unix()+10); err != keyNotExistErr {
		t.Errorf("got '%v', want '%v' error", err, keyNotExistErr)
	}
	if err := dm.Expireat(key, time.Now().UTC().Unix()+10); err != nil {
		t.Errorf("got '%v', expected 'nil' error", err)
	}

}

func TestMapTTL(t *testing.T) {
	key := "test"
	var dm DataMap
	dm.Init()
	dm.hash[key] = &Data{}
	if _, err := dm.TTL("i'm groot"); err != keyNotExistErr {
		t.Errorf("got '%v', want '%v' error", err, keyNotExistErr)
	}
	ttl, err := dm.TTL(key)
	if err != nil {
		t.Errorf("got '%v' error, expected 'nil' error", err)
	}
	if ttl != "-1" {
		t.Errorf("got '%s' ttl, expected '-1' ttl with no expire", ttl)
	}
	want := int64(100)
	dm.hash[key] = &Data{ttl: want}
	ttl, err = dm.TTL(key)
	if err != nil {
		t.Errorf("got '%v' error, expected 'nil' error", err)
	}
	if ttl != string(want) {
		t.Errorf("got '%s' ttl, expected '%d' ttl", ttl, want)
	}
}
