package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sampleString() string {
	return strings.Join([]string{"1-3 a: abcde", "1-3 b: cdefg", "2-9 c: ccccccccc", ""}, "\n")
}

func sampleData() []Pw {
	return []Pw{
		{1, 3, 'a', "abcde"},
		{1, 3, 'b', "cdefg"},
		{2, 9, 'c', "ccccccccc"},
	}
}

func TestParseData(t *testing.T) {
	got := parseData(sampleString())
	if diff := cmp.Diff(sampleData(), got, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Bad at reading: %s", diff)
	}
}

func TestFindInvalid(t *testing.T) {
	tests := []struct {
		desc string
		give []Pw
		want []Pw
	}{
		{
			"Sample data",
			sampleData(),
			[]Pw{{1, 3, 'b', "cdefg"}},
		},
		{
			"Too long",
			[]Pw{{1, 2, 'a', "aaa"}},
			[]Pw{{1, 2, 'a', "aaa"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := findInvalid(tc.give)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Bad results: %s", diff)
			}
		})
	}
}
