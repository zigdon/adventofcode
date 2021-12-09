package main

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	sample := [][]int{
		{2, 1, 9, 9, 9, 4, 3, 2, 1, 0},
		{3, 9, 8, 7, 8, 9, 4, 9, 2, 1},
		{9, 8, 5, 6, 7, 8, 9, 8, 9, 2},
		{8, 7, 6, 7, 8, 9, 6, 7, 8, 9},
		{9, 8, 9, 9, 9, 6, 5, 6, 7, 8},
	}

	want := hMap{}
	for y, l := range sample {
		for x, h := range l {
			want[coord{x, y}] = h
		}
	}

	got := readFile("sample.txt")

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad at reading:\n%s", diff)
	}
}

func TestFindLocalLow(t *testing.T) {
	data := readFile("sample.txt")
	got, _ := findLocalLow(data)
	sort.Ints(got)
	want := []int{0, 1, 5, 5}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad lows:\n%s", diff)
	}
}

func TestFindRisk(t *testing.T) {
	data := readFile("sample.txt")
	got := findRisk(data)
	want := 15
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad risk:\n%s", diff)
	}
}

func TestFindBasinSize(t *testing.T) {
	data := readFile("sample.txt")
	tests := []struct {
		x, y int
		want int
	}{
		{1, 0, 3},
		{9, 0, 9},
		{2, 2, 14},
		{6, 4, 9},
	}

	for _, tc := range tests {
		got := data.basinSize(coord{tc.x, tc.y})
		if got != tc.want {
			t.Errorf("bad size at %d,%d: want %d, got %d", tc.x, tc.y, tc.want, got)
		}
	}

}
