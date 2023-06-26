package expect // import "github.com/HalCanary/facility/expect"


FUNCTIONS

func DeepEqual(t *testing.T, u, v any) bool
    If `!reflect.DeepEqual(u, v)`, call `t.Error()` and return false.

func Equal[T comparable](t *testing.T, u, v T) bool
    If `u != v`, call `t.Error()` and return false.

func True(t *testing.T, v bool) bool
    If `!v`, call `t.Error()` and return false.

