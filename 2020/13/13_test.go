package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sample() []string {
	return []string{
		"939",
		"7,13,x,x,59,x,31,19",
	}
}

func TestParseSchedule(t *testing.T) {
	gotS, gotB := parseSchedule(sample())

	want := map[int]int{
		7:  0,
		13: 1,
		59: 4,
		31: 6,
		19: 7,
	}
	if gotS != 939 {
		t.Errorf("bad start: %d", gotS)
	}
	if diff := cmp.Diff(want, gotB); diff != "" {
		t.Errorf("bad schedule: %s", diff)
	}
}

func TestNextBus(t *testing.T) {
	start, sched := parseSchedule(sample())
	wait, line := nextBus(start, sched)
	if wait != 5 {
		t.Errorf("bad wait: %d", wait)
	}
	if line != 59 {
		t.Errorf("bad line: %d", line)
	}
}

func TestValidateTime(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{input: sample()[1], want: 1068781},
		{input: "17,x,13,19", want: 3417},
		{input: "67,7,59,61", want: 754018},
		{input: "67,x,7,59,61", want: 779210},
		{input: "67,7,x,59,61", want: 1261476},
		{input: "1789,37,47,1889", want: 1202161486},
	}

	for i, tc := range tests {
		_, sched := parseSchedule([]string{"0", tc.input})
		got := validateTime(tc.want, sched)
		if !got {
			t.Errorf("bad validation of %d #%d", tc.want, i)
		}
	}
}

func TestFindSequence(t *testing.T) {
	tests := []struct {
		input string
		want  int
		max   int
	}{
		{input: sample()[1], want: 1068781, max: 59},
		{input: "17,x,13,19", want: 3417, max: 19},
		{input: "67,7,59,61", want: 754018, max: 67},
		{input: "67,x,7,59,61", want: 779210, max: 67},
		{input: "67,7,x,59,61", want: 1261476, max: 67},
		{input: "1789,37,47,1889", want: 1202161486, max: 1889},
	}

	for i, tc := range tests {
		_, sched := parseSchedule([]string{"0", tc.input})
		got := findSequence(tc.want-2*tc.max, tc.want+2000, sched, false)
		if tc.want != got {
			t.Errorf("bad sequence found for #%d: want %d, got %d", i, tc.want, got)
		}
	}
}
