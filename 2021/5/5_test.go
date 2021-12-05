package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	got := readFile("sample.txt")
	want := []*line{
		{point{0, 9}, point{5, 9}, false},
		{point{8, 0}, point{0, 8}, true},
		{point{9, 4}, point{3, 4}, false},
		{point{2, 2}, point{2, 1}, false},
		{point{7, 0}, point{7, 4}, false},
		{point{6, 4}, point{2, 0}, true},
		{point{0, 9}, point{2, 9}, false},
		{point{3, 4}, point{1, 4}, false},
		{point{0, 0}, point{8, 8}, true},
		{point{5, 5}, point{8, 2}, true},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad lines:\n%s", diff)
	}

}

func TestFindDanger(t *testing.T) {
	ls := readFile("sample.txt")
	got := findDanger(ls, 2, false)
	if got != 5 {
		t.Errorf("bad at finding danger, got %d, want 5", got)
	}
	got = findDanger(ls, 2, true)
	if got != 12 {
		t.Errorf("bad at finding diagonal danger, got %d, want 12", got)
	}
}
