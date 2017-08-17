package server

import (
	"testing"
)

func TestSet(t *testing.T) {
	s := "hello world"
	var d data
	err := d.SSet(s)
	if err != nil {
		t.Fatalf("d.SSet(%s) error: %v", s, err)
	}
	got, _ := d.value.(string)
	if got != s {
		t.Errorf("got %s, want: %s", got, s)
	}
}

func TestGet(t *testing.T) {
	want := "hello world"
	d := data{value: want}
	got, err := d.SGet()
	if err != nil {
		t.Fatalf("d.SGet error: %v", err)
	}
	if got != want {
		t.Errorf("d.SGet = %s, want: %s", got, want)
	}

}

func TestLSet(t *testing.T) {
	s := []string{"hello", "world"}
	var d data
	err := d.LSet(s)
	if err != nil {
		t.Fatalf("d.LSet(%v) error: %v", s, err)
	}
	got, _ := d.value.([]string)
	if len(got) != len(s) {
		t.Errorf("got %v, want: %v", got, s)
	}
}

func TestLGet(t *testing.T) {
	want := []string{"hello", "world"}
	d := data{value: want}
	got, err := d.LGet()
	if err != nil {
		t.Fatalf("d.LGet() error: %v", err)
	}
	if len(got) != len(want) {
		t.Errorf("d.LGet() = %v, want: %v", got, want)
	}
}

func TestTTL(t *testing.T) {
	ttl := int64(123456)
	d := data{ttl: ttl}
	if d.TTL() != ttl {
		t.Errorf("d.TTL() = %d, want: %d", d.TTL(), ttl)
	}
}

func TestSetTTL(t *testing.T) {
	ttl := int64(123456)
	var d data
	d.SetTTL(ttl)
	if d.ttl != ttl {
		t.Errorf("d.SetTTL put %d but got %d", ttl, d.ttl)
	}
}

func TestHSet(t *testing.T) {
	dict := map[string]string{"one": "hello", "two": "world"}
	var d data
	err := d.HSet(dict)
	if err != nil {
		t.Fatalf("d.HSet(%v) error: %v", dict, err)
	}
	got, _ := d.value.(map[string]string)
	if len(got) != len(dict) {
		t.Errorf("got %v, want: %v", got, dict)
	}
}

func TestHGet(t *testing.T) {
	want := map[string]string{"one": "hello", "two": "world"}
	d := data{value: want}
	got, err := d.HGet()
	if err != nil {
		t.Fatalf("d.HGet() error: %v", err)
	}
	if len(got) != len(want) {
		t.Errorf("d.HGet() = %v, want: %v", got, want)
	}

}
