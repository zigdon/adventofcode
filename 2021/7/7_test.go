package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	got := readFile("sample.txt")
	want := []int{16, 1, 2, 0, 4, 2, 7, 1, 2, 14}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad file:\n%s", diff)
	}
}

func TestRequiredFuel(t *testing.T) {
	fleet := readFile("sample.txt")
	tests := []struct {
		dest int
		cost int
	}{
		{2, 37},
		{1, 41},
		{3, 39},
		{10, 71},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("align to %d", tc.dest), func(t *testing.T) {
			got := align(fleet, tc.dest)
			if got != tc.cost {
				t.Errorf("bad cost, want %d, got %d", tc.cost, got)
			}
		})
	}
}

func TestSum(t *testing.T) {
	for i := 1; i < 20; i++ {
		want := 0
		for n := 0; n <= i; n++ {
			want += n
		}
		got := sum(i)
		if got != want {
			t.Errorf("couldn't count to %d: want %d, got %d", i, want, got)
		}
	}
}

func TestRequiredFuel2(t *testing.T) {
	fleet := readFile("sample.txt")
	tests := []struct {
		dest int
		cost int
	}{
		{2, 206},
		{5, 168},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("align to %d", tc.dest), func(t *testing.T) {
			got := align2(fleet, tc.dest)
			if got != tc.cost {
				t.Errorf("bad cost, want %d, got %d", tc.cost, got)
			}
		})
	}
}

func TestPlan(t *testing.T) {
	fleet := readFile("sample.txt")
	got, cost := plan(fleet, align)
	if got != 2 {
		t.Errorf("bad plan, want 2, got %d", got)
	}
	if cost != 37 {
		t.Errorf("bad cost, want 37, got %d", cost)
	}
	got, cost = plan(fleet, align2)
	if got != 5 {
		t.Errorf("bad plan, want 5, got %d", got)
	}
	if cost != 168 {
		t.Errorf("bad cost, want 168, got %d", cost)
	}
}
