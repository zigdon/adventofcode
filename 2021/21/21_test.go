package main

import (
	"fmt"
	"math"
	"testing"
)

func TestReadFile(t *testing.T) {
	a, b := readFile("sample.txt")
	if a != 4 || b != 8 {
		t.Errorf("bad starting posititons, want %d and %d, got %d and %d", 4, 8, a, b)
	}
}

func TestDetDie(t *testing.T) {
	d := deterministict()
	for i := 1; i <= 100; i++ {
		got := d.roll()
		if got != i {
			t.Errorf("bad detRoll: want %d, got %d", i, got)
		}
	}
	got := d.roll()
	if got != 1 {
		t.Errorf("die didn't reset: got %d", got)
	}
}

func TestDiracDie(t *testing.T) {
	tests := []struct {
		id   int64
		want []int
	}{
		{0, []int{1, 1, 1, 1, 1}},
		{1, []int{2, 1, 1, 1, 1}},
		{2, []int{3, 1, 1, 1, 1}},
		{4, []int{2, 2, 1, 1, 1}},
		{200, []int{3, 1, 2, 2, 3, 1, 1}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d", tc.id), func(t *testing.T) {
			d := dirac(tc.id)
			n := 0
			for len(tc.want) > 0 {
				n++
				got := d.roll()
				if tc.want[0] != got {
					t.Errorf("bad roll #%d: want %d, got %d", n, tc.want[0], got)
				}
				tc.want = tc.want[1:]
			}
		})
	}
}

func TestDetCast(t *testing.T) {
	rolls := []int{1, 2, 3}
	want := []int{1, 5, 15}
	d := deterministict()
	for i, r := range rolls {
		got, _ := d.cast(r)
		if got != want[i] {
			t.Errorf("bad cast after case #%d: want %d, got %d", i, want[i], got)
		}
	}
}

func TestGame(t *testing.T) {
	g := newGame(2, []int{4, 8}, 1000)
	d := deterministict()
	winner := g.runGame(d)

	if winner != 0 {
		t.Errorf("wrong winner: %d", winner)
	}

	if g.loserScore() != 745 {
		t.Errorf("loser has wrong score: %d", g.loserScore())
	}

	if d.rolled != 993 {
		t.Errorf("wrong number of rolls: %d", d.rolled)
	}
}

func TestAllGames(t *testing.T) {
	a, b := readFile("sample.txt")
	g := newGame(2, []int{a, b}, 21)
	cache = make(map[key]wins)
	var got wins
	opts := int64(math.Pow(3, 6))
	for id := int64(0); id < opts; id++ {
		wins := allGames(g, id)
		t.Logf("%d: %v", id, wins)
		got.a += wins.a
		got.b += wins.b
	}
	want := wins{444356092776315, 341960390180808}
	if got.a != want.a || got.b != want.b {
		t.Errorf("bad all games:\nwant: %15d, %15d\n got: %15d, %15d",
			want.a, want.b, got.a, got.b)
	}

}
