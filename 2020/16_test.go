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
	vs := map[string]*validator{
		"class": newValidator("class: 1-3 or 5-7"),
		"row":   newValidator("row: 6-11 or 33-44"),
		"seat":  newValidator("seat: 13-40 or 45-50"),
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

func TestIdentifyFields(t *testing.T) {
	vs := map[string]*validator{
		"class": newValidator("class: 0-1 or 4-19"),
		"row":   newValidator("row: 0-5 or 8-19"),
		"seat":  newValidator("seat: 0-13 or 16-19"),
	}

	ts := []*ticket{
		newTicket("3,9,18"),
		newTicket("15,1,5"),
		newTicket("5,14,9"),
	}

	got := identifyFields(vs, ts)
	want := []string{"row", "class", "seat"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad at identifying:\n%s", diff)
	}
}
