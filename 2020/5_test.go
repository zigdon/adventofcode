package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	sample = []*seat{
		{
			BSP: "FBFBBFFRLR",
			Row: 44,
			Col: 5,
			Id:  357,
		},
		{
			BSP: "BFFFBBFRRR",
			Row: 70,
			Col: 7,
			Id:  567,
		},
		{
			BSP: "FFFBBBFRRR",
			Row: 14,
			Col: 7,
			Id:  119,
		},
		{
			BSP: "BBFFBBFRLL",
			Row: 102,
			Col: 4,
			Id:  820,
		},
	}
)

func TestFindOpenSeat(t *testing.T) {
	tests := []struct {
		in   []bool
		want int
	}{
		{[]bool{false, false, false}, -1},
		{[]bool{false, false, true, false, true}, 3},
		{[]bool{true, false, true, false, true}, 1},
	}

	for i, tc := range tests {
		got := findOpenSeat(tc.in)
		if got != tc.want {
			t.Errorf("Found wrong seat in #%d: want %d, got %d", i, tc.want, got)
		}
	}
}

func TestParseBSP(t *testing.T) {
	for i, s := range sample {
		testSeat := &seat{BSP: s.BSP}
		testSeat.parseBSP()
		if diff := cmp.Diff(testSeat, s); diff != "" {
			t.Errorf("Bad test seat %d: %s", i, diff)
		}
	}
}
