package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sample() [][]string {
	return [][]string{
		{
			"L.LL.LL.LL",
			"LLLLLLL.LL",
			"L.L.L..L..",
			"LLLL.LL.LL",
			"L.LL.LL.LL",
			"L.LLLLL.LL",
			"..L.L.....",
			"LLLLLLLLLL",
			"L.LLLLLL.L",
			"L.LLLLL.LL",
		},
		{
			"#.##.##.##",
			"#######.##",
			"#.#.#..#..",
			"####.##.##",
			"#.##.##.##",
			"#.#####.##",
			"..#.#.....",
			"##########",
			"#.######.#",
			"#.#####.##",
		},
		{
			"#.LL.L#.##",
			"#LLLLLL.L#",
			"L.L.L..L..",
			"#LLL.LL.L#",
			"#.LL.LL.LL",
			"#.LLLL#.##",
			"..L.L.....",
			"#LLLLLLLL#",
			"#.LLLLLL.L",
			"#.#LLLL.##",
		},
		{
			"#.##.L#.##",
			"#L###LL.L#",
			"L.#.#..#..",
			"#L##.##.L#",
			"#.##.LL.LL",
			"#.###L#.##",
			"..#.#.....",
			"#L######L#",
			"#.LL###L.L",
			"#.#L###.##",
		},
		{
			"#.#L.L#.##",
			"#LLL#LL.L#",
			"L.L.L..#..",
			"#LLL.##.L#",
			"#.LL.LL.LL",
			"#.LL#L#.##",
			"..L.L.....",
			"#L#LLLL#L#",
			"#.LLLLLL.L",
			"#.#L#L#.##",
		}, {
			"#.#L.L#.##",
			"#LLL#LL.L#",
			"L.#.L..#..",
			"#L##.##.L#",
			"#.#L.LL.LL",
			"#.#L#L#.##",
			"..L.L.....",
			"#L#L##L#L#",
			"#.LLLLLL.L",
			"#.#L#L#.##",
		},
	}
}

func TestCountNear(t *testing.T) {
	samples := sample()
	tests := []struct {
		sampleID int
		x, y     int
		want     int
	}{
		{sampleID: 1, x: 0, y: 0, want: 2},
		{sampleID: 1, x: 0, y: 1, want: 3},
	}

	for i, tc := range tests {
		var b *board
		b = InitFromString(samples[tc.sampleID])
		got := b.CountNear(tc.x, tc.y)
		if got != tc.want {
			t.Errorf("bad count at #%d: want %d got %d", i, tc.want, got)
		}
	}
}

func TestEvolve(t *testing.T) {
	samples := sample()
	tests := []struct {
		input      []string
		want       []string
		wantChange bool
	}{
		{
			input:      []string{"...", ".L.", "..."},
			want:       []string{"...", ".#.", "..."},
			wantChange: true,
		},
		{
			input:      samples[0],
			want:       samples[1],
			wantChange: true,
		},
		{
			input:      samples[1],
			want:       samples[2],
			wantChange: true,
		},
		{
			input:      samples[2],
			want:       samples[3],
			wantChange: true,
		},
		{
			input:      samples[3],
			want:       samples[4],
			wantChange: true,
		},
		{
			input:      samples[4],
			want:       samples[5],
			wantChange: true,
		},
		{
			input: samples[5],
			want:  samples[5],
		},
	}

	for i, tc := range tests {
		var input, want *board
		input = InitFromString(tc.input)
		want = InitFromString(tc.want)
		got := input.Evolve(4, false)
		if diff := cmp.Diff(want.Spaces, got.Spaces); diff != "" {
			t.Errorf("Bad new board in test #%d:\n%s", i, diff)
		}
		if tc.wantChange != (got.LastChange > 0) {
			t.Errorf("Bad at detecting change for test #%d: want %v got %v",
				i, tc.wantChange, got.LastChange)
		}
	}
}

func TestCountOccupied(t *testing.T) {
	var b *board
	b = InitFromString(sample()[5])
	got := b.CountOccupied()
	if got != 37 {
		t.Errorf("wrong number of passengers: want 37 got %d", got)
	}
}
