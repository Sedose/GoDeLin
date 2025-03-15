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

func TestSliceAny(t *testing.T) {
	s := Slice[int]{1, 2, 3, 4, 5}
	if !s.Any(func(i int) bool { return i > 3 }) {
		t.Errorf("expected at least one element to be greater than 3")
	}
	if s.Any(func(i int) bool { return i < 0 }) {
		t.Errorf("expected no element to be less than 0")
	}
}

func TestSliceAnyEmpty(t *testing.T) {
	var s Slice[int]
	if s.Any(func(i int) bool { return i > 0 }) {
		t.Errorf("expected empty slice to return false for Any")
	}
}

func TestAssociate(t *testing.T) {
	s := Slice[int]{1, 2, 3, 4}
	result := Associate(s, func(n int) (int, string) {
		return n, "num" + string(rune('0'+n))
	})
	expected := map[int]string{
		1: "num1",
		2: "num2",
		3: "num3",
		4: "num4",
	}

	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("expected %q for key %d, got %q", v, k, result[k])
		}
	}
}

func TestAssociate_Overwriting(t *testing.T) {
	// When multiple elements map to the same key, later ones overwrite previous ones.
	s := Slice[int]{1, 2, 3, 4}
	result := Associate(s, func(n int) (int, string) {
		return n % 2, "value" + string(rune('0'+n))
	})
	// For key 0: 2 and 4 map to key 0, with 4 overwriting 2.
	// For key 1: 1 and 3 map to key 1, with 3 overwriting 1.
	expected := map[int]string{
		0: "value4",
		1: "value3",
	}

	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("expected %q for key %d, got %q", v, k, result[k])
		}
	}
}
