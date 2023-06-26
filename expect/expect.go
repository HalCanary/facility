package expect

import (
	"reflect"
	"testing"
)

// If `u != v`, call `t.Error()` and return false.
func Equal[T comparable](t *testing.T, u, v T) bool {
	if u != v {
		t.Helper()
		t.Errorf("Error: %#v != %#v", u, v)
		return false
	}
	return true
}

// If `!v`, call `t.Error()` and return false.
func True(t *testing.T, v bool) bool {
	if !v {
		t.Helper()
		t.Errorf("Error!")
		return false
	}
	return true
}

// If `!reflect.DeepEqual(u, v)`, call `t.Error()` and return false.
func DeepEqual(t *testing.T, u, v any) bool {
	if !reflect.DeepEqual(u, v) {
		t.Helper()
		t.Errorf("Error: %#v != %#v", u, v)
		return false
	}
	return true
}
