package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	got := readFile("sample.txt")
	want := []input{
		{In: []string{"BE", "CFBEGAD", "CBDGEF", "FGAECD", "CGEB", "FDCGE", "AGEBFD", "FECDB", "FABCD", "EDB"},
			Out: []string{"FDGACBE", "CEFDB", "CEFBGD", "GCBE"}},
		{In: []string{"EDBFGA", "BEGCD", "CBG", "GC", "GCADEBF", "FBGDE", "ACBGFD", "ABCDE", "GFCBED", "GFEC"},
			Out: []string{"FCGEDB", "CGB", "DGEBACF", "GC"}},
		{In: []string{"FGAEBD", "CG", "BDAEC", "GDAFB", "AGBCFD", "GDCBEF", "BGCAD", "GFAC", "GCB", "CDGABEF"},
			Out: []string{"CG", "CG", "FDCAGB", "CBG"}},
		{In: []string{"FBEGCD", "CBD", "ADCEFB", "DAGEB", "AFCB", "BC", "AEFDC", "ECDAB", "FGDECA", "FCDBEGA"},
			Out: []string{"EFABCD", "CEDBA", "GADFEC", "CB"}},
		{In: []string{"AECBFDG", "FBG", "GF", "BAFEG", "DBEFA", "FCGE", "GCBEA", "FCAEGB", "DGCEAB", "FCBDGA"},
			Out: []string{"GECF", "EGDCABF", "BGF", "BFGEA"}},
		{In: []string{"FGEAB", "CA", "AFCEBG", "BDACFEG", "CFAEDG", "GCFDB", "BAEC", "BFADEG", "BAFGC", "ACF"},
			Out: []string{"GEBDCFA", "ECBA", "CA", "FADEGCB"}},
		{In: []string{"DBCFG", "FGD", "BDEGCAF", "FGEC", "AEGBDF", "ECDFAB", "FBEDC", "DACGB", "GDCEBF", "GF"},
			Out: []string{"CEFG", "DCBEF", "FCGE", "GBCADFE"}},
		{In: []string{"BDFEGC", "CBEGAF", "GECBF", "DFCAGE", "BDACG", "ED", "BEDF", "CED", "ADCBEFG", "GEBCD"},
			Out: []string{"ED", "BCGAFE", "CDGBA", "CBGEF"}},
		{In: []string{"EGADFB", "CDBFEG", "CEGD", "FECAB", "CGB", "GBDEFCA", "CG", "FGCDAB", "EGFDB", "BFCEG"},
			Out: []string{"GBDFCAE", "BGC", "CG", "CGB"}},
		{In: []string{"GCAFB", "GCF", "DCAEBFG", "ECAGB", "GF", "ABCDEG", "GAEF", "CAFBGE", "FDBAC", "FEGBDC"},
			Out: []string{"FGAE", "CFGAB", "FG", "BAGCE"}},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad reading:\n%s", diff)
	}
}

func TestPrune(t *testing.T) {
	tests := []struct {
		desc string
		in   *wireMap
		src  string
		dst  rune
		want *wireMap
	}{
		{
			desc: "no change",
			in: &wireMap{
				Possible: map[rune]map[rune]bool{
					'A': {'A': true, 'B': false, 'C': true},
					'B': {'A': true, 'B': true, 'C': true},
					'C': {'A': true, 'B': false, 'C': true},
				},
			},
			src: "B",
			dst: 'B',
			want: &wireMap{
				Possible: map[rune]map[rune]bool{
					'A': {'A': true, 'B': false, 'C': true},
					'B': {'A': true, 'B': true, 'C': true},
					'C': {'A': true, 'B': false, 'C': true},
				},
			},
		},
		{
			desc: "some change",
			in: &wireMap{
				Possible: map[rune]map[rune]bool{
					'A': {'A': true, 'B': false, 'C': true},
					'B': {'A': true, 'B': true, 'C': true},
					'C': {'A': true, 'B': false, 'C': true},
				},
			},
			src: "A",
			dst: 'C',
			want: &wireMap{
				Possible: map[rune]map[rune]bool{
					'A': {'A': true, 'B': false, 'C': true},
					'B': {'A': true, 'B': true, 'C': false},
					'C': {'A': true, 'B': false, 'C': false},
				},
			},
		},
		{
			desc: "multiple sources",
			in: &wireMap{
				Possible: map[rune]map[rune]bool{
					'A': {'A': true, 'B': false, 'C': true},
					'B': {'A': true, 'B': true, 'C': true},
					'C': {'A': true, 'B': false, 'C': true},
				},
			},
			src: "AB",
			dst: 'C',
			want: &wireMap{
				Possible: map[rune]map[rune]bool{
					'A': {'A': true, 'B': false, 'C': true},
					'B': {'A': true, 'B': true, 'C': true},
					'C': {'A': true, 'B': false, 'C': false},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			tc.in.prune(tc.src, tc.dst)
			if diff := cmp.Diff(tc.want, tc.in); diff != "" {
				t.Errorf("failed to prune:\n%s", diff)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	tests := []struct {
		desc   string
		limits []struct{ segs, opts string }
		want   *wireMap
	}{
		{
			desc: "nop",
			limits: []struct{ segs, opts string }{
				{"CD", "CD"},
			},
			want: &wireMap{
				Possible: map[rune]map[rune]bool{
					'A': {'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true},
					'B': {'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true},
					'C': {'A': false, 'B': false, 'C': true, 'D': true, 'E': false, 'F': false, 'G': false},
					'D': {'A': false, 'B': false, 'C': true, 'D': true, 'E': false, 'F': false, 'G': false},
					'E': {'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true},
					'F': {'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true},
					'G': {'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true},
				},
				Solved: map[rune]rune{},
				Opts:   map[rune]int{},
				Pairs:  map[rune]string{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			w := newWireMap()
			w.Possible['C'] = map[rune]bool{
				'A': false, 'B': false, 'C': true, 'D': true, 'E': false, 'F': false, 'G': false,
			}
			w.Possible['D'] = map[rune]bool{
				'A': false, 'B': false, 'C': true, 'D': true, 'E': false, 'F': false, 'G': false,
			}
			for _, l := range tc.limits {
				w.limit(l.segs, l.opts)
			}
			if diff := cmp.Diff(tc.want, w); diff != "" {
				t.Errorf("failed to prune:\n%s", diff)
			}

		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		desc string
		pos  []string
		want bool
	}{
		{
			desc: "already solved",
			pos:  []string{"A", "B", "C", "D", "E", "F", "G"},
			want: true,
		},
		{
			desc: "one move away",
			pos:  []string{"A", "AB", "AC", "D", "E", "F", "G"},
			want: true,
		},
		{
			desc: "two moves away",
			pos:  []string{"A", "AB", "CAB", "D", "E", "F", "G"},
			want: true,
		},
		{
			desc: "not solved",
			pos:  []string{"AB", "AB", "C", "D", "E", "F", "G"},
		},
		{
			desc: "one move, still not solved",
			pos:  []string{"A", "ABC", "BC", "D", "E", "F", "G"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			wm := newWireMap()
			wm.init(tc.pos)
			got := wm.check()
			if got != tc.want {
				t.Errorf("bad check: want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestCount(t *testing.T) {
	data := readFile("sample.txt")
	flat := [][]string{}
	for _, d := range data {
		flat = append(flat, d.Out)
	}

	got := countEasy(flat)
	if got != 26 {
		t.Errorf("can't even count the easy ones: want 26, got %d", got)
	}
}

func TestSolve(t *testing.T) {
	wm := newWireMap()
	wm.solve(strings.Fields("ACEDGFB CDFBE GCDFA FBCAD DAB CEFABD CDFGEB EAFB CAGEDB AB CDFEB FCADB CDFEB CDBAF"))

	want := map[rune]rune{
		'A': 'C', 'B': 'F', 'C': 'G', 'D': 'A', 'E': 'B', 'F': 'D', 'G': 'E',
	}

	if diff := cmp.Diff(want, wm.Solved); diff != "" {
		t.Errorf("bad solution:\n%s", diff)
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "ACEDGFB CDFBE GCDFA FBCAD DAB CEFABD CDFGEB EAFB CAGEDB AB CDFEB FCADB CDFEB CDBAF",
			want: "5353",
		},
	}

	data := readFile("sample.txt")
	n := []string{"8394", "9781", "1197", "9361", "4873", "8418", "4548", "1625", "8717", "4315"}
	for i, d := range data {
		tests = append(tests, struct {
			in   string
			want string
		}{
			in:   strings.Join(append(d.In, d.Out...), " "),
			want: n[i]})
	}

	for _, tc := range tests {
		wm := newWireMap()
		in := strings.Fields(tc.in)
		wm.solve(in)
		got := wm.decode(in[10:])
		if got != tc.want {
			t.Errorf("bad decode: want %s, got %s", tc.want, got)
		}
	}
}
