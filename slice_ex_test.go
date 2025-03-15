package godelin

import "testing"

func TestSliceAll(t *testing.T) {
	s := Slice[int]{1, 2, 3, 4, 5}
	if !s.All(func(i int) bool { return i < 10 }) {
		t.Errorf("expected all elements to be less than 10")
	}
	if s.All(func(i int) bool { return i%2 == 0 }) {
		t.Errorf("expected not all elements to be even")
	}
}

func TestSliceAllEmpty(t *testing.T) {
	var s Slice[int]
	if !s.All(func(i int) bool { return i < 10 }) {
		t.Errorf("expected empty slice to return true for All")
	}
}
