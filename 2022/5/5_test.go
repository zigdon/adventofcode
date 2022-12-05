package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	stacks, insts, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	wantS := []*Stack{
		{[]byte("ZN")},
		{[]byte("MCD")},
		{[]byte("P")},
	}
	if diff := cmp.Diff(wantS, stacks); diff != "" {
		t.Error(diff)
	}

	wantI := []Instruction{
		{1, 2, 1, false},
		{3, 1, 3, false},
		{2, 2, 1, false},
		{1, 1, 2, false},
	}
	if diff := cmp.Diff(wantI, insts); diff != "" {
		t.Error(diff)
	}
}

func TestOne(t *testing.T) {
	stacks, insts, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := one(stacks, insts)
	if diff := cmp.Diff("CMZ", got); diff != "" {
		t.Error(diff)
	}
}

func TestTwo(t *testing.T) {
	stacks, insts, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := two(stacks, insts)
	if diff := cmp.Diff("MCD", got); diff != "" {
		t.Error(diff)
	}
}
