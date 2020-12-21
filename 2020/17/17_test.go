package main

import (
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type point struct {
	xyz coord
	val bool
}

func coordSort(data []coord) func(int, int) bool {
	return func(j, k int) bool {
		a := data[j]
		b := data[k]
		if a.X < b.X {
			return false
		}
		if a.X > b.X {
			return true
		}
		if a.Y < b.Y {
			return false
		}
		if a.Y > b.Y {
			return true
		}
		if a.Z < b.Z {
			return false
		}
		if a.Z > b.Z {
			return true
		}
		return false
	}
}

func _TestGetSet(t *testing.T) {
	tests := []struct {
		pts  []point
		want []coord
	}{
		{
			pts: []point{
				{coord{1, 1, 1, 0}, true},
				{coord{1, 2, 1, 0}, true},
				{coord{2, -1, 2, 0}, true},
			},
			want: []coord{coord{1, 1, 1, 0}, coord{1, 2, 1, 0}, coord{2, -1, 2, 0}},
		},
		{
			pts: []point{
				{coord{1, 1, 1, 0}, true},
				{coord{1, 1, 1, 0}, false},
				{coord{2, -1, 2, 0}, false},
				{coord{2, -1, 2, 0}, true},
			},
			want: []coord{coord{2, -1, 2, 0}},
		},
	}

	for _, tc := range tests {
		s := newSpace()
		for _, p := range tc.pts {
			s.set(p.xyz, p.val)
		}

		got := s.getTrue()
		sort.Slice(got, coordSort(got))
		sort.Slice(tc.want, coordSort(tc.want))

		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("bad values: %s", diff)
		}
	}
}

func _TestLoadSlice(t *testing.T) {
	data := strings.Join([]string{".#.", "..#", "###"}, "\n")
	s := newSpace()
	s.loadSlice(data, 2, 0)
	got := s.getTrue()
	want := []coord{
		{1, 0, 2, 0},
		{2, 1, 2, 0},
		{0, 2, 2, 0},
		{1, 2, 2, 0},
		{2, 2, 2, 0}}
	sort.Slice(got, coordSort(got))
	sort.Slice(want, coordSort(want))
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad slicing: %s", diff)
	}
}

func _TestCount(t *testing.T) {
	s := newSpace()
	data := strings.Join([]string{".#.", "..#", "###"}, "\n")
	s.loadSlice(data, 0, 0)
	data = strings.Join([]string{"#.#"}, "\n")
	s.loadSlice(data, 1, 0)
	tests := []struct {
		x, y, z int
		want    int
	}{
		{0, 0, 0, 2},
		{1, 1, 0, 7},
		{0, 0, 2, 1},
	}

	for _, tc := range tests {
		got := s.count(coord{tc.x, tc.y, tc.z, 0})
		if got != tc.want {
			t.Errorf("bad count at %d,%d,%d: want %d, got %d", tc.x, tc.y, tc.z, tc.want, got)
		}
	}
}

func _TestEvolve(t *testing.T) {
	s := newSpace()
	data := strings.Join([]string{".#.", "..#", "###"}, "\n")
	s.loadSlice(data, 0, 0)
	t.Log("Before:")
	t.Log(s.dump())

	want := newSpace()
	want.loadSlice(strings.Join([]string{
		"#..",
		"..#",
		".#."}, "\n"), -1, 0)
	want.loadSlice(strings.Join([]string{
		"#.#",
		".##",
		".#."}, "\n"), 0, 0)
	want.loadSlice(strings.Join([]string{
		"#..",
		"..#",
		".#."}, "\n"), 1, 0)
	s.evolve()
	t.Log("After 1:")
	t.Log(s.dump())
	if diff := cmp.Diff(want.dump(), s.dump()); diff != "" {
		t.Errorf("bad evolution gen 1: %s", diff)
	}

	want = newSpace()
	want.loadSlice(strings.Join([]string{
		".....",
		".....",
		"..#..",
		".....",
		"....."}, "\n"), -2, 0)

	want.loadSlice(strings.Join([]string{
		"..#..",
		".#..#",
		"....#",
		".#...",
		"....."}, "\n"), -1, 0)

	want.loadSlice(strings.Join([]string{
		"##...",
		"##...",
		"#....",
		"....#",
		".###."}, "\n"), 0, 0)

	want.loadSlice(strings.Join([]string{
		"..#..",
		".#..#",
		"....#",
		".#...",
		"....."}, "\n"), 1, 0)

	want.loadSlice(strings.Join([]string{
		".....",
		".....",
		"..#..",
		".....",
		"....."}, "\n"), 2, 0)

	s.evolve()
	t.Log("After 2:")
	t.Log(s.dump())
	if diff := cmp.Diff(want.dump(), s.dump()); diff != "" {
		t.Errorf("bad evolution gen 2: %s", diff)
	}

	// gens 3-6
	for i := 3; i <= 6; i++ {
		s.evolve()
	}
	cnt := len(s.getTrue())
	if cnt != 112 {
		t.Errorf("wrong number of actives: want 112 got %d", cnt)
	}
}

func TestEvolve4d(t *testing.T) {
	s := newSpace()
	data := strings.Join([]string{".#.", "..#", "###"}, "\n")
	s.loadSlice(data, 0, 0)
	t.Log("Before:")
	t.Log(s.dump())

	want := newSpace()
	want.loadSlice(strings.Join([]string{
		"#..",
		"..#",
		".#."}, "\n"), -1, -1)
	want.loadSlice(strings.Join([]string{
		"#..",
		"..#",
		".#."}, "\n"), 0, -1)
	want.loadSlice(strings.Join([]string{
		"#..",
		"..#",
		".#."}, "\n"), 1, -1)
	want.loadSlice(strings.Join([]string{
		"#..",
		"..#",
		".#."}, "\n"), -1, 0)
	want.loadSlice(strings.Join([]string{
		"#.#",
		".##",
		".#."}, "\n"), 0, 0)
	want.loadSlice(strings.Join([]string{
		"#..",
		"..#",
		".#."}, "\n"), 1, 0)
	want.loadSlice(strings.Join([]string{
		"#..",
		"..#",
		".#."}, "\n"), -1, 1)
	want.loadSlice(strings.Join([]string{
		"#..",
		"..#",
		".#."}, "\n"), 0, 1)
	want.loadSlice(strings.Join([]string{
		"#..",
		"..#",
		".#."}, "\n"), 1, 1)

	s.evolve4d()
	t.Log("After 1:")
	got := s.dump()
	t.Log(got)
	if diff := cmp.Diff(want.dump(), got); diff != "" {
		t.Errorf("bad 4d evolution gen 1: %s", diff)
	}
}
