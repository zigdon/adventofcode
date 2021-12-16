package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	data := readFile("sample.txt")
	want := board{
		Danger: [][]int{
			{1, 1, 6, 3, 7, 5, 1, 7, 4, 2},
			{1, 3, 8, 1, 3, 7, 3, 6, 7, 2},
			{2, 1, 3, 6, 5, 1, 1, 3, 2, 8},
			{3, 6, 9, 4, 9, 3, 1, 5, 6, 9},
			{7, 4, 6, 3, 4, 1, 7, 1, 1, 1},
			{1, 3, 1, 9, 1, 2, 8, 1, 3, 7},
			{1, 3, 5, 9, 9, 1, 2, 4, 2, 1},
			{3, 1, 2, 5, 4, 2, 1, 6, 3, 9},
			{1, 2, 9, 3, 1, 3, 8, 5, 2, 1},
			{2, 3, 1, 1, 9, 4, 4, 5, 8, 1},
		},
		Best: 0,
		Attempt: &trail{
			Steps: []*step{
				{Seen: map[coord]bool{{0, 0}: true}},
			},
		},
		Exit: coord{9, 9},
	}

	if diff := cmp.Diff(want, data); diff != "" {
		t.Errorf("bad at reading map:\n%s", diff)
	}
}

func TestFindPath(t *testing.T) {
	b := readFile("sample.txt")
	_, score := b.solve(0, 0)
	if score != 40 {
		t.Errorf("bad score, wanted 40, got %d", score)
	}
}

func TestMove(t *testing.T) {
	tests := []struct {
		d    dir
		want coord
	}{
		{
			d:    UP,
			want: coord{5, 4},
		},
		{
			d:    DOWN,
			want: coord{5, 6},
		},
		{
			d:    LEFT,
			want: coord{4, 5},
		},
		{
			d:    RIGHT,
			want: coord{6, 5},
		},
	}

	start := &step{
		Pos: coord{5, 5},
		Seen: map[coord]bool{
			{5, 5}: true,
			{4, 5}: true,
			{3, 5}: true,
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s", tc.d), func(t *testing.T) {
			got := move(start, tc.d)
			if cmp.Diff(tc.want, got.Pos) != "" {
				t.Errorf("bad move: want %v, got %v", tc.want, got.Pos)
			}
			for k := range start.Seen {
				if !got.Seen[k] {
					t.Errorf("seen not copied properly, missing %v", k)
				}
			}
			if !got.Seen[got.Pos] {
				t.Error("move not marked!")
			}
		})
	}
}

func TestGuess(t *testing.T) {
	tests := []struct {
		desc string
		in   *step
		want *step
	}{
		{
			desc: "normal move",
			in: &step{
				Pos:    coord{1, 1},
				Danger: 5,
				Seen:   map[coord]bool{{1, 1}: true},
			},
			want: &step{
				Pos:    coord{2, 1},
				Danger: 13,
				Seen:   map[coord]bool{{1, 1}: true, {2, 1}: true},
			},
		},
		{
			desc: "avoid repeating moves",
			in: &step{
				Pos:    coord{1, 1},
				Danger: 5,
				Seen:   map[coord]bool{{1, 1}: true, {2, 1}: true},
			},
			want: &step{
				Pos:    coord{1, 2},
				Danger: 6,
				Seen:   map[coord]bool{{1, 1}: true, {2, 1}: true, {1, 2}: true},
			},
		},
		{
			desc: "no move",
			in: &step{
				Pos:    coord{1, 1},
				Danger: 23,
				Seen: map[coord]bool{
					{1, 1}: true,
					{0, 1}: true,
					{1, 0}: true,
					{2, 1}: true,
					{1, 2}: true,
				},
			},
		},
	}

	b := board{
		Danger: [][]int{
			{1, 1, 6},
			{1, 3, 8},
			{2, 1, 3},
		},
		Best: 20,
		Exit: coord{2, 2},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			b.Attempt = &trail{Steps: []*step{tc.in}}
			got := b.guess()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad guess:\n%s", diff)
			}
		})
	}
}
