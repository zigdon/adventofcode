package main

import (
	"testing"
)

func sample() []string {
	return []string{
		"nop +0",
		"acc +1",
		"jmp +4",
		"acc +3",
		"jmp -3",
		"acc -99",
		"acc +1",
		"jmp -4",
		"acc +6",
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		inst     []string
		wantAcc  int
		wantLoop int
	}{
		{
			[]string{"nop +0"},
			0, -1,
		},
		{
			[]string{"acc +1", "acc -2", "nop +0", "acc +5"},
			4, -1,
		},
		{
			[]string{"acc +1", "jmp +2", "acc -2", "nop +0", "acc +5"},
			6, -1,
		},
		{
			[]string{"jmp +3", "acc -2", "jmp +3", "acc +5", "jmp -3", "nop +0"},
			3, -1,
		},
		{
			[]string{"acc +1", "acc -2", "nop +0", "jmp -2"},
			-1, 1,
		},
		{
			sample(),
			5, 1,
		},
	}

	for i, tc := range tests {
		code := compile(tc.inst)
		got, loop := run(code)

		if got != tc.wantAcc {
			t.Errorf("Bad acc value for test #%d: want %d, got %d", i, tc.wantAcc, got)
		}
		if loop != tc.wantLoop {
			t.Errorf("Bad loop detected for test #%d: want %d, got %d", i, tc.wantLoop, loop)
		}
	}
}
