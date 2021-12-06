package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCount(t *testing.T) {
	fish := readFile("sample.txt")
	got := fish.Size()
	if got != 5 {
		t.Errorf("bad at counting fish, wanted 5, got %d", got)
	}

}

func TestReadFile(t *testing.T) {
	got := readFile("sample.txt")
	//            0  1  2  3  4  5  6  7  8
	want := &school{
		Ages: []int64{0, 1, 1, 2, 1, 0, 0, 0, 0},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad school:\n%s", diff)
	}
}

func TestBreed(t *testing.T) {
	fish := readFile("sample.txt")
	want := []int64{5, 6, 7, 9, 10, 10, 10, 10, 11, 12, 15, 17, 19, 20, 20, 21, 22, 26}

	for day := 0; day < 18; day++ {
		got := fish.Breed(1)
		if got != want[day] {
			t.Errorf("wrong number of fish at day %d: got %d, want %d", day+1, got, want[day])
		}
	}

	got := fish.Breed(80 - 18)
	if got != 5934 {
		t.Errorf("wrong number of fish after 80 days: got %d, want %d", got, 5934)
	}

	got = fish.Breed(256 - 80)
	if got != 26984457539 {
		t.Errorf("wrong number of fish after 256 days: got %d, want %d", got, int64(26984457539))
	}
}
