package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	got := readFile("sample.txt")
	want := &Cave{
		Floor: false,
		Min:   &Point{494, 0},
		Max:   &Point{503, 9},
		Scan: map[Point]Object{
			{494, 9}: Rock,
			{495, 9}: Rock,
			{496, 9}: Rock,
			{497, 9}: Rock,
			{498, 9}: Rock,
			{499, 9}: Rock,
			{500, 9}: Rock,
			{501, 9}: Rock,
			{502, 9}: Rock,
			{502, 8}: Rock,
			{502, 7}: Rock,
			{502, 6}: Rock,
			{502, 5}: Rock,
			{502, 4}: Rock,
			{503, 4}: Rock,
			{496, 6}: Rock,
			{497, 6}: Rock,
			{498, 6}: Rock,
			{498, 5}: Rock,
			{498, 4}: Rock,
			{500, 0}: Source,
		},
	}

	t.Logf("\n%s", got)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}

func TestDrop(t *testing.T) {
	cave := readFile("sample.txt")
	tests := []*Point{
		{500, 8},
		{499, 8},
		{501, 8},
		{500, 7},
	}

	source := Point{500, 0}
	for n, tc := range tests {
		got := cave.Drop(source, 0)
		if cmp.Diff(tc, got) != "" {
			t.Errorf("Error after %d drops: want %s, got %s\n%s", n, tc, got, cave)
		}
	}

	for n := len(tests); n < 24; n++ {
		if cave.Drop(source, 0) == nil {
			t.Errorf("Drop %d is off the ledge", n)
		}
	}
	p := cave.Drop(source, 0)
	if p != nil {
		t.Errorf("Last drop ended up at %s\n%s", p, cave)
	}
	t.Log(cave)
}

func TestOne(t *testing.T) {
	data := readFile("sample.txt")
	got := one(data)
	want := 24
	t.Logf("\n%s", data)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}

func TestTwo(t *testing.T) {
	data := readFile("sample.txt")
	got := two(data)
	want := 93
	t.Logf("\n%s", data)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}
