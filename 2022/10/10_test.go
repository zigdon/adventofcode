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

	want := []struct {
		name string
		arg  int
	}{
		{"noop", 0},
		{"addx", 3},
		{"addx", -5},
	}

	for n, op := range got {
		if op.Name != want[n].name {
			t.Errorf("Wrong op #%d: want %s, got %s", n, want[n].name, op.Name)
		}
		if op.Arg != want[n].arg {
			t.Errorf("Wrong arg #%d: want %d, got %d", n, want[n].arg, op.Arg)
		}
	}
}

func TestOne(t *testing.T) {
	data, err := readFile("sample2.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := one(data)
	want := 13140

	if got != want {
		t.Errorf("Wrong read: want %d, got %d", want, got)
	}
}

func TestTwo(t *testing.T) {
	data, err := readFile("sample2.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := two(data)
	want := []string{
		"##..##..##..##..##..##..##..##..##..##..",
		"###...###...###...###...###...###...###.",
		"####....####....####....####....####....",
		"#####.....#####.....#####.....#####.....",
		"######......######......######......####",
		"#######.......#######.......#######.....",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}
