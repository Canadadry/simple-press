package paginator

import (
	"reflect"
	"testing"
)

func TestComputeSliceBound(t *testing.T) {
	tests := []struct {
		limit         int
		offset        int
		length        int
		expectedStart int
		expectedEnd   int
	}{
		{0, 0, 0, 0, 0},
		{100, 0, 2, 0, 2},
		{1000, 0, 1000, 0, 100},
		{1000, 10, 1000, 10, 110},
		{50, 12, 23, 12, 23},
		{50, 35, 23, 22, 23},
		{-1, -1, 10, 0, 1},
		{-1, -1, 0, 0, 0},
	}

	for i, tt := range tests {
		s, e := ComputeSliceBound(tt.limit, tt.offset, tt.length)
		if s != tt.expectedStart {
			t.Fatalf("%d : expected start %d but got %d", i, tt.expectedStart, s)
		}
		if e != tt.expectedEnd {
			t.Fatalf("%d : expected end %d but got %d", i, tt.expectedEnd, e)
		}
		if (e-s) > tt.limit && tt.limit > 0 {
			t.Fatalf("%d : set a limit of %d but got a range of %d ", i, tt.limit, e-s)
		}
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		current  int
		total    int
		max      int
		link     string
		expected Pages
	}{
		{
			current: 1,
			total:   1,
			max:     5,
			link:    "test",
			expected: Pages{
				Current:  1,
				Previous: 0,
				Next:     0,
				Total:    1,
				All:      []int{1},
				Link:     "test",
			},
		},
		{
			current: 1,
			total:   2,
			max:     5,
			link:    "test",
			expected: Pages{
				Current:  1,
				Previous: 0,
				Next:     2,
				Total:    2,
				All:      []int{1, 2},
				Link:     "test",
			},
		},
		{
			current: 2,
			total:   4,
			max:     5,
			link:    "test",
			expected: Pages{
				Current:  2,
				Previous: 1,
				Next:     3,
				Total:    4,
				All:      []int{1, 2, 3, 4},
				Link:     "test",
			},
		},
	}

	for i, tt := range tests {
		result := New(tt.current, tt.total, tt.max, tt.link)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Fatalf("%d : \ngot %#v\nexp %#v\n", i, result, tt.expected)
		}
	}
}
