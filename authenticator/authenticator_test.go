package main

import (
	"reflect"
	"testing"
)

func TestXorBytes(t *testing.T) {
	result := xorBytes([]byte{0}, []byte{0})
	expected := []byte{0}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %#v, but received %#v", expected, result)
	}

	result = xorBytes([]byte{1}, []byte{1})
	expected = []byte{0}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %#v, but received %#v", expected, result)
	}

	result = xorBytes([]byte{0, 1, 1}, []byte{0, 1, 1})
	expected = []byte{0, 0, 0}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %#v, but received %#v", expected, result)
	}

	result = xorBytes([]byte{1, 1, 1}, []byte{0, 0, 0})
	expected = []byte{1, 1, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %#v, but received %#v", expected, result)
	}
}
