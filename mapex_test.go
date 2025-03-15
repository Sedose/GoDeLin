package godelin

import (
	"strconv"
	"testing"
)

// TestMapExSlice demonstrates storing a slice in a MapEx,
// appending items to it, and confirming that subsequent calls
// to GetOrPut for the same key retrieve the existing slice.
func TestMapExSlice(t *testing.T) {
	m := NewMapEx[string, []int]()

	slice := m.GetOrPut("nums", func() []int { return nil })
	slice = append(slice, 42)
	m.data["nums"] = slice

	sliceAgain := m.GetOrPut("nums", func() []int { return nil })
	sliceAgain = append(sliceAgain, 7)
	m.data["nums"] = sliceAgain

	finalSlice := m.GetOrPut("nums", func() []int { return nil })
	if len(finalSlice) != 2 || finalSlice[0] != 42 || finalSlice[1] != 7 {
		t.Errorf("expected [42 7], got %v", finalSlice)
	}
}

func TestMapExAppendingMultipleKeys(t *testing.T) {
	m := NewMapEx[string, []string]()

	valA := m.GetOrPut("A", func() []string { return []string{} })
	valA = append(valA, "apple")
	m.data["A"] = valA

	valB := m.GetOrPut("B", func() []string { return []string{} })
	valB = append(valB, "banana")
	m.data["B"] = valB

	valA2 := m.GetOrPut("A", func() []string { return []string{} })
	valA2 = append(valA2, "avocado")
	m.data["A"] = valA2

	finalA := m.GetOrPut("A", func() []string { return []string{} })
	if len(finalA) != 2 || finalA[0] != "apple" || finalA[1] != "avocado" {
		t.Errorf("expected [apple avocado], got %v", finalA)
	}
	finalB := m.GetOrPut("B", func() []string { return []string{} })
	if len(finalB) != 1 || finalB[0] != "banana" {
		t.Errorf("expected [banana], got %v", finalB)
	}
}

// Below are the basic tests from before, ensuring we also check simpler use cases.

func TestMapExIntString(t *testing.T) {
	m := NewMapEx[int, string]()

	calls := 0
	defValFunc := func() string {
		calls++
		return "default"
	}

	if v := m.GetOrPut(1, defValFunc); v != "default" {
		t.Errorf("expected default, got %v", v)
	}
	if calls != 1 {
		t.Errorf("expected defValFunc to be called once, got %d", calls)
	}
	if v := m.GetOrPut(1, defValFunc); v != "default" {
		t.Errorf("expected default, got %v", v)
	}
	if calls != 1 {
		t.Errorf("expected defValFunc not to be called again, got %d", calls)
	}
}

func TestMapExStringInt(t *testing.T) {
	m := NewMapEx[string, int]()

	if v := m.GetOrPut("keyA", func() int { return 123 }); v != 123 {
		t.Errorf("expected 123, got %v", v)
	}
	if v := m.GetOrPut("keyA", func() int { return 456 }); v != 123 {
		t.Errorf("expected 123, got %v", v)
	}
}

func TestMapExMultipleKeys(t *testing.T) {
	m := NewMapEx[int, int]()

	values := []int{10, 20, 30}
	for i := range values {
		res := m.GetOrPut(i, func() int { return values[i] })
		if res != values[i] {
			t.Errorf("expected %d, got %d", values[i], res)
		}
	}
	for i := range values {
		res := m.GetOrPut(i, func() int { return 999 })
		if res != values[i] {
			t.Errorf("expected %d, got %d", values[i], res)
		}
	}
}

type customKey struct {
	X, Y int
}

func TestMapExStructKey(t *testing.T) {
	m := NewMapEx[customKey, string]()

	key := customKey{X: 2, Y: 3}
	res := m.GetOrPut(key, func() string { return "val" })
	if res != "val" {
		t.Errorf("expected val, got %v", res)
	}
	res = m.GetOrPut(key, func() string { return "other" })
	if res != "val" {
		t.Errorf("expected val, got %v", res)
	}
}

func TestMapExFuncAsDefaultValue(t *testing.T) {
	m := NewMapEx[int, func() string]()
	defFunc := func() func() string {
		return func() string { return "result" }
	}
	res := m.GetOrPut(5, defFunc)
	if res() != "result" {
		t.Errorf("expected result, got %v", res())
	}
	res2 := m.GetOrPut(5, func() func() string {
		return func() string { return "changed" }
	})
	if res2() != "result" {
		t.Errorf("expected original result, got %v", res2())
	}
}

func TestMapExDifferentCalls(t *testing.T) {
	m := NewMapEx[string, int]()
	call1 := m.GetOrPut("A", func() int { return 1 })
	call2 := m.GetOrPut("A", func() int { return 2 })
	if call1 != 1 || call2 != 1 {
		t.Errorf("expected repeated calls to return 1, got %v and %v", call1, call2)
	}
}

func TestMapExIncrement(t *testing.T) {
	m := NewMapEx[string, int]()
	for i := 0; i < 5; i++ {
		m.GetOrPut("counter", func() int { return 0 })
		m.data["counter"]++
	}
	if m.data["counter"] != 5 {
		t.Errorf("expected 5, got %v", m.data["counter"])
	}
}

func TestMapExGetOrPutWithStringConversion(t *testing.T) {
	m := NewMapEx[int, string]()
	for i := 0; i < 3; i++ {
		val := m.GetOrPut(i, func() string { return "val" + strconv.Itoa(i) })
		if val != "val"+strconv.Itoa(i) {
			t.Errorf("expected %s, got %s", "val"+strconv.Itoa(i), val)
		}
	}
}
