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

	want := []int{}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}

func TestOne(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := one(data)
	want := 0

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
	want := 0

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}
