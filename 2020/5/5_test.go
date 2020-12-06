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

func TestParseBSP(t *testing.T) {
	for i, s := range sample {
		testSeat := &seat{BSP: s.BSP}
		testSeat.parseBSP()
		if diff := cmp.Diff(testSeat, s); diff != "" {
			t.Errorf("Bad test seat %d: %s", i, diff)
		}
	}
}
