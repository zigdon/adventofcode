package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	got, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	want := grid{
		Trees: [][]int{
			{3, 0, 3, 7, 3},
			{2, 5, 5, 1, 2},
			{6, 5, 3, 3, 2},
			{3, 3, 5, 4, 9},
			{3, 5, 3, 9, 0},
		},
		W: 5,
		H: 5,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}

func TestIsVisible(t *testing.T) {
	g, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	want := [][]bool{
		{true, true, true, true, true},
		{true, true, true, false, true},
		{true, true, false, true, true},
		{true, false, true, false, true},
		{true, true, true, true, true},
	}

	got := make([][]bool, g.H)
	g.Do(func(x, y int) {
		if got[y] == nil {
			got[y] = make([]bool, g.W)
		}
		got[y][x] = g.IsVisible(x, y)
	})

	t.Log("\n" + g.String())
	for y := range want {
		for x := range want[y] {
			if want[y][x] != got[y][x] {
				t.Errorf("Wrong visibility (%d, %d): want %v, got %v", x, y, want[y][x], got[y][x])
			}
		}
	}
}

func TestScenic(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	tests := []struct{ x, y, want int }{
		{2, 1, 4},
		{2, 3, 8},
	}

	for _, tc := range tests {
		t.Logf("%d, %d\n%s", tc.x, tc.y, data)
		got := data.ScenicScore(tc.x, tc.y)
		if got != tc.want {
			t.Errorf("Bad score (%d,%d): want %d, got %d", tc.x, tc.y, tc.want, got)
		}
	}

}

func TestOne(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := one(data)
	want := 21

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}

func TestTwo(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := two(data)
	want := 8

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}
