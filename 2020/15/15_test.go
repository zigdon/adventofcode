package main

import "testing"

func TestGetSequence(t *testing.T) {
	tests := []struct {
		starting []int
		want     int
		pos      int
	}{
		// Part 1
		{starting: []int{0, 3, 6}, want: 436},
		{starting: []int{0, 3, 6}, pos: 4, want: 0},
		{starting: []int{0, 3, 6}, pos: 5, want: 3},
		{starting: []int{0, 3, 6}, pos: 6, want: 3},
		{starting: []int{0, 3, 6}, pos: 7, want: 1},
		{starting: []int{0, 3, 6}, pos: 8, want: 0},
		{starting: []int{0, 3, 6}, pos: 9, want: 4},
		{starting: []int{0, 3, 6}, pos: 10, want: 0},
		{starting: []int{1, 3, 2}, want: 1},
		{starting: []int{2, 1, 3}, want: 10},
		{starting: []int{1, 2, 3}, want: 27},
		{starting: []int{2, 3, 1}, want: 78},
		{starting: []int{3, 2, 1}, want: 438},
		{starting: []int{3, 1, 2}, want: 1836},
		// Part 2
		{starting: []int{0, 3, 6}, pos: 30000000, want: 175594},
		{starting: []int{1, 3, 2}, pos: 30000000, want: 2578},
		{starting: []int{2, 1, 3}, pos: 30000000, want: 3544142},
		{starting: []int{1, 2, 3}, pos: 30000000, want: 261214},
		{starting: []int{2, 3, 1}, pos: 30000000, want: 6895259},
		{starting: []int{3, 2, 1}, pos: 30000000, want: 18},
		{starting: []int{3, 1, 2}, pos: 30000000, want: 362},
	}

	for i, tc := range tests {
		pos := tc.pos
		if pos == 0 {
			pos = 2020
		}
		got := getSequence(tc.starting, pos)
		if got != tc.want {
			t.Errorf("bad sequence for test #%d, position %d: want %d got %d", i, pos, tc.want, got)
		}
	}
}
