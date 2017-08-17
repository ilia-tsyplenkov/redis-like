package server

import (
	"fmt"
	"testing"
)

func TestMapParser(t *testing.T) {
	slice := []string{"one"}
	if _, err := mapParser(slice); err != fewArgsErr {
		t.Error("slice with len < 2 should not be allowed: %v", slice)
	}
	slice = []string{"one", "two", "three"}
	if _, err := mapParser(slice); err != missValueErr {
		t.Fatalf("slice with odd items number shoud not be allowed: %v", slice)
	}
	slice = append(slice, "four")
	got, err := mapParser(slice)
	if err != nil {
		t.Fatalf("mapParser(%v) unexpected error: %v", slice, err)
	}
	want := map[string]string{"one": "two", "three": "four"}
	for key := range got {
		if got[key] != want[key] {
			t.Fatalf("mapParser(%v) = %v, want: %v", slice, got, want)
		}
	}
}

func TestParamsParser(t *testing.T) {
	slice := []string{}
	if _, _, err := paramsParser(slice); err != keyNotSpecifiedErr {
		t.Fatalf("empty slice %v should not be allowed", slice)
	}

	slice = []string{"key"}
	key, val, err := paramsParser(slice)
	if err != nil {
		t.Fatalf("paramsParser(%v) error: %v", slice, err)
	}

	if key != "key" {
		t.Fatalf("paramsParser(%v) = %s key, want: 'key'", slice, key)
	}
	if len(val) > 0 {
		t.Fatalf("there is no value in the %v slice", slice)
	}

	slice = []string{"key", "val1", "val2"}
	key, val, err = paramsParser(slice)
	if err != nil {
		t.Fatalf("paramsParser(%v) error: %v", slice, err)
	}
	want := fmt.Sprintf("%v", slice[1:])
	if want != fmt.Sprintf("%v", val) {
		t.Fatalf("paramsParser(%v) = %v, want: %v", slice, val, want)
	}

}

func TestWrongDataHandler(t *testing.T) {
	var dm DataMap
	cmd := "I am Groot"
	if _, err := DataHandler(&dm, cmd, []string{"key"}); err != unknownCmdErr {
		t.Fatalf("got '%v', want: '%v'", err, unknownCmdErr)
	}
}

func TestSetDataHandler(t *testing.T) {
	var dm DataMap
	cmd := "set"
	dm.Init()
	if _, err := DataHandler(&dm, cmd, []string{"key"}); err != fewArgsErr {
		t.Fatalf("got '%v', want: '%v'", err, fewArgsErr)
	}
	if _, err := DataHandler(&dm, cmd, []string{"key", "hello", "world"}); err != manyArgsErr {
		t.Fatalf("got '%v', want: '%v'", err, manyArgsErr)
	}
	key := "test"
	dm.hash[key] = &data{value: ""}
	if _, err := DataHandler(&dm, cmd, []string{key, "hello"}); err != nil {
		t.Fatalf("got '%v', expected 'nil'", err)
	}

}

func TestGetDataHandler(t *testing.T) {
	var dm DataMap
	cmd := "get"
	dm.Init()
	if _, err := DataHandler(&dm, cmd, []string{"key", "hello"}); err != manyArgsErr {
		t.Fatalf("got '%v', want: '%v'", err, manyArgsErr)
	}
	key := "test"
	dm.hash[key] = &data{value: "hello"}
	if _, err := DataHandler(&dm, cmd, []string{key}); err != nil {
		t.Fatalf("got '%v', expected 'nil'", err)
	}
}

func TestLSetDataHandler(t *testing.T) {
	var dm DataMap
	cmd := "lset"
	dm.Init()
	if _, err := DataHandler(&dm, cmd, []string{"key"}); err != fewArgsErr {
		t.Fatalf("got '%v', want '%v'", err, fewArgsErr)
	}
	key := "test"
	dm.hash[key] = &data{value: []string{}}
	if _, err := DataHandler(&dm, cmd, []string{key, "hello", "world"}); err != nil {
		t.Fatalf("got '%v', expected 'nil' error", err)
	}
}

func TestLGetDataHandler(t *testing.T) {
	var dm DataMap
	dm.Init()
	cmd := "lget"
	if _, err := DataHandler(&dm, cmd, []string{"key", "hello"}); err != manyArgsErr {
		t.Fatalf("got '%v', expected '%v' error", err, manyArgsErr)
	}
	key := "test"
	dm.hash[key] = &data{value: []string{"hello", "world"}}
	if _, err := DataHandler(&dm, cmd, []string{key}); err != nil {
		t.Fatalf("got '%v', expected 'nil' error", err)
	}
}

func TestLGetItDataHandler(t *testing.T) {
	var dm DataMap
	dm.Init()
	cmd := "lgetit"
	if _, err := DataHandler(&dm, cmd, []string{"key"}); err != fewArgsErr {
		t.Fatalf("got '%v', expected '%v' error for not enough args", err, fewArgsErr)
	}
	if _, err := DataHandler(&dm, cmd, []string{"key", "ind1", "ind2"}); err != manyArgsErr {
		t.Fatalf("got '%v', expected '%v' error for not enough args", err, manyArgsErr)
	}
	key := "test"
	dm.hash[key] = &data{value: []string{"hello", "world"}}
	have := "world"
	got, err := DataHandler(&dm, cmd, []string{key, "1"})
	if err != nil {
		t.Fatalf("got '%v', expected 'nil' erorr for valid case", err)
	}
	if got != have {
		t.Fatalf("got %q, want %q", got, have)
	}

}

func TestLUpdateDataHandler(t *testing.T) {
	var dm DataMap
	dm.Init()
	cmd := "lupdate"
	if _, err := DataHandler(&dm, cmd, []string{"key"}); err != fewArgsErr {
		t.Fatalf("got '%v', expected '%v' error for not enough args", err, fewArgsErr)
	}
	if _, err := DataHandler(&dm, cmd, []string{"key", "ind1", "ind2", "value"}); err != manyArgsErr {
		t.Fatalf("got '%v', expected '%v' error for not enough args", err, manyArgsErr)
	}
	key := "test"
	dm.hash[key] = &data{value: []string{"hello", "world"}}
	have := "bye"
	_, err := DataHandler(&dm, cmd, []string{key, "1", have})
	if err != nil {
		t.Fatalf("got '%v', expected nil error for valid case", err)
	}
	got, _ := dm.LGetIt(key, 1)
	if got != have {
		t.Fatalf("got %q, expected %q", got, have)
	}
}

func TestHSetDataHandler(t *testing.T) {
	var dm DataMap
	cmd := "hset"
	dm.Init()
	if _, err := DataHandler(&dm, cmd, []string{"key", "hello"}); err != fewArgsErr {
		t.Fatalf("got '%v', want '%v'", err, fewArgsErr)
	}
	if _, err := DataHandler(&dm, cmd, []string{"key", "one", "two", "three"}); err != missValueErr {
		t.Fatalf("got '%v', want '%v'", err, missValueErr)
	}
	key := "test"
	dm.hash[key] = &data{value: map[string]string{}}
	if _, err := DataHandler(&dm, cmd, []string{key, "key", "value"}); err != nil {
		t.Fatalf("got '%v', expected 'nil' error", err)
	}
}

func TestHGetDataHandler(t *testing.T) {
	var dm DataMap
	dm.Init()
	cmd := "hget"
	if _, err := DataHandler(&dm, cmd, []string{"key", "hello"}); err != manyArgsErr {
		t.Fatalf("got '%v', expected '%v' error", err, manyArgsErr)
	}
	key := "test"
	dm.hash[key] = &data{value: map[string]string{"key": "value"}}
	if _, err := DataHandler(&dm, cmd, []string{key}); err != nil {
		t.Fatalf("got '%v', expected 'nil' error", err)
	}
}

func TestHGetValDataHandler(t *testing.T) {
	var dm DataMap
	dm.Init()
	cmd := "hgetval"
	if _, err := DataHandler(&dm, cmd, []string{"key"}); err != fewArgsErr {
		t.Fatalf("got '%v', want '%v'", err, fewArgsErr)
	}
	if _, err := DataHandler(&dm, cmd, []string{"outKey", "inKey", "value"}); err != manyArgsErr {
		t.Fatalf("got '%v', want '%v'", err, manyArgsErr)
	}
	key := "test"
	dm.hash[key] = &data{value: map[string]string{"key": "value"}}
	have := "value"
	got, err := DataHandler(&dm, cmd, []string{key, "key"})
	if err != nil {
		t.Fatalf("got '%v', want 'nil' error for valid case", err)
	}
	if got != have {
		t.Fatalf("got %q, expected %q", got, have)
	}
}

func TestHUpdateDataHandler(t *testing.T) {
	var dm DataMap
	dm.Init()
	cmd := "hupdate"
	if _, err := DataHandler(&dm, cmd, []string{"key"}); err != fewArgsErr {
		t.Fatalf("got '%v', want '%v' error for not enough args case", err, fewArgsErr)
	}
	if _, err := DataHandler(&dm, cmd, []string{"outKey", "inKey", "value", "value1"}); err != manyArgsErr {
		t.Fatalf("got '%v', want '%v' error for too many args case", err, manyArgsErr)
	}
	key := "test"
	have := "bye"
	dm.hash[key] = &data{value: map[string]string{"hello": "world"}}
	if _, err := DataHandler(&dm, cmd, []string{key, "hello", have}); err != nil {
		t.Fatalf("got '%v', want 'nil' error for valid case", err)
	}
	got, _ := dm.HGetVal(key, "hello")
	if got != have {
		t.Fatalf("got %q, want %q", got, have)
	}
}
