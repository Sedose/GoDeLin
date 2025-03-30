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
		inputSlice []int
		predicate  func(int) bool
	}
	testCases := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			"positive case",
			args{
				[]int{1, 2, 3, 4, 5},
				func(i int) bool { return i < 10 },
			},
			true,
		},
		{
			"negative case",
			args{
				[]int{1, 2, 3, 4, 5},
				func(i int) bool { return i%2 == 0 },
			},
			false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := All(testCase.args.inputSlice, testCase.args.predicate); actual != testCase.expected {
				t.Errorf("All() = %v, expected %v", actual, testCase.expected)
			}
		})
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
		{
			"positive case",
			args{
				[]int{1, 2, 3, 4, 5, 6},
				func(i int) bool { return i%2 == 0 },
			},
			true,
		},
		{
			"negative case",
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

func TestGetOrPut(t *testing.T) {
	testCases := []struct {
		name         string
		initialMap   map[string]int
		key          string
		defaultValue func(string) int
		expected     int
		finalMap     map[string]int
	}{
		{
			name:       "key exists in map",
			initialMap: map[string]int{"a": 10, "b": 20},
			key:        "a",
			defaultValue: func(k string) int {
				return 99
			},
			expected: 10,
			finalMap: map[string]int{"a": 10, "b": 20},
		},
		{
			name:       "key does not exist, default applied",
			initialMap: map[string]int{"x": 1},
			key:        "y",
			defaultValue: func(k string) int {
				return 42
			},
			expected: 42,
			finalMap: map[string]int{"x": 1, "y": 42},
		},
		{
			name:       "default function uses key",
			initialMap: map[string]int{},
			key:        "key42",
			defaultValue: func(k string) int {
				if k == "key42" {
					return 123
				}
				return 0
			},
			expected: 123,
			finalMap: map[string]int{"key42": 123},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := GetOrPut(testCase.initialMap, testCase.key, testCase.defaultValue)
			if result != testCase.expected {
				t.Errorf("GetOrPut() = %v, expected %v", result, testCase.expected)
			}
			if !reflect.DeepEqual(testCase.initialMap, testCase.finalMap) {
				t.Errorf("map after GetOrPut = %v, expected %v", testCase.initialMap, testCase.finalMap)
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
	var empty []string
	result := GroupBy(empty, func(s string) (string, string) {
		return s, s
	})
	expected := map[string][]string{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Associate() with empty slice = %v, expected %v", result, expected)
	}
}

func TestChunkedBy(t *testing.T) {
	input := []int{
		10, 20, 30, 40,
		31, 31, 33, 34,
		21, 22, 23, 24,
		11, 12, 13, 14,
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
		{
			"test distinct",
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
	}{
		{
			"test distinctBy",
			args{
				[]string{"a", "A", "b", "B", "c", "C"},
				func(s string) string { return strings.ToLower(s) },
			},
			[]string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DistinctBy(tt.args.s, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DistinctBy() = %v, expected %v", got, tt.want)
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
			t.Errorf("expected %d elements for key:'%v'. got %d", len(vs), k, len(avs))
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
	// Updated from []*Pair[string, []int] to []Pair[string, []int]
	want := []Pair[string, []int]{
		{"a", []int{1, 2, 3, 4}},
		{"b", []int{1, 2}},
		{"c", []int{1, 2, 3}},
	}
	sort.Slice(want, func(i, j int) bool {
		return want[i].First < want[j].First
	})

	got := Items(m) // now returns []Pair[K, V]
	sort.Slice(got, func(i, j int) bool {
		return got[i].First < got[j].First
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
		{
			"test mapping",
			args{
				[]int{1, 2, 3, 4, 5},
				func(i int) int { return i * i },
			},
			[]int{1, 4, 9, 16, 25},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.elems, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestFlatMapIndexed(t *testing.T) {
	testCases := []struct {
		name       string
		inputSlice []int
		transform  func(int, int) []int
		expected   []int
	}{
		{
			name:       "index and square",
			inputSlice: []int{1, 2, 3},
			transform:  func(i, v int) []int { return []int{i, v * v} },
			expected:   []int{0, 1, 1, 4, 2, 9},
		},
		{
			name:       "repeat element by index",
			inputSlice: []int{7, 8, 9},
			transform: func(i, v int) []int {
				repeated := make([]int, i)
				for j := 0; j < i; j++ {
					repeated[j] = v
				}
				return repeated
			},
			expected: []int{8, 9, 9},
		},
		{
			name:       "empty input slice",
			inputSlice: []int{},
			transform:  func(i, v int) []int { return []int{i + v} },
			expected:   []int{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := FlatMapIndexed(testCase.inputSlice, testCase.transform)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("FlatMapIndexed() = %v, expected %v", actual, testCase.expected)
			}
		})
	}
}

func TestFold(t *testing.T) {
	testCases := []struct {
		name       string
		inputSlice []int
		initial    int
		operation  func(int, int) int
		expected   int
	}{
		{
			name:       "sum of elements",
			inputSlice: []int{1, 2, 3, 4},
			initial:    0,
			operation:  func(acc, v int) int { return acc + v },
			expected:   10,
		},
		{
			name:       "product of elements",
			inputSlice: []int{1, 2, 3, 4},
			initial:    1,
			operation:  func(acc, v int) int { return acc * v },
			expected:   24,
		},
		{
			name:       "subtract all elements",
			inputSlice: []int{1, 2, 3},
			initial:    10,
			operation:  func(acc, v int) int { return acc - v },
			expected:   4, // 10 - 1 - 2 - 3
		},
		{
			name:       "empty input slice",
			inputSlice: []int{},
			initial:    42,
			operation:  func(acc, v int) int { return acc + v },
			expected:   42,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := Fold(testCase.inputSlice, testCase.initial, testCase.operation)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Fold() = %v, expected %v", actual, testCase.expected)
			}
		})
	}
}

func TestFoldIndexed(t *testing.T) {
	testCases := []struct {
		name       string
		inputSlice []int
		initial    int
		operation  func(int, int, int) int
		expected   int
	}{
		{
			name:       "weighted sum by index",
			inputSlice: []int{1, 2, 3, 4},
			initial:    0,
			operation:  func(index, acc, v int) int { return acc + index*v },
			expected:   20, // 0+1*2+2*3+3*4 = 2+6+12
		},
		{
			name:       "subtract with index multiplier",
			inputSlice: []int{5, 6, 7},
			initial:    100,
			operation:  func(index, acc, v int) int { return acc - index*v },
			expected:   100 - 1*6 - 2*7, // 100 - 6 - 14 = 80
		},
		{
			name:       "empty slice",
			inputSlice: []int{},
			initial:    123,
			operation:  func(i, acc, v int) int { return acc + v + i },
			expected:   123,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := FoldIndexed(testCase.inputSlice, testCase.initial, testCase.operation)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("FoldIndexed() = %v, expected %v", actual, testCase.expected)
			}
		})
	}
}

func TestFoldMapEntries(t *testing.T) {
	testCases := []struct {
		name     string
		inputMap map[string]int
		initial  int
		combine  func(int, string, int) int
		expected int
	}{
		{
			name:     "sum values",
			inputMap: map[string]int{"a": 10, "b": 20, "c": 30},
			initial:  0,
			combine:  func(acc int, _ string, v int) int { return acc + v },
			expected: 60,
		},
		{
			name:     "sum keys length and values",
			inputMap: map[string]int{"apple": 5, "banana": 6},
			initial:  0,
			combine:  func(acc int, k string, v int) int { return acc + len(k) + v },
			expected: (5 + 5) + (6 + 6), // 10 + 12 = 22
		},
		{
			name:     "product of key length and value",
			inputMap: map[string]int{"a": 3, "go": 5},
			initial:  1,
			combine: func(acc int, k string, v int) int {
				return acc * len(k) * v
			},
			expected: 1 * 1 * 3 * 2 * 5, // 3 * 10 = 30
		},
		{
			name:     "empty map",
			inputMap: map[string]int{},
			initial:  100,
			combine:  func(acc int, _ string, _ int) int { return acc + 1 },
			expected: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FoldMapEntries(tc.inputMap, tc.initial, tc.combine)
			if result != tc.expected {
				t.Errorf("FoldMapEntries() = %v, expected %v", result, tc.expected)
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
		{
			"map indexed",
			args{
				[]int{1, 2, 3, 4, 5},
				func(index, i int) int { return index * i },
			},
			[]int{0, 2, 6, 12, 20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapIndexed(tt.args.s, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
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
		{
			"reduce",
			args{
				[]int{1, 2, 3, 4, 5},
				func(acc, v int) int { return acc + v },
			},
			15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.s, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
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
		{
			"reduce indexed",
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
			if got := ReduceIndexed(tt.args.s, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReduceIndexed() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestTakeWhile(t *testing.T) {
	got := TakeWhile(alphabet(), func(s rune) bool { return s < 'f' })
	expected := []rune{'a', 'b', 'c', 'd', 'e'}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("TakeWhile() = %v expected %v", got, expected)
	}
}

func TestMapEntries(t *testing.T) {
	type args struct {
		inputMap  map[string]int
		transform func(string, int) (string, int)
	}
	testCases := []struct {
		name     string
		args     args
		expected map[string]int
	}{
		{
			name: "rename keys and double values",
			args: args{
				inputMap: map[string]int{"a": 1, "b": 2},
				transform: func(k string, v int) (string, int) {
					return k + "_x", v * 2
				},
			},
			expected: map[string]int{"a_x": 2, "b_x": 4},
		},
		{
			name: "identity transformation",
			args: args{
				inputMap: map[string]int{"m": 3, "n": 4},
				transform: func(k string, v int) (string, int) {
					return k, v
				},
			},
			expected: map[string]int{"m": 3, "n": 4},
		},
		{
			name: "key squash",
			args: args{
				inputMap: map[string]int{"x": 10, "y": 20},
				transform: func(k string, v int) (string, int) {
					return "z", v
				},
			},
			expected: map[string]int{"z": 20}, // last one wins
		},
		{
			name: "empty input map",
			args: args{
				inputMap: map[string]int{},
				transform: func(k string, v int) (string, int) {
					return k, v
				},
			},
			expected: map[string]int{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := MapEntries(testCase.args.inputMap, testCase.args.transform)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("MapEntries() = %v, expected %v", actual, testCase.expected)
			}
		})
	}
}

func TestUnzip(t *testing.T) {
	type args struct {
		pairs []Pair[string, int]
	}
	testCases := []struct {
		name       string
		args       args
		wantFirst  []string
		wantSecond []int
	}{
		{
			name: "multiple pairs",
			args: args{
				pairs: []Pair[string, int]{
					{"a", 10},
					{"b", 20},
					{"c", 30},
				},
			},
			wantFirst:  []string{"a", "b", "c"},
			wantSecond: []int{10, 20, 30},
		},
		{
			name:       "empty slice",
			args:       args{pairs: []Pair[string, int]{}},
			wantFirst:  []string{},
			wantSecond: []int{},
		},
		{
			name: "single pair",
			args: args{
				pairs: []Pair[string, int]{
					{"x", 99},
				},
			},
			wantFirst:  []string{"x"},
			wantSecond: []int{99},
		},
		{
			name: "repeated values",
			args: args{
				pairs: []Pair[string, int]{
					{"z", 1},
					{"z", 1},
				},
			},
			wantFirst:  []string{"z", "z"},
			wantSecond: []int{1, 1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gotFirst, gotSecond := Unzip(testCase.args.pairs)
			if !reflect.DeepEqual(gotFirst, testCase.wantFirst) {
				t.Errorf("Unzip() firsts = %v, expected %v", gotFirst, testCase.wantFirst)
			}
			if !reflect.DeepEqual(gotSecond, testCase.wantSecond) {
				t.Errorf("Unzip() seconds = %v, expected %v", gotSecond, testCase.wantSecond)
			}
		})
	}
}

func TestZip(t *testing.T) {
	type args struct {
		left  []string
		right []int
	}
	testCases := []struct {
		name string
		args args
		want []Pair[string, int]
	}{
		{
			name: "same length slices",
			args: args{
				[]string{"a", "b", "c"},
				[]int{10, 20, 30},
			},
			want: []Pair[string, int]{
				{"a", 10},
				{"b", 20},
				{"c", 30},
			},
		},
		{
			name: "left slice longer",
			args: args{
				[]string{"a", "b", "c", "d"},
				[]int{1, 2, 3},
			},
			want: []Pair[string, int]{
				{"a", 1},
				{"b", 2},
				{"c", 3},
			},
		},
		{
			name: "right slice longer",
			args: args{
				[]string{"x", "y"},
				[]int{7, 8, 9, 10},
			},
			want: []Pair[string, int]{
				{"x", 7},
				{"y", 8},
			},
		},
		{
			name: "left slice empty",
			args: args{
				[]string{},
				[]int{100},
			},
			want: []Pair[string, int]{},
		},
		{
			name: "right slice empty",
			args: args{
				[]string{"a", "b", "c"},
				[]int{},
			},
			want: []Pair[string, int]{},
		},
		{
			name: "both slices empty",
			args: args{
				[]string{},
				[]int{},
			},
			want: []Pair[string, int]{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := Zip(testCase.args.left, testCase.args.right)
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("Zip() = %v, expected %v", got, testCase.want)
			}
		})
	}
}

func TestWindowed(t *testing.T) {
	testCases := []struct {
		name        string
		input       []int
		size        int
		step        int
		expected    [][]int
		shouldPanic bool
	}{
		{
			name:  "empty slice",
			input: []int{},
			size:  1,
			step:  1,
			// now expect an empty slice rather than nil
			expected: [][]int{},
		},
		{
			name:        "size is zero causes panic",
			input:       []int{1, 2, 3},
			size:        0,
			step:        1,
			shouldPanic: true,
		},
		{
			name:        "step is zero causes panic",
			input:       []int{1, 2, 3},
			size:        1,
			step:        0,
			shouldPanic: true,
		},
		{
			name:     "basic window size=2 step=2",
			input:    []int{1, 2, 3, 4, 5},
			size:     2,
			step:     2,
			expected: [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:  "overlapping windows size=2 step=1",
			input: []int{1, 2, 3, 4},
			size:  2,
			step:  1,
			expected: [][]int{
				{1, 2},
				{2, 3},
				{3, 4},
				{4},
			},
		},
		{
			name:  "size=3 step=1 with partial last window",
			input: []int{1, 2, 3, 4, 5},
			size:  3,
			step:  1,
			expected: [][]int{
				{1, 2, 3},
				{2, 3, 4},
				{3, 4, 5},
				{4, 5},
				{5},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected a panic but did not get one")
					}
				}()
			}

			result := Windowed(tc.input, tc.size, tc.step)
			if !reflect.DeepEqual(result, tc.expected) && !tc.shouldPanic {
				t.Errorf("Windowed() = %v, expected %v", result, tc.expected)
			}
		})
	}
}

func alphabet() []rune {
	ret := make([]rune, 0)
	for r := 'a'; r <= 'z'; r++ {
		ret = append(ret, r)
	}
	return ret
}
