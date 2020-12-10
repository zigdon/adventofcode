package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sample1() []int {
	return []int{16, 10, 15, 5, 1, 11, 7, 19, 6, 12, 4}
}

func sample2() []int {
	return []int{28, 33, 18, 42, 31, 14, 46, 20, 48, 47, 24, 23, 49,
		45, 19, 38, 39, 11, 1, 32, 25, 35, 8, 17, 7, 9, 4, 2, 34, 10, 3}
}

func TestGetDist(t *testing.T) {
	tests := []struct {
		data []int
		want map[int]int
	}{
		{data: sample1(), want: map[int]int{1: 7, 3: 5}},
		{data: sample2(), want: map[int]int{1: 22, 3: 10}},
	}

	for _, tc := range tests {
		got := getDist(tc.data)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("bad dist: %s", diff)
		}
	}
}
