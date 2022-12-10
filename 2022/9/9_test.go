package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	inst, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	want := []Inst{
		{1, 0, 4},  // R 4
		{0, 1, 4},  // U 4
		{-1, 0, 3}, // L 3
		{0, -1, 1}, // D 1
		{1, 0, 4},  // R 4
		{0, -1, 1}, // D 1
		{-1, 0, 5}, // L 5
		{1, 0, 2},  // R 2
	}

	if diff := cmp.Diff(want, inst); diff != "" {
		t.Error(diff)
	}
}

func TestNudge(t *testing.T) {
	tests := []struct {
		desc   string
		h      *Point
		dx, dy int
		want   *Point
	}{
		{
			desc: "overlap -> no move",
			dx:   1,
		},
		{
			desc: "return to overlap",
			h:    &Point{1, 0},
			dx:   -1,
		},
		{
			desc: "drag right",
			h:    &Point{1, 0},
			dx:   1,
			want: &Point{1, 0},
		},
		{
			desc: "drag up",
			h:    &Point{0, 1},
			dy:   1,
			want: &Point{0, 1},
		},
		{
			desc: "diag -> no move",
			h:    &Point{1, 0},
			dx:   -1,
			dy:   1,
		},
		{
			desc: "to diag -> no move",
			h:    &Point{1, 0},
			dy:   1,
		},
		{
			desc: "diag -> up",
			h:    &Point{1, 1},
			dy:   1,
			want: &Point{1, 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			r := NewRope(2)
			if tc.want == nil {
				tc.want = &Point{}
			}
			if tc.h != nil {
				r.Head().Set(tc.h)
			}
			r.Nudge(0, tc.dx, tc.dy)
			if r.Tail().x != tc.want.x || r.Tail().y != tc.want.y {
				t.Errorf("Bad tail: want %s, got %s", tc.want, r.Tail())
			}
		})
	}
}

func TestOne(t *testing.T) {
	i, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := one(i)
	want := 13

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}

func TestTwo(t *testing.T) {
	i, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := two(i)
	want := 1

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}
