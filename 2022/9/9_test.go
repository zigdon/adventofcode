package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	inst, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	want := []Inst{
		{1, 0, 4},  // R 4
		{0, 1, 4},  // U 4
		{-1, 0, 3}, // L 3
		{0, -1, 1}, // D 1
		{1, 0, 4},  // R 4
		{0, -1, 1}, // D 1
		{-1, 0, 5}, // L 5
		{1, 0, 2},  // R 2
	}

	if diff := cmp.Diff(want, inst); diff != "" {
		t.Error(diff)
	}
}

func TestOne(t *testing.T) {
	i, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := one(i)
	want := 13

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}

func TestTwo(t *testing.T) {
	i, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := two(i)
	want := 0

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}
