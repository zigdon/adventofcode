package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	want := []string{
		"[({(<(())[]>[[{[]{<()<>>", // incomplete
		"[(()[<>])]({[<{<<[]>>(",   // incomplete
		"{([(<{}[<>[]}>{[]{[(<()>", // 12 }
		"(((({<>}<{<{<>}{[]{[]{}",  // incomplete
		"[[<[([]))<([[{}[[()]]]",   // 8 )
		"[{[{({}]{}}([{[{{{}}([]",  // 7 ]
		"{<[[]]>}<{[{[{[]{()[[[]",  // incomplete
		"[<(<(<(<{}))><([]([]()",   // 10 )
		"<{([([[(<>()){}]>(<<{{",   // 16 >
		"<{([{{}}[<[[[<>{}]]]>[]]", // incomplete
	}

	got := readFile("sample.txt")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad read:\n%s", diff)
	}
}

func TestIsValid(t *testing.T) {
	data := readFile("sample.txt")
	want := []struct {
		pos int
		r   rune
	}{
		{0, 0},
		{0, 0},
		{12, '}'},
		{0, 0},
		{8, ')'},
		{7, ']'},
		{0, 0},
		{10, ')'},
		{16, '>'},
		{0, 0},
	}

	for i, l := range data {
		w := want[i]
		p, c := isValid(l)
		if p != w.pos || c != w.r {
			t.Errorf("wrong validation of %q: want %v, got (%d, %q)", l, w, p, c)
		}
	}
}

func TestScoreSet(t *testing.T) {
	data := readFile("sample.txt")
	want := 26397
	got := scoreSet(data)
	if want != got {
		t.Errorf("bad score: want %d, got %d", want, got)
	}
}

func TestCompleteLine(t *testing.T) {
	want := []struct {
		c string
		v int64
	}{
		{"}}]])})]", 288957},
		{")}>]})", 5566},
		{"}}>}>))))", 1480781},
		{"]]}}]}]}>", 995444},
		{"])}>", 294},
	}
	data := readFile("sample.txt")
	incomplete := []string{}
	for _, l := range data {
		if _, r := isValid(l); r == 0 {
			incomplete = append(incomplete, l)
		}
	}

	for i, l := range incomplete {
		got, val := complete(l)
		if diff := cmp.Diff(want[i].c, got); diff != "" {
			t.Errorf("bad complete for %q: want %q, got %q", l, want[i].c, got)
		}
		if val != want[i].v {
			t.Errorf("bad complete value for %q: want %d, got %d", l, want[i].v, val)
		}
	}
}
