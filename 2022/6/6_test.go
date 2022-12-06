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

	want := []string{
		"mjqjpqmgbljsphdztnvjfqwrcgsmlb",
		"bvwbjplbgvbhsrlpgdmjqwftvncz",
		"nppdvjthqldpwncqszvftbrmjlhg",
		"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg",
		"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}

func TestOne(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	want := []int{7, 5, 6, 10, 11}
	for i, l := range data {
		got := one(l)

		if diff := cmp.Diff(want[i], got); diff != "" {
			t.Errorf(diff)
		}
	}
}

func TestTwo(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	want := []int{19, 23, 23, 29, 26}
	for i, l := range data {
		got := two(l)

		if diff := cmp.Diff(want[i], got); diff != "" {
			t.Errorf(diff)
		}
	}
}
