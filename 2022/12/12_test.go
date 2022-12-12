package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	got := readFile("sample.txt")
	if diff := cmp.Diff(Point{8, 5}, got.Size); diff != "" {
		t.Errorf("Wrong size:\n%s", diff)
	}
	if diff := cmp.Diff(Point{0, 0}, got.Start); diff != "" {
		t.Errorf("Wrong start:\n%s", diff)
	}
	if diff := cmp.Diff(Point{5, 2}, got.End); diff != "" {
		t.Errorf("Wrong end:\n%s", diff)
	}
	want := [][]int{
		{0, 0, 1, 16, 15, 14, 13, 12},
		{0, 1, 2, 17, 24, 23, 23, 11},
		{0, 2, 2, 18, 25, 25, 23, 10},
		{0, 2, 2, 19, 20, 21, 22, 9},
		{0, 1, 3, 4, 5, 6, 7, 8},
	}

	for y, l := range want {
		for x, a := range l {
			p := Point{x, y}
			if got.Alt[p] != a {
				t.Errorf("Wrong height at (%s): want %d, got %d", p, a, got.Alt[p])
			}
		}
	}
}

func TestOne(t *testing.T) {
	data := readFile("sample.txt")
	p := one(data)
	if p == nil {
		t.Fatalf("No path found")
	}
	if len(p.Steps)-1 != 31 {
		t.Logf("\n%s", data)
		t.Errorf("Bad path: want %d, got %d:\n%s", 31, len(p.Steps), p)
	}
}

func TestTwo(t *testing.T) {
	data := readFile("sample.txt")
	p := two(data)
	if p == nil {
		t.Fatalf("No path found")
	}
	if len(p.Steps)-1 != 29 {
		t.Logf("\n%s", data)
		t.Errorf("Bad path: want %d, got %d:\n%s", 29, len(p.Steps), p)
	}
}
