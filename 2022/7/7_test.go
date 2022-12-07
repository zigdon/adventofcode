package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	got, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	want := &fs{
		Dirs: map[string][]string{
			"/a":   {"f", "g", "h.lst"},
			"/a/e": {"i"},
			"/d":   {"j", "d.log", "d.ext", "k"},
			"":     {"b.txt", "c.dat"},
		},
		Files: map[string]int{
			"/a/e/i":   584,
			"/a/f":     29116,
			"/a/g":     2557,
			"/a/h.lst": 62596,
			"/d/j":     4060174,
			"/d/d.log": 8033020,
			"/d/d.ext": 5626152,
			"/d/k":     7214296,
			"/b.txt":   14848514,
			"/c.dat":   8504156,
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}

func TestDU(t *testing.T) {
	fs, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	t.Logf("fs: %#v", fs)

	want := map[string]int{
		"/a/e": 584,
		"/a":   94853,
		"/d":   24933642,
		"":     48381165,
	}

	for d, size := range want {
		got := fs.Du(d)
		if got != size {
			t.Errorf("%s: want %d, got %d", d, size, got)
		}
	}

}

func TestOne(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := one(data)
	want := 95437

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}

func TestTwo(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := two(data)
	want := 24933642

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}
