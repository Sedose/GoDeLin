package godelin

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestAll(t *testing.T) {
	type args struct {
		elems []int
		fn    func(int) bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"positive case",
			args{
				[]int{1, 2, 3, 4, 5},
				func(i int) bool { return i < 10 },
			},
			true,
		},
		{"negative case",
			args{
				[]int{1, 2, 3, 4, 5},
				func(i int) bool { return i%2 == 0 },
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := All(tt.args.elems, tt.args.fn); got != tt.want {
				t.Errorf("All() = %v, expected %v", got, tt.want)
			}
		})
	}
}

// Test when key is already present in the map
func TestGetOrPut_ExistingKey(t *testing.T) {
	m := map[string]int{"a": 42}

	got := GetOrPut(m, "a", func() int { return 100 })
	expected := 42

	if got != expected {
		t.Errorf("GetOrPut() returned %v, expected %v", got, expected)
	}
}

// Test when key is missing, default value should be inserted
func TestGetOrPut_MissingKey(t *testing.T) {
	m := map[string]int{}

	got := GetOrPut(m, "b", func() int { return 100 })
	expected := 100

	if got != expected {
		t.Errorf("GetOrPut() returned %v, expected %v", got, expected)
	}

	// Check if the key was actually inserted
	if m["b"] != 100 {
		t.Errorf("GetOrPut() did not insert the expected value into the map")
	}
}

// Test GetOrPut with slice values
func TestGetOrPut_SliceValue(t *testing.T) {
	m := map[string][]int{}

	got := GetOrPut(m, "c", func() []int { return []int{1, 2, 3} })
	expected := []int{1, 2, 3}

	if len(got) != len(expected) {
		t.Errorf("GetOrPut() returned slice of length %d, expected %d", len(got), len(expected))
	}

	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("GetOrPut() returned %v, expected %v", got, expected)
		}
	}

	// Check if map was updated correctly
	if len(m["c"]) != 3 {
		t.Errorf("GetOrPut() did not insert the expected slice into the map")
	}
}

func TestGetOrPut_StructValue(t *testing.T) {
	type TestStruct struct {
		Value string
	}
	m := map[int]TestStruct{}

	got := GetOrPut(m, 10, func() TestStruct { return TestStruct{Value: "Hello"} })
	expected := TestStruct{Value: "Hello"}

	if got != expected {
		t.Errorf("GetOrPut() returned %v, expected %v", got, expected)
	}

	// Check if map was updated
	if m[10] != expected {
		t.Errorf("GetOrPut() did not insert the expected struct into the map")
	}
}

func TestAny(t *testing.T) {
	type args struct {
		elems []int
		fn    func(int) bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"positive case",
			args{
				[]int{1, 2, 3, 4, 5, 6},
				func(i int) bool { return i%2 == 0 },
			},
			true,
		},
		{"negative case",
			args{
				[]int{1, 2, 3, 4, 5, 6},
				func(i int) bool { return i > 7 },
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Any(tt.args.elems, tt.args.fn); got != tt.want {
				t.Errorf("Any() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	testCases := []struct {
		name      string
		slice     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "filter even numbers",
			slice:     []int{1, 2, 3, 4, 5, 6},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  []int{2, 4, 6},
		},
		{
			name:      "no elements match predicate",
			slice:     []int{1, 3, 5},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  []int{},
		},
		{
			name:      "all elements match predicate",
			slice:     []int{2, 4, 6},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  []int{2, 4, 6},
		},
		{
			name:      "empty slice",
			slice:     []int{},
			predicate: func(x int) bool { return x > 0 },
			expected:  []int{},
		},
		{
			name:      "predicate matches only first element",
			slice:     []int{9, 1, 2, 3},
			predicate: func(x int) bool { return x == 9 },
			expected:  []int{9},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := Filter(testCase.slice, testCase.predicate)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Filter() = %v, expected %v", actual, testCase.expected)
			}
		})
	}
}

func TestAssociate_StringGrouping(t *testing.T) {
	// Group strings by their first letter.
	fruits := []string{"apple", "apricot", "banana", "avocado", "blueberry"}
	result := GroupBy(fruits, func(fruit string) (string, string) {
		return fruit[:1], fruit
	})
	expected := map[string][]string{
		"a": {"apple", "apricot", "avocado"},
		"b": {"banana", "blueberry"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Associate() = %v, expected %v", result, expected)
	}
}

func TestAssociate_IntGrouping(t *testing.T) {
	// Group integers by even/odd.
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := GroupBy(numbers, func(n int) (bool, int) {
		return n%2 == 0, n
	})
	expected := map[bool][]int{
		false: {1, 3, 5, 7, 9},
		true:  {2, 4, 6, 8, 10},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Associate() with ints = %v, expected %v", result, expected)
	}
}

func TestAssociate_EmptySlice(t *testing.T) {
	// Edge case: empty slice.
	var empty []string
	result := GroupBy(empty, func(s string) (string, string) {
		return s, s
	})
	expected := map[string][]string{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Associate() with empty slice = %v, expected %v", result, expected)
	}
}

func TestChunked(t *testing.T) {
	type args struct {
		s []int
		n int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{"exact multiple",
			args{
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				3,
			},
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		},
		{"extra elements",
			args{
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				2,
			},
			[][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9}},
		},
		{"not enough elements",
			args{
				[]int{1, 2, 3},
				5,
			},
			[][]int{{1, 2, 3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Chunked(tt.args.s, tt.args.n); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("Chunked() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestChunkedBy(t *testing.T) {
	input := []int{
		10,
		20,
		30,
		40,
		31,
		31,
		33,
		34,
		21,
		22,
		23,
		24,
		11,
		12,
		13,
		14,
	}
	output := ChunkedBy(input, func(prev, next int) bool {
		return prev < next
	})
	expected := [][]int{
		{10, 20, 30, 40},
		{31},
		{31, 33, 34},
		{21, 22, 23, 24},
		{11, 12, 13, 14},
	}
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("ChunkedBy() = %v, expected %v", output, expected)
	}
}

func TestDistinct(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"test distinct",
			args{
				[]int{1, 1, 2, 3, 3, 4, 4, 4, 4, 5, 5, 5},
			},
			[]int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Distinct(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Distinct() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestDistinctBy(t *testing.T) {
	type args struct {
		s  []string
		fn func(string) string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{{"test distinctBy",
		args{[]string{"a", "A", "b", "B", "c", "C"},
			func(s string) string { return strings.ToLower(s) },
		},
		[]string{"a", "b", "c"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DistinctBy(tt.args.s, tt.args.fn); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("DistinctBy() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestDrop(t *testing.T) {
	type args struct {
		elements []int
		n        int
	}
	testCases := []struct {
		name     string
		args     args
		expected []int
	}{
		{
			name: "drop less than slice length",
			args: args{
				elements: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				n:        4,
			},
			expected: []int{5, 6, 7, 8, 9},
		},
		{
			name: "drop more than slice length",
			args: args{
				elements: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				n:        10,
			},
			expected: nil, // Expect nil when n >= len(elements)
		},
		{
			name: "drop exactly slice length",
			args: args{
				elements: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				n:        9,
			},
			expected: nil, // Expect nil when n >= len(elements)
		},
		{
			name: "drop zero elements",
			args: args{
				elements: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				n:        0,
			},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "drop from empty slice",
			args: args{
				elements: []int{},
				n:        3,
			},
			expected: []int{}, // Expect nil when slice is empty
		},
		{
			name: "drop negative number of elements",
			args: args{
				elements: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				n:        -1,
			},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := Drop(testCase.args.elements, testCase.args.n)
			if !reflect.DeepEqual(got, testCase.expected) {
				t.Errorf("Drop() = %v, expected %v", got, testCase.expected)
			}
		})
	}
}

func TestDropLast(t *testing.T) {
	testCases := []struct {
		name      string
		input     []int
		numToDrop int
		expected  []int
	}{
		{
			name:      "drop fewer than slice length",
			input:     []int{1, 2, 3, 4, 5},
			numToDrop: 2,
			expected:  []int{1, 2, 3},
		},
		{
			name:      "drop exactly slice length",
			input:     []int{1, 2, 3},
			numToDrop: 3,
			expected:  nil,
		},
		{
			name:      "drop more than slice length",
			input:     []int{1, 2, 3},
			numToDrop: 5,
			expected:  nil,
		},
		{
			name:      "drop zero elements",
			input:     []int{1, 2, 3},
			numToDrop: 0,
			expected:  []int{1, 2, 3},
		},
		{
			name:      "drop negative number of elements",
			input:     []int{1, 2, 3},
			numToDrop: -1,
			expected:  []int{1, 2, 3},
		},
		{
			name:      "empty slice",
			input:     []int{},
			numToDrop: 3,
			expected:  []int{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := DropLast(testCase.input, testCase.numToDrop)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("DropLast() = %v, expected %v", actual, testCase.expected)
			}
		})
	}
}

func TestDropLastWhile(t *testing.T) {
	testCases := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "drop elements greater than 3 at end",
			input:     []int{1, 2, 5, 4},
			predicate: func(x int) bool { return x > 3 },
			expected:  []int{1, 2},
		},
		{
			name:      "no elements match predicate",
			input:     []int{1, 2, 3},
			predicate: func(x int) bool { return x > 5 },
			expected:  []int{1, 2, 3},
		},
		{
			name:      "all elements match predicate",
			input:     []int{5, 6, 7},
			predicate: func(x int) bool { return x > 1 },
			expected:  []int{},
		},
		{
			name:      "empty slice",
			input:     []int{},
			predicate: func(x int) bool { return x > 1 },
			expected:  []int{},
		},
		{
			name:      "predicate matches only last element",
			input:     []int{1, 2, 3, 4},
			predicate: func(x int) bool { return x == 4 },
			expected:  []int{1, 2, 3},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := DropLastWhile(testCase.input, testCase.predicate)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("DropLastWhile() = %v, expected %v", actual, testCase.expected)
			}
		})
	}
}

func TestDropWhile(t *testing.T) {
	testCases := []struct {
		name      string
		slice     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "drop elements less than 4",
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x < 4 },
			expected:  []int{4, 5},
		},
		{
			name:      "no elements match predicate",
			slice:     []int{1, 2, 3},
			predicate: func(x int) bool { return x > 10 },
			expected:  []int{1, 2, 3},
		},
		{
			name:      "all elements match predicate",
			slice:     []int{5, 6, 7},
			predicate: func(x int) bool { return x > 1 },
			expected:  []int{},
		},
		{
			name:      "empty slice",
			slice:     []int{},
			predicate: func(x int) bool { return x < 5 },
			expected:  []int{},
		},
		{
			name:      "predicate matches first element only",
			slice:     []int{2, 3, 4},
			predicate: func(x int) bool { return x == 2 },
			expected:  []int{3, 4},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := DropWhile(testCase.slice, testCase.predicate)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("DropWhile() = %v, expected %v", actual, testCase.expected)
			}
		})
	}
}

func TestFilterIndexed(t *testing.T) {
	testCases := []struct {
		name       string
		inputSlice []int
		predicate  func(int, int) bool
		expected   []int
	}{
		{
			name:       "index equals value",
			inputSlice: []int{0, 1, 2, 3, 4, 8, 6},
			predicate:  func(index, value int) bool { return index == value },
			expected:   []int{0, 1, 2, 3, 4, 6},
		},
		{
			name:       "index less than value",
			inputSlice: []int{1, 2, 3, 2, 1},
			predicate:  func(index, value int) bool { return index < value },
			expected:   []int{1, 2, 3},
		},
		{
			name:       "no elements match",
			inputSlice: []int{10, 10, 10},
			predicate:  func(index, value int) bool { return index > value },
			expected:   []int{},
		},
		{
			name:       "all elements match",
			inputSlice: []int{5, 6, 7},
			predicate:  func(index, value int) bool { return index < 10 },
			expected:   []int{5, 6, 7},
		},
		{
			name:       "empty slice",
			inputSlice: []int{},
			predicate:  func(index, value int) bool { return true },
			expected:   []int{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := FilterIndexed(testCase.inputSlice, testCase.predicate)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("FilterIndexed() = %v, expected %v", actual, testCase.expected)
			}
		})
	}
}

func TestFlatMap(t *testing.T) {
	testCases := []struct {
		name       string
		inputSlice []int
		transform  func(int) []int
		expected   []int
	}{
		{
			name:       "expand with squares",
			inputSlice: []int{1, 2, 3},
			transform:  func(x int) []int { return []int{x, x * x} },
			expected:   []int{1, 1, 2, 4, 3, 9},
		},
		{
			name:       "flatten empty slices",
			inputSlice: []int{1, 2, 3},
			transform:  func(x int) []int { return []int{} },
			expected:   []int{},
		},
		{
			name:       "identity",
			inputSlice: []int{1, 2, 3},
			transform:  func(x int) []int { return []int{x} },
			expected:   []int{1, 2, 3},
		},
		{
			name:       "empty input slice",
			inputSlice: []int{},
			transform:  func(x int) []int { return []int{x, x + 1} },
			expected:   []int{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := FlatMap(testCase.inputSlice, testCase.transform)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("FlatMap() = %v, expected %v", actual, testCase.expected)
			}
		})
	}
}

func TestFold(t *testing.T) {
	type args struct {
		s       []int
		initial int
		fn      func(int, int) int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"summation by fold",
			args{
				[]int{1, 2, 3, 4, 5},
				0,
				func(acc, v int) int { return acc + v },
			},
			15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Fold(tt.args.s, tt.args.initial, tt.args.fn); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("Fold() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestFoldIndexed(t *testing.T) {
	type args struct {
		s       []int
		initial int
		fn      func(int, int, int) int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"fold indexed",
			args{
				[]int{1, 2, 3, 4, 5},
				0,
				func(index, acc, v int) int { return acc + index*v },
			},
			40,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FoldIndexed(tt.args.s, tt.args.initial, tt.args.fn); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("FoldIndexed() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestGroupByWithOriginalTypesForKeyAndValue(t *testing.T) {

	input := []string{"a", "abc", "ab", "def", "abcd"}
	want := map[int][]string{
		1: {"a"},
		2: {"ab"},
		3: {"abc", "def"},
		4: {"abcd"},
	}
	got := GroupBy(input, func(str string) (int, string) {
		return len(str), str
	})
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GroupBy() = %v, expected %v", got, want)
	}
}

type wrapped struct {
	value string
}

func (w wrapped) String() string {
	return fmt.Sprintf("Wrapped:'%s'", w.value)
}

func TestGroupByWithNewTypesForKeyAndValue(t *testing.T) {

	input := []string{"a", "abc", "ab", "def", "abcd"}
	want := map[float64][]*wrapped{
		10.0: {&wrapped{"a"}},
		20.0: {&wrapped{"ab"}},
		30.0: {&wrapped{"abc"}, &wrapped{"def"}},
		40.0: {&wrapped{"abcd"}},
	}
	got := GroupBy(input, func(str string) (float64, *wrapped) {
		return float64(len(str)) * 10.0, &wrapped{str}
	})
	for k, vs := range want {
		avs, ok := got[k]
		if !ok {
			t.Errorf("expected key '%v' not found", k)
			return
		}
		if len(vs) != len(avs) {
			t.Errorf("expected %d elements for key:'%v'. got %d",
				len(vs), k, len(avs))
			return
		}
		for i := 0; i < len(vs); i++ {
			av := avs[i]
			v := vs[i]
			if av.value != v.value {
				t.Errorf("expected value: %s, got %s", v, av)
			}
		}

	}

}

func TestItems(t *testing.T) {
	m := map[string][]int{
		"a": {1, 2, 3, 4},
		"b": {1, 2},
		"c": {1, 2, 3},
	}
	want := []*Pair[string, []int]{
		{"a", []int{1, 2, 3, 4}},
		{"b", []int{1, 2}},
		{"c", []int{1, 2, 3}},
	}
	sort.Slice(want, func(i, j int) bool {
		return want[i].Fst < want[j].Fst
	})

	got := Items(m)
	sort.Slice(got, func(i, j int) bool {
		return got[i].Fst < got[j].Fst
	})

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Items() = %v, expected %v", got, want)
	}
}

func TestMap(t *testing.T) {
	type args struct {
		elems []int
		fn    func(int) int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"test mapping",
			args{
				[]int{1, 2, 3, 4, 5},
				func(i int) int { return i * i },
			},
			[]int{1, 4, 9, 16, 25},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.elems, tt.args.fn); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("Map() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestFlatMapIndexed(t *testing.T) {
	type args struct {
		elems []int
		fn    func(idx, val int) []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"test mapping",
			args{
				[]int{1, 2, 3, 4, 5},
				func(idx, val int) []int { return []int{idx, val * val} },
			},
			[]int{0, 1, 1, 4, 2, 9, 3, 16, 4, 25},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlatMapIndexed(tt.args.elems, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlatMap() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestMapIndexed(t *testing.T) {
	type args struct {
		s  []int
		fn func(int, int) int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"map indexed",
			args{
				[]int{1, 2, 3, 4, 5},
				func(index, i int) int { return index * i },
			},
			[]int{0, 2, 6, 12, 20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapIndexed(tt.args.s, tt.args.fn); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("MapIndexed() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestPartition(t *testing.T) {
	tests := []struct {
		name          string
		input         []int
		predicate     func(int) bool
		expectedTrue  []int
		expectedFalse []int
	}{
		{
			name:          "Even and odd numbers",
			input:         []int{1, 2, 3, 4, 5, 6},
			predicate:     func(n int) bool { return n%2 == 0 },
			expectedTrue:  []int{2, 4, 6},
			expectedFalse: []int{1, 3, 5},
		},
		{
			name:          "Greater than 3",
			input:         []int{1, 2, 3, 4, 5},
			predicate:     func(n int) bool { return n > 3 },
			expectedTrue:  []int{4, 5},
			expectedFalse: []int{1, 2, 3},
		},
		{
			name:          "Empty slice",
			input:         []int{},
			predicate:     func(n int) bool { return n > 0 },
			expectedTrue:  []int{},
			expectedFalse: []int{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			trueResult, falseResult := Partition(test.input, test.predicate)
			if !reflect.DeepEqual(trueResult, test.expectedTrue) {
				t.Errorf("expected trueList %v, got %v", test.expectedTrue, trueResult)
			}
			if !reflect.DeepEqual(falseResult, test.expectedFalse) {
				t.Errorf("expected falseList %v, got %v", test.expectedFalse, falseResult)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	type args struct {
		s  []int
		fn func(int, int) int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"reduce",
			args{
				[]int{1, 2, 3, 4, 5},
				func(acc, v int) int { return acc + v },
			},
			15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.s, tt.args.fn); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("Reduce() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestReduceIndexed(t *testing.T) {
	type args struct {
		s  []string
		fn func(int, string, string) string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"reduce indexed",
			args{
				[]string{"a", "b", "c", "d"},
				func(index int, acc, v string) string {
					return fmt.Sprintf("%s%s%d", acc, v, index)
				},
			},

			"ab1c2d3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReduceIndexed(tt.args.s, tt.args.fn); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("ReduceIndexed() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	Reverse(s)
	want := []int{7, 6, 5, 4, 3, 2, 1}
	if !reflect.DeepEqual(s, want) {
		t.Errorf("Reversed() = %v, expected %v", s, want)
	}
}

func TestReversed(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"reversed",
			args{
				[]int{1, 2, 3, 4, 5},
			},
			[]int{5, 4, 3, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reversed(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reversed() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestTake(t *testing.T) {
	type args struct {
		elems []int
		n     int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"take less than slice length",
			args{
				[]int{1, 2, 3, 4, 5, 6, 7, 8},
				4,
			},
			[]int{1, 2, 3, 4},
		},
		{"take more than slice length",
			args{
				[]int{1, 2, 3, 4, 5, 6, 7, 8},
				10,
			},
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{"take slice length",
			args{
				[]int{1, 2, 3, 4, 5, 6, 7, 8},
				8,
			},
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{"empty list",
			args{
				[]int{},
				8,
			},
			[]int{},
		},
		{"take 0 elements",
			args{
				[]int{1, 2, 3, 4, 5, 6, 7, 8},
				0,
			},
			[]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Take(tt.args.elems, tt.args.n); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("Take() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestTakeLast(t *testing.T) {
	type args struct {
		s []int
		n int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"less than slice len",
			args{
				[]int{1, 2, 3, 4, 5, 6, 7, 8},
				2,
			},
			[]int{7, 8},
		},
		{"more than slice len",
			args{
				[]int{1, 2, 3, 4, 5, 6, 7, 8},
				9,
			},
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{"exactly slice len",
			args{
				[]int{1, 2, 3, 4, 5, 6, 7, 8},
				8,
			},
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{"empty list",
			args{
				[]int{},
				8,
			},
			[]int{},
		},
		{"take 0 elements",
			args{
				[]int{1, 2, 3, 4, 5, 6, 7, 8},
				0,
			},
			[]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TakeLast(tt.args.s, tt.args.n); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("TakeLast() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestTake2(t *testing.T) {
	got := Take(alphabet(), 2)
	expected := []rune{'a', 'b'}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("TakeLast() = %v expected %v", got, expected)
	}
}

func TestTakeLast2(t *testing.T) {
	got := TakeLast(alphabet(), 2)
	expected := []rune{'y', 'z'}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("TakeLast() = %v expected %v", got, expected)
	}
}

func TestTakeLastWhile(t *testing.T) {
	got := TakeLastWhile(alphabet(), func(s rune) bool { return s > 'w' })
	expected := []rune{'x', 'y', 'z'}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("TakeLastWhile() = %v expected %v", got, expected)
	}
}

func TestTakeWhile(t *testing.T) {
	got := TakeWhile(alphabet(), func(s rune) bool { return s < 'f' })
	expected := []rune{'a', 'b', 'c', 'd', 'e'}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("TakeWhile() = %v expected %v", got, expected)
	}
}

func TestWindowed(t *testing.T) {
	type args struct {
		s    []int
		size int
		step int
	}
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{"size = 5, step = 1",
			args{input, 5, 1},
			[][]int{
				{1, 2, 3, 4, 5},
				{2, 3, 4, 5, 6},
				{3, 4, 5, 6, 7},
				{4, 5, 6, 7, 8},
				{5, 6, 7, 8, 9},
				{6, 7, 8, 9, 10},
				{7, 8, 9, 10},
				{8, 9, 10},
				{9, 10},
				{10},
			},
		},
		{"size = 5, step = 3",
			args{input, 5, 3},
			[][]int{
				{1, 2, 3, 4, 5},
				{4, 5, 6, 7, 8},
				{7, 8, 9, 10},
				{10},
			},
		},
		{"size = 3, step = 4",
			args{input, 3, 4},
			[][]int{
				{1, 2, 3},
				{5, 6, 7},
				{9, 10},
			},
		},

		{"slice smaller than size",
			args{[]int{1, 2, 3}, 4, 1},
			[][]int{
				{1, 2, 3},
				{2, 3},
				{3},
			},
		},
		{"slice smaller than size and step",
			args{[]int{1, 2, 3}, 4, 4},
			[][]int{
				{1, 2, 3},
			},
		},
		{"slice larger than size and smaller than step",
			args{[]int{1, 2, 3}, 2, 4},
			[][]int{
				{1, 2},
			},
		},
		{"empty slice",
			args{[]int{}, 4, 4},
			[][]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Windowed(tt.args.s, tt.args.size, tt.args.step); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("Windowed() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestZip(t *testing.T) {
	s1 := []string{"a", "b", "c", "d"}
	s2 := []int{1, 2, 3}
	got := Zip(s1, s2)
	want := []*Pair[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Zip() = %v, expected %v", got, want)
	}
}

func TestUnZip(t *testing.T) {
	ps := []*Pair[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}
	want1 := []string{"a", "b", "c"}
	want2 := []int{1, 2, 3}
	got1, got2 := Unzip(ps)
	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("Zip() first list = %v, expected %v", got1, want1)
	}
	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("Zip() first list = %v, expected %v", got2, want2)
	}
}

func alphabet() []rune {
	ret := make([]rune, 0)
	for r := 'a'; r <= 'z'; r++ {
		ret = append(ret, r)
	}
	return ret
}

func TestGetOrInsert(t *testing.T) {

	m := map[int]int{1: 10, 2: 20}
	fn := func(k int) int { return k * 10 }

	// case where key is present
	expected := 20
	got := GetOrInsert(m, 2, fn)
	if got != expected {
		t.Errorf("GetOrInsert() = %d expected = %d", got, expected)
	}

	// case where key is not present, but populated by invoking the function
	expected = 30
	got = GetOrInsert(m, 3, fn)
	if got != expected {
		t.Errorf("GetOrInsert() = %d expected = %d", got, expected)
	}

	// check that the new value was stored in the map as well
	expected = 30
	got, ok := m[3]
	if !ok {
		t.Errorf("GetOrInsert did not insert new value in map!!")
	}
	if got != expected {
		t.Errorf("value in map = %d expected = %d", got, expected)
	}

}

func TestFoldItems(t *testing.T) {
	type args struct {
		m       map[int]int
		initial map[string]string
		fn      func(map[string]string, int, int) map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			"test fold over map items",
			args{
				map[int]int{1: 10, 2: 20, 3: 30},
				make(map[string]string),
				func(acc map[string]string, k, v int) map[string]string {
					acc[fmt.Sprintf("entry_%d", k)] = fmt.Sprintf(
						"%d->%d",
						k,
						v,
					)
					return acc
				},
			},
			map[string]string{
				"entry_1": "1->10",
				"entry_2": "2->20",
				"entry_3": "3->30",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FoldItems(tt.args.m, tt.args.initial, tt.args.fn); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("FoldItems() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestTransformMap(t *testing.T) {
	type args struct {
		m  map[string][]int
		fn func(k string, v []int) (string, []int, bool)
	}
	tests := []struct {
		name string
		args args
		want map[string][]int
	}{
		{
			"filter entries",
			args{
				map[string][]int{
					"a": {1, 2, 3, 4},
					"b": {1, 2},
					"c": {1, 2, 3},
				},
				func(k string, v []int) (string, []int, bool) {
					if len(v) < 3 {
						return k, v, false
					}
					return k, v, true
				},
			},
			map[string][]int{
				"a": {1, 2, 3, 4},
				"c": {1, 2, 3},
			},
		},
		{
			"map entries",
			args{
				map[string][]int{
					"a": {1, 2, 3, 4},
					"b": {5, 6},
				},
				func(k string, v []int) (string, []int, bool) {
					newK := strings.ToUpper(k)
					newV := Map(v, func(i int) int { return i * 10 })
					return newK, newV, true
				},
			},
			map[string][]int{
				"A": {10, 20, 30, 40},
				"B": {50, 60},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TransformMap(tt.args.m, tt.args.fn); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("TransformMap() = %v, expected %v", got, tt.want)
			}
		})
	}
}
