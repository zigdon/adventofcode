package main

import (
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	data := readFile("sample.txt")
	want := &caves{
		Edges: map[string]*set{
			"start": newSet("A", "b"),
			"A":     newSet("start", "b", "c", "end"),
			"b":     newSet("start", "A", "d", "end"),
			"c":     newSet("A"),
			"d":     newSet("b"),
			"end":   newSet("A", "b"),
		},
		SmallCaves: newSet("start", "b", "c", "d", "end"),
	}

	if diff := cmp.Diff(want, data); diff != "" {
		t.Errorf("bad cave map:\n%s", diff)
	}
}

func TestFindPaths(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"sample.txt", []string{
			"start,A,b,A,c,A,end",
			"start,A,b,A,end",
			"start,A,b,end",
			"start,A,c,A,b,A,end",
			"start,A,c,A,b,end",
			"start,A,c,A,end",
			"start,A,end",
			"start,b,A,c,A,end",
			"start,b,A,end",
			"start,b,end"},
		},
		{"sample2.txt", []string{
			"start,HN,dc,HN,end",
			"start,HN,dc,HN,kj,HN,end",
			"start,HN,dc,end",
			"start,HN,dc,kj,HN,end",
			"start,HN,end",
			"start,HN,kj,HN,dc,HN,end",
			"start,HN,kj,HN,dc,end",
			"start,HN,kj,HN,end",
			"start,HN,kj,dc,HN,end",
			"start,HN,kj,dc,end",
			"start,dc,HN,end",
			"start,dc,HN,kj,HN,end",
			"start,dc,end",
			"start,dc,kj,HN,end",
			"start,kj,HN,dc,HN,end",
			"start,kj,HN,dc,end",
			"start,kj,HN,end",
			"start,kj,dc,HN,end",
			"start,kj,dc,end",
		}},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			data := readFile(tc.input)
			got := data.findPaths("start", emptySet())
			sort.Strings(got)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad paths:\n%s", diff)
			}
		})
	}
}

func TestFindLongerPaths(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"sample.txt", []string{
			"start,A,b,A,b,A,c,A,end",
			"start,A,b,A,b,A,end",
			"start,A,b,A,b,end",
			"start,A,b,A,c,A,b,A,end",
			"start,A,b,A,c,A,b,end",
			"start,A,b,A,c,A,c,A,end",
			"start,A,b,A,c,A,end",
			"start,A,b,A,end",
			"start,A,b,d,b,A,c,A,end",
			"start,A,b,d,b,A,end",
			"start,A,b,d,b,end",
			"start,A,b,end",
			"start,A,c,A,b,A,b,A,end",
			"start,A,c,A,b,A,b,end",
			"start,A,c,A,b,A,c,A,end",
			"start,A,c,A,b,A,end",
			"start,A,c,A,b,d,b,A,end",
			"start,A,c,A,b,d,b,end",
			"start,A,c,A,b,end",
			"start,A,c,A,c,A,b,A,end",
			"start,A,c,A,c,A,b,end",
			"start,A,c,A,c,A,end",
			"start,A,c,A,end",
			"start,A,end",
			"start,b,A,b,A,c,A,end",
			"start,b,A,b,A,end",
			"start,b,A,b,end",
			"start,b,A,c,A,b,A,end",
			"start,b,A,c,A,b,end",
			"start,b,A,c,A,c,A,end",
			"start,b,A,c,A,end",
			"start,b,A,end",
			"start,b,d,b,A,c,A,end",
			"start,b,d,b,A,end",
			"start,b,d,b,end",
			"start,b,end",
		}},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			data := readFile(tc.input)
			got := data.findLongerPaths()
			sort.Strings(got)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad paths:\n%s", diff)
			}
		})
	}
}

func TestCountPaths(t *testing.T) {
	tests := []struct {
		input string
		f     func(*caves) []string
		want  int
	}{
		{
			input: "sample3.txt",
			f: func(c *caves) []string {
				return c.findPaths("start", emptySet())
			},
			want: 226,
		},
		{
			input: "sample2.txt",
			f: func(c *caves) []string {
				return c.findLongerPaths()
			},
			want: 103,
		},
		{
			input: "sample3.txt",
			f: func(c *caves) []string {
				return c.findLongerPaths()
			},
			want: 3509,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s(%d)", tc.input, tc.want), func(t *testing.T) {
			data := readFile(tc.input)
			got := tc.f(data)
			if len(got) != tc.want {
				t.Errorf("bad count of paths: want %d, got %d", tc.want, len(got))
			}
		})
	}
}
