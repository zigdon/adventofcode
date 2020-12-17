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

	want := map[int]bool{
		7:  true,
		13: true,
		59: true,
		31: true,
		19: true,
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
