package dirForStudy

import (
	"math"
	"reflect"
	"sort"
	"testing"
)

func TestSortSlices(t *testing.T) {
	tests := []struct {
		name  string
		input [][]int
		want  []int
	}{
		{
			name:  "no slices",
			input: [][]int{},
			want:  []int{},
		},
		{
			name:  "single slice",
			input: [][]int{{3, 1, 2}},
			want:  []int{1, 2, 3},
		},
		{
			name:  "multiple slices",
			input: [][]int{{3, 5}, {4, 2}, {1}},
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "with duplicates",
			input: [][]int{{2, 1, 2}, {3, 2}},
			want:  []int{1, 2, 2, 2, 3},
		},
		{
			name:  "negatives and positives",
			input: [][]int{{-1, 3}, {0, -2}},
			want:  []int{-2, -1, 0, 3},
		},
		{
			name:  "empty inner slices",
			input: [][]int{{}, {}},
			want:  []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sortSlices(tc.input...)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("sortSlices(%v) = %v; want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestFromMapToSlice(t *testing.T) {
	tests := []struct {
		name  string
		input [][]int
		want  []int
	}{
		{
			name:  "no slices",
			input: [][]int{},
			want:  []int{},
		},
		{
			name:  "single slice",
			input: [][]int{{3, 1, 2}},
			want:  []int{1, 2, 3},
		},
		{
			name:  "multiple slices",
			input: [][]int{{5, 3}, {4, 2}, {1}},
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "with duplicates across slices",
			input: [][]int{{2, 1, 2}, {3, 2, 1}},
			want:  []int{1, 2, 3},
		},
		{
			name:  "negatives and positives",
			input: [][]int{{-1, 3}, {0, -2, 3}},
			want:  []int{-2, -1, 0, 3},
		},
		{
			name:  "empty inner slices",
			input: [][]int{{}, {}},
			want:  []int{},
		},
		{
			name:  "unsorted input preserved uniqueness",
			input: [][]int{{10, 5, 7}, {7, 5, 10}, {8, 6}},
			want:  []int{5, 6, 7, 8, 10},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sortFromMapToSlice(tc.input...)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("sortFromMapToSlice(%v) = %v; want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestAddSumKey(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]struct{}
		want  map[string]struct{}
	}{
		{
			name:  "empty map",
			input: map[string]struct{}{},
			want:  map[string]struct{}{},
		},
		{
			name:  "single key",
			input: map[string]struct{}{"x": {}},
			want:  map[string]struct{}{"x": {}},
		},
		{
			name:  "two keys",
			input: map[string]struct{}{"b": {}, "a": {}},
			want: map[string]struct{}{
				"a":  {},
				"b":  {},
				"ba": {},
			},
		},
		{
			name:  "composite already exists",
			input: map[string]struct{}{"a": {}, "b": {}, "ba": {}},
			want: map[string]struct{}{
				"a":  {},
				"b":  {},
				"ba": {},
			},
		},
		{
			name:  "multiple keys",
			input: map[string]struct{}{"d": {}, "b": {}, "a": {}},
			want: map[string]struct{}{
				"a":  {},
				"b":  {},
				"d":  {},
				"ba": {},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := addSumKey(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("addSumKey() = %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want []int
	}{
		{
			name: "both empty",
			a:    []int{},
			b:    []int{},
			want: []int{},
		},
		{
			name: "a non-empty, b empty",
			a:    []int{1, 2, 3},
			b:    []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "a empty, b non-empty",
			a:    []int{},
			b:    []int{1, 2, 3},
			want: []int{},
		},
		{
			name: "no common elements",
			a:    []int{1, 3, 5},
			b:    []int{2, 4, 6},
			want: []int{1, 3, 5},
		},
		{
			name: "some common, with duplicates in a",
			a:    []int{1, 2, 3, 2},
			b:    []int{2, 4},
			want: []int{1, 3},
		},
		{
			name: "duplicates only in a",
			a:    []int{1, 1, 2},
			b:    []int{2},
			want: []int{1, 1},
		},
		{
			name: "b covers all of a",
			a:    []int{1, 2},
			b:    []int{1, 2, 3},
			want: []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := subtractSlices(tc.a, tc.b)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("substactSlices(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestIntersectionSlices(t *testing.T) {
	tests := []struct {
		name string
		arr1 []int
		arr2 []int
		want []int
	}{
		{
			name: "both empty",
			arr1: []int{},
			arr2: []int{},
			want: []int{},
		},
		{
			name: "first empty",
			arr1: []int{},
			arr2: []int{1, 2, 3},
			want: []int{},
		},
		{
			name: "second empty",
			arr1: []int{1, 2, 3},
			arr2: []int{},
			want: []int{},
		},
		{
			name: "no intersection",
			arr1: []int{1, 2, 3},
			arr2: []int{4, 5, 6},
			want: []int{},
		},
		{
			name: "simple intersection",
			arr1: []int{1, 2, 3},
			arr2: []int{2, 3, 4},
			want: []int{2, 3},
		},
		{
			name: "with duplicates in both",
			arr1: []int{1, 2, 2, 3},
			arr2: []int{2, 2, 4},
			want: []int{2, 2},
		},
		{
			name: "order preserved from first",
			arr1: []int{5, 1, 2, 3},
			arr2: []int{1, 3, 5},
			want: []int{5, 1, 3},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := intersectionSlices(tc.arr1, tc.arr2)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("intersectionSlices(%v, %v) = %v; want %v",
					tc.arr1, tc.arr2, got, tc.want)
			}
		})
	}
}

func TestMirrorMap(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]int
		want  map[int][]string
	}{
		{
			name:  "empty map",
			input: map[string]int{},
			want:  map[int][]string{},
		},
		{
			name:  "single pair",
			input: map[string]int{"apple": 1},
			want:  map[int][]string{1: {"apple"}},
		},
		{
			name:  "multiple distinct values",
			input: map[string]int{"a": 1, "b": 2, "c": 3},
			want: map[int][]string{
				1: {"a"},
				2: {"b"},
				3: {"c"},
			},
		},
		{
			name:  "grouping keys by same value",
			input: map[string]int{"red": 10, "blue": 20, "green": 10, "yellow": 20},
			want: map[int][]string{
				10: {"green", "red"},
				20: {"blue", "yellow"},
			},
		},
		{
			name:  "multiple groups with single and multiple",
			input: map[string]int{"x": 0, "y": 1, "z": 0, "w": 2},
			want: map[int][]string{
				0: {"x", "z"},
				1: {"y"},
				2: {"w"},
			},
		},
		{
			name:  "keys with same prefix",
			input: map[string]int{"aa": 5, "ab": 5, "ac": 6, "ad": 5},
			want: map[int][]string{
				5: {"aa", "ab", "ad"},
				6: {"ac"},
			},
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			got := mirrorMap(tc.input)

			// Sort slices in both got and want before comparing
			for k := range got {
				sort.Strings(got[k])
			}
			for k := range tc.want {
				sort.Strings(tc.want[k])
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("groupKeysByValue(%v) = %v; want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestMinAndMax(t *testing.T) {
	tests := []struct {
		name    string
		arr     []float64
		wantMin float64
		wantMax float64
	}{
		{
			name:    "ascending",
			arr:     []float64{1.1, 2.2, 3.3, 4.4},
			wantMin: 1.1,
			wantMax: 4.4,
		},
		{
			name:    "descending",
			arr:     []float64{5.5, 4.4, 3.3, 2.2},
			wantMin: 2.2,
			wantMax: 5.5,
		},
		{
			name:    "all equal",
			arr:     []float64{7.7, 7.7, 7.7},
			wantMin: 7.7,
			wantMax: 7.7,
		},
		{
			name:    "mixed positive and negative",
			arr:     []float64{-3.5, 0.0, 2.5, -1.2},
			wantMin: -3.5,
			wantMax: 2.5,
		},
		{
			name:    "single element",
			arr:     []float64{9.9},
			wantMin: 9.9,
			wantMax: 9.9,
		},
		{
			name:    "with zeros",
			arr:     []float64{0.0, 0.0, 0.0},
			wantMin: 0.0,
			wantMax: 0.0,
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			gotMin, gotMax := minAndMax(tc.arr)
			if gotMin != tc.wantMin || gotMax != tc.wantMax {
				t.Errorf("minAndMax(%v) = (%v, %v); want (%v, %v)",
					tc.arr, gotMin, gotMax, tc.wantMin, tc.wantMax)
			}
		})
	}
}

func TestFilterSlice(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		filter func(int) bool
		want   []int
	}{
		{
			name:   "empty slice",
			arr:    []int{},
			filter: func(int) bool { return true },
			want:   []int{},
		},
		{
			name:   "no removal",
			arr:    []int{1, 2, 3, 4},
			filter: func(int) bool { return true },
			want:   []int{1, 2, 3, 4},
		},
		{
			name:   "remove all",
			arr:    []int{1, 2, 3, 4},
			filter: func(int) bool { return false },
			want:   []int{},
		},
		{
			name:   "remove evens",
			arr:    []int{1, 2, 3, 4, 5},
			filter: func(n int) bool { return n%2 != 0 },
			want:   []int{1, 3, 5},
		},
		{
			name:   "remove odds",
			arr:    []int{1, 2, 3, 4, 5},
			filter: func(n int) bool { return n%2 == 0 },
			want:   []int{2, 4},
		},
		{
			name: "remove non-primes",
			arr:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			filter: func(n int) bool {
				if n < 2 {
					return false
				}
				for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
					if n%i == 0 {
						return false
					}
				}
				return true
			},
			want: []int{2, 3, 5, 7},
		},
		{
			name:   "remove > threshold",
			arr:    []int{1, 5, 10, 15, 20},
			filter: func(n int) bool { return n <= 10 },
			want:   []int{1, 5, 10},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := filterSlice(append([]int(nil), tc.arr...), tc.filter)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("filterSlice(%v) = %v; want %v", tc.arr, got, tc.want)
			}
		})
	}
}

func TestSplitMaps(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		num  int
		want [][]int
	}{
		{
			name: "empty slice",
			arr:  []int{},
			num:  3,
			want: [][]int{},
		},
		{
			name: "chunk size greater than length",
			arr:  []int{1, 2, 3},
			num:  5,
			want: [][]int{{1, 2, 3}},
		},
		{
			name: "exact division",
			arr:  []int{1, 2, 3, 4},
			num:  2,
			want: [][]int{{1, 2}, {3, 4}},
		},
		{
			name: "non-exact division",
			arr:  []int{1, 2, 3, 4, 5},
			num:  2,
			want: [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name: "chunk size of one",
			arr:  []int{7, 8, 9},
			num:  1,
			want: [][]int{{7}, {8}, {9}},
		},
		{
			name: "chunk size equal length",
			arr:  []int{10, 11, 12},
			num:  3,
			want: [][]int{{10, 11, 12}},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := splitMaps(tc.arr, tc.num)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("splitMaps(%v, %d) = %v; want %v",
					tc.arr, tc.num, got, tc.want)
			}
		})
	}
}

func TestCombineMaps(t *testing.T) {
	tests := []struct {
		name string
		m1   map[string]int
		m2   map[string]int
		want map[string]int
	}{
		{
			name: "both empty",
			m1:   map[string]int{},
			m2:   map[string]int{},
			want: map[string]int{},
		},
		{
			name: "first empty",
			m1:   map[string]int{},
			m2:   map[string]int{"a": 1, "b": 2},
			want: map[string]int{"a": 1, "b": 2},
		},
		{
			name: "second empty",
			m1:   map[string]int{"x": 5, "y": 10},
			m2:   map[string]int{},
			want: map[string]int{"x": 5, "y": 10},
		},
		{
			name: "no overlapping keys",
			m1:   map[string]int{"a": 1, "b": 2},
			m2:   map[string]int{"c": 3, "d": 4},
			want: map[string]int{"a": 1, "b": 2, "c": 3, "d": 4},
		},
		{
			name: "overlap m2 has higher values",
			m1:   map[string]int{"a": 1, "b": 2},
			m2:   map[string]int{"b": 5, "c": 3},
			want: map[string]int{"a": 1, "b": 5, "c": 3},
		},
		{
			name: "overlap m1 has higher values",
			m1:   map[string]int{"a": 10, "b": 2},
			m2:   map[string]int{"a": 5, "c": 3},
			want: map[string]int{"a": 10, "b": 2, "c": 3},
		},
		{
			name: "mixed overlaps",
			m1:   map[string]int{"a": 1, "b": 8, "c": 3},
			m2:   map[string]int{"a": 4, "b": 2, "d": 6},
			want: map[string]int{"a": 4, "b": 8, "c": 3, "d": 6},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := combineMaps(tc.m1, tc.m2)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("combineMaps(%v, %v) = %v; want %v",
					tc.m1, tc.m2, got, tc.want)
			}
		})
	}
}

func TestGroupByStruct(t *testing.T) {
	tests := []struct {
		name  string
		input []Item
		want  map[string][]int
	}{
		{
			name:  "empty slice",
			input: []Item{},
			want:  map[string][]int{},
		},
		{
			name:  "single item",
			input: []Item{{Category: "A", Value: 1}},
			want:  map[string][]int{"A": {1}},
		},
		{
			name: "multiple categories",
			input: []Item{
				{Category: "A", Value: 1},
				{Category: "B", Value: 2},
				{Category: "A", Value: 3},
				{Category: "C", Value: 4},
			},
			want: map[string][]int{
				"A": {1, 3},
				"B": {2},
				"C": {4},
			},
		},
		{
			name: "all same category",
			input: []Item{
				{Category: "X", Value: 5},
				{Category: "X", Value: 6},
				{Category: "X", Value: 7},
			},
			want: map[string][]int{"X": {5, 6, 7}},
		},
		{
			name: "interleaved categories",
			input: []Item{
				{Category: "A", Value: 1},
				{Category: "B", Value: 2},
				{Category: "A", Value: 3},
				{Category: "B", Value: 4},
			},
			want: map[string][]int{
				"A": {1, 3},
				"B": {2, 4},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := groupByStruct(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("groupByStruct(%v) = %v; want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestDeleteDuplicates(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "empty slice",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "no duplicates",
			input: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "all duplicates",
			input: []string{"x", "x", "x", "x"},
			want:  []string{"x"},
		},
		{
			name:  "mixed duplicates",
			input: []string{"a", "b", "a", "c", "b", "d"},
			want:  []string{"a", "b", "c", "d"},
		},
		{
			name:  "adjacent duplicates",
			input: []string{"foo", "foo", "bar", "bar", "baz", "baz"},
			want:  []string{"foo", "bar", "baz"},
		},
		{
			name:  "non-adjacent duplicates preserve order",
			input: []string{"1", "2", "1", "3", "2", "4"},
			want:  []string{"1", "2", "3", "4"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := deleteDuplicates(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("deleteDups(%v) = %v; want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestFindBinary(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		targets []int
		want    []int
	}{
		{
			name:    "no targets",
			input:   []int{1, 2, 3},
			targets: []int{},
			want:    []int{},
		},
		{
			name:    "all targets found",
			input:   []int{1, 2, 3, 4, 5},
			targets: []int{1, 3, 5},
			want:    []int{1, 3, 5},
		},
		{
			name:    "none found",
			input:   []int{10, 20, 30},
			targets: []int{1, 2, 3},
			want:    []int{},
		},
		{
			name:    "some found, some not",
			input:   []int{2, 4, 6, 8},
			targets: []int{1, 2, 3, 4, 8},
			want:    []int{2, 4, 8},
		},
		{
			name:    "duplicates in targets",
			input:   []int{1, 3, 5, 7},
			targets: []int{3, 3, 5, 5, 9},
			want:    []int{3, 3, 5, 5},
		},
		{
			name:    "unsorted targets preserved order",
			input:   []int{1, 2, 3, 4, 5},
			targets: []int{5, 1, 4},
			want:    []int{5, 1, 4},
		},
		{
			name:    "single element slice",
			input:   []int{42},
			targets: []int{42, 7},
			want:    []int{42},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := findBinary(tc.input, tc.targets)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("findBinary(%v, %v) = %v; want %v",
					tc.input, tc.targets, got, tc.want)
			}
		})
	}
}
