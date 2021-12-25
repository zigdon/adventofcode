package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	algo, img := readFile("sample.txt")

	if len(algo.Bits) != 512 {
		t.Errorf("too short, wanted 512, got %d", len(algo.Bits))
	}

	want := image{
		Pixels: map[coord]bool{
			{0, 0}: true, {1, 0}: false, {2, 0}: false, {3, 0}: true, {4, 0}: false,
			{0, 1}: true, {1, 1}: false, {2, 1}: false, {3, 1}: false, {4, 1}: false,
			{0, 2}: true, {1, 2}: true, {2, 2}: false, {3, 2}: false, {4, 2}: true,
			{0, 3}: false, {1, 3}: false, {2, 3}: true, {3, 3}: false, {4, 3}: false,
			{0, 4}: false, {1, 4}: false, {2, 4}: true, {3, 4}: true, {4, 4}: true,
		},
		Max: coord{4, 4},
	}

	if diff := cmp.Diff(want.String(), img.String()); diff != "" {
		t.Errorf("bad image:\n%s", diff)
		t.Logf("Want:\n%s", want.String())
		t.Logf("Got:\n%s", img.String())
	}

	if !want.Max.eq(img.Max) {
		t.Errorf("bad max: want %s, got %s", want.Max, img.Max)
	}

	if !img.Min.eq(coord{0, 0}) {
		t.Errorf("bad mmin: want %s, got %s", want.Min, img.Min)
	}
}

func TestApply(t *testing.T) {
	algo, img := readFile("sample.txt")
	want := newImage([]string{
		".##.##.",
		"#..#.#.",
		"##.#..#",
		"####..#",
		".#..##.",
		"..##..#",
		"...#.#.",
	})

	got := img.apply(algo)

	t.Run("first algo", func(t *testing.T) {
		if !want.eq(got) {
			t.Errorf("bad after algo:\nWant:\n%s\nGot:\n%s", want.String(), got.String())
		}
	})

	got = got.apply(algo)
	want = newImage([]string{
		"...............",
		"...............",
		"...............",
		"..........#....",
		"....#..#.#.....",
		"...#.#...###...",
		"...#...##.#....",
		"...#.....#.#...",
		"....#.#####....",
		".....#.#####...",
		"......##.##....",
		".......###.....",
		"...............",
		"...............",
		"...............",
	})

	t.Run("second algo", func(t *testing.T) {
		if !want.eq(got) {
			t.Errorf("bad after second algo:\nWant:\n%s\nGot:\n%s", want.String(), got.String())
		}
	})

	if got.count() != 35 {
		t.Errorf("bad count: want 35, got %d", got.count())
	}
}

func TestEnhance(t *testing.T) {
	algo, img := readFile("sample.txt")
	cnt := img.enhance(algo, 2).count()
	if cnt != 35 {
		t.Errorf("bad count after 2 iterationns: want 35, got %d", cnt)
	}
	cnt = img.enhance(algo, 50).count()
	if cnt != 3351 {
		t.Errorf("bad count after 50 iterationns: want 3351, got %d", cnt)
	}
}
