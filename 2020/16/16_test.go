package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewValidator(t *testing.T) {
	tests := []struct {
		line string
		name string
		good []int
		bad  []int
	}{{
		line: "seat: 45-475 or 489-960",
		name: "seat",
		good: []int{45, 300, 500, 960},
		bad:  []int{40, 480, 961},
	}}

	for _, tc := range tests {
		v := newValidator(tc.line)
		if v.name != tc.name {
			t.Errorf("bad name: want %q, got %q", tc.name, v.name)
		}
		for _, n := range tc.good {
			if !v.ok(n) {
				t.Errorf("bad invalid %v: %d", v.rules, n)
			}
		}
		for _, n := range tc.bad {
			if v.ok(n) {
				t.Errorf("bad valid %v: %d", v.rules, n)
			}
		}
	}
}

func TestNewTicket(t *testing.T) {
	tic := newTicket("191,139,59,79,149,83,67,73,167,181,173,61,53,137,71,163,179,193,107,197")
	want := []int{191, 139, 59, 79, 149, 83, 67, 73, 167, 181, 173, 61, 53, 137, 71, 163, 179, 193, 107, 197}
	if diff := cmp.Diff(want, tic.fs); diff != "" {
		t.Errorf("bad ticket:\n%s", diff)
	}
}

func TestFindBadFields(t *testing.T) {
	vs := []*validator{
		newValidator("class: 1-3 or 5-7"),
		newValidator("row: 6-11 or 33-44"),
		newValidator("seat: 13-40 or 45-50"),
	}

	tests := []struct {
		tic  *ticket
		want []int
	}{
		{tic: newTicket("7,3,47"), want: []int{}},
		{tic: newTicket("40,4,50"), want: []int{4}},
		{tic: newTicket("55,2,20"), want: []int{55}},
		{tic: newTicket("38,6,12"), want: []int{12}},
	}

	for _, tc := range tests {
		if diff := cmp.Diff(tc.want, findBadFields(tc.tic, vs)); diff != "" {
			t.Errorf("wrong bad fields: %s", diff)
		}
	}
}
