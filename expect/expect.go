package expect

import (
	"reflect"
	"testing"
)

func Equal[T comparable](t *testing.T, u, v T) {
	if u != v {
		t.Helper()
		t.Errorf("Error: %#v != %#v", u, v)
	}
}

func True(t *testing.T, v bool) {
	if !v {
		t.Helper()
		t.Errorf("Error!")
	}
}

func DeepEqual(t *testing.T, u, v any) {
	if !reflect.DeepEqual(u, v) {
		t.Helper()
		t.Errorf("Error: %#v != %#v", u, v)
	}
}
