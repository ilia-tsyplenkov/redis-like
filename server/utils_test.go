package server

import (
	"fmt"
	"testing"
)

func TestMapParser(t *testing.T) {
	slice := []string{"one"}
	if _, err := MapParser(slice); err != fewArgsErr {
		t.Error("slice with len < 2 should not be allowed: %v", slice)
	}
	slice = []string{"one", "two", "three"}
	if _, err := MapParser(slice); err != missValueErr {
		t.Errorf("slice with odd items number shoud not be allowed: %v", slice)
	}
	slice = append(slice, "four")
	got, err := MapParser(slice)
	if err != nil {
		t.Errorf("MapParser(%v) unexpected error: %v", slice, err)
	}
	want := map[string]string{"one": "two", "three": "four"}
	for key := range got {
		if got[key] != want[key] {
			t.Errorf("MapParser(%v) = %v, want: %v", slice, got, want)
		}
	}
}

func TestParamsParser(t *testing.T) {
	slice := []string{}
	if _, _, err := ParamsParser(slice); err != keyNotSpecifiedErr {
		t.Errorf("empty slice %v should not be allowed", slice)
	}

	slice = []string{"key"}
	key, val, err := ParamsParser(slice)
	if err != nil {
		t.Errorf("ParamsParser(%v) error: %v", slice, err)
	}

	if key != "key" {
		t.Errorf("ParamsParser(%v) = %s key, want: 'key'", slice, key)
	}
	if len(val) > 0 {
		t.Errorf("there is no value in the %v slice", slice)
	}

	slice = []string{"key", "val1", "val2"}
	key, val, err = ParamsParser(slice)
	if err != nil {
		t.Errorf("ParamsParser(%v) error: %v", slice, err)
	}
	want := fmt.Sprintf("%v", slice[1:])
	if want != fmt.Sprintf("%v", val) {
		t.Errorf("ParamsParser(%v) = %v, want: %v", slice, val, want)
	}

}

func TestWrongCommandHandler(t *testing.T) {
	var dm DataMap
	cmd := "I am Groot"
	if _, err := CommandHandler(dm, cmd, []string{"key"}); err != unknownCmdErr {
		t.Errorf("got '%v', want: '%v'", err, unknownCmdErr)
	}
}

func TestSetCommandHandler(t *testing.T) {
	var dm DataMap
	cmd := "set"
	dm.Init()
	if _, err := CommandHandler(dm, cmd, []string{"key"}); err != fewArgsErr {
		t.Errorf("got '%v', want: '%v'", err, fewArgsErr)
	}
	if _, err := CommandHandler(dm, cmd, []string{"key", "hello", "world"}); err != manyArgsErr {
		t.Errorf("got '%v', want: '%v'", err, manyArgsErr)
	}
	key := "test"
	dm.hash[key] = &Data{Value: ""}
	if _, err := CommandHandler(dm, cmd, []string{key, "hello"}); err != nil {
		t.Errorf("got '%v', expected 'nil'", err)
	}

}

func TestGetCommnadHandler(t *testing.T) {
	var dm DataMap
	cmd := "get"
	dm.Init()
	if _, err := CommandHandler(dm, cmd, []string{"key", "hello"}); err != manyArgsErr {
		t.Errorf("got '%v', want: '%v'", err, manyArgsErr)
	}
	key := "test"
	dm.hash[key] = &Data{Value: "hello"}
	if _, err := CommandHandler(dm, cmd, []string{key}); err != nil {
		t.Errorf("got '%v', expected 'nil'", err)
	}
}

func TestLSetCommandHandler(t *testing.T) {
	var dm DataMap
	cmd := "lset"
	dm.Init()
	if _, err := CommandHandler(dm, cmd, []string{"key"}); err != fewArgsErr {
		t.Errorf("got '%v', want '%v'", err, fewArgsErr)
	}
	key := "test"
	dm.hash[key] = &Data{Value: []string{}}
	if _, err := CommandHandler(dm, cmd, []string{key, "hello", "world"}); err != nil {
		t.Errorf("got '%v', expected 'nil' error", err)
	}
}

func TestLGetCommandHandler(t *testing.T) {
	var dm DataMap
	dm.Init()
	cmd := "lget"
	if _, err := CommandHandler(dm, cmd, []string{"key", "hello"}); err != manyArgsErr {
		t.Errorf("got '%v', expected '%v' error", err, manyArgsErr)
	}
	key := "test"
	dm.hash[key] = &Data{Value: []string{"hello", "world"}}
	if _, err := CommandHandler(dm, cmd, []string{key}); err != nil {
		t.Errorf("got '%v', expected 'nil' error", err)
	}
}
