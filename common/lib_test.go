package common

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAsInts(t *testing.T) {
	tests := []struct {
		desc string
		in   []interface{}
		want []int
	}{
		{
			"good",
			[]interface{}{"10", "20", "30", "10"},
			[]int{10, 20, 30, 10},
		},
		{
			"ignore errors",
			[]interface{}{"10", "twenty", "30", "30.1"},
			[]int{10, 30},
		},
		{
			"now you're being silly",
			[]interface{}{"{", 20.5},
			[]int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := AsInts(tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad ints:\n%s", diff)
			}
		})
	}
}

func TestAsIntGrid(t *testing.T) {
	tests := []struct {
		desc string
		in   []interface{}
		want [][]int
	}{
		{
			desc: "single string",
			in:   []interface{}{"1"},
			want: [][]int{{1}},
		},
		{
			desc: "string array",
			in:   []interface{}{"1", "2"},
			want: [][]int{{1}, {2}},
		},
		{
			desc: "multiple lists",
			in:   []interface{}{[]interface{}{"1", "2"}, []interface{}{"3", "4"}},
			want: [][]int{{1, 2}, {3, 4}},
		},
		{
			desc: "multiple mixed lists",
			in:   []interface{}{[]interface{}{"1", 2}, []interface{}{3, "4"}},
			want: [][]int{{1, 2}, {3, 4}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := AsIntGrid(tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad int grid:\n%s", diff)
			}
		})
	}
}

func TestAsStrings(t *testing.T) {
	tests := []struct {
		desc string
		in   []interface{}
		want []string
	}{
		{
			"good",
			[]interface{}{"10", "20", "30", "10"},
			[]string{"10", "20", "30", "10"},
		},
		{
			"ignore non-strings",
			[]interface{}{"10", "twenty", []string{"30"}, "30.1"},
			[]string{"10", "twenty", "30.1"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := AsStrings(tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad strings:\n%s", diff)
			}
		})
	}
}

func TestIgnoreBlankLines(t *testing.T) {
	tests := []struct {
		desc     string
		in       interface{}
		want     interface{}
		wantSkip bool
	}{
		{
			desc: "fine line",
			in:   "one",
			want: "one",
		},
		{
			desc:     "blank line",
			in:       "    ",
			want:     "",
			wantSkip: true,
		},
		{
			desc: "not string",
			in:   []int{1, 2},
			want: []int{1, 2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, got, err := IgnoreBlankLines(0, tc.in)
			if (err == skipLine) != tc.wantSkip {
				t.Errorf("unexpected skipLine: want %v, got %v", tc.wantSkip, err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad blanks:\n%s", diff)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		desc string
		in   interface{}
		sep  string
		want interface{}
	}{
		{
			desc: "single word",
			in:   "one",
			sep:  "/",
			want: []string{"one"},
		},
		{
			desc: "many words",
			in:   "here are/words",
			sep:  "/",
			want: []string{"here are", "words"},
		},
		{
			desc: "long separator",
			in:   "here are words",
			sep:  " are ",
			want: []string{"here", "words"},
		},
		{
			desc: "empty bits",
			sep:  "/",
			in:   "/word ///another ",
			want: []string{"word ", "another "},
		},
		{
			desc: "substrings",
			sep:  "/",
			in:   []string{"first/entry", "//second/"},
			want: [][]string{{"first", "entry"}, {"second"}},
		},
		{
			desc: "not string",
			in:   []int{1, 2},
			sep:  "/",
			want: []int{1, 2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			f := Split(tc.sep)
			_, got, err := f(0, tc.in)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad split:\n%s", diff)
			}
		})
	}
}

func TestSplitWords(t *testing.T) {
	tests := []struct {
		desc string
		in   interface{}
		want interface{}
	}{
		{
			desc: "single word",
			in:   "one",
			want: []string{"one"},
		},
		{
			desc: "many words",
			in:   "here are words",
			want: []string{"here", "are", "words"},
		},
		{
			desc: "many spaces",
			in:   "here are   words",
			want: []string{"here", "are", "words"},
		},
		{
			desc: "trim spaces",
			in:   "word!    ",
			want: []string{"word!"},
		},
		{
			desc: "not string",
			in:   []int{1, 2},
			want: []int{1, 2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, got, err := SplitWords(0, tc.in)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad words:\n%s", diff)
			}
		})
	}
}

func TestInts(t *testing.T) {
	tests := []struct {
		desc    string
		in      interface{}
		want    interface{}
		wantErr bool
	}{
		{
			desc: "good",
			in:   "10",
			want: 10,
		},
		{
			desc:    "errors",
			in:      "twenty",
			want:    0,
			wantErr: true,
		},
		{
			desc: "not strings",
			in:   20,
			want: 20,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, got, err := Ints(0, tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad ints:\n%s", diff)
			}
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: want %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestRange(t *testing.T) {
	sample := strings.Join([]string{
		"10 two",
		"   ",
		"15",
		"165",
	}, "\n")
	samplePath := filepath.Join(t.TempDir(), "sample.txt")
	err := ioutil.WriteFile(samplePath, []byte(sample), 0644)
	if err != nil {
		t.Fatalf("Can't write sample.txt: %v", err)
	}

	tests := []struct {
		desc string
		fs   []Transformation
		want []interface{}
	}{
		{
			desc: "(2-): skip blanks, ints",
			fs:   []Transformation{Range(2, -1, IgnoreBlankLines, Ints)},
			want: []interface{}{"10 two", "   ", 15, 165},
		},
		{
			desc: "skip blanks, (2-): ints",
			fs:   []Transformation{IgnoreBlankLines, Range(2, -1, Ints)},
			want: []interface{}{"10 two", 15, 165},
		},
		{
			desc: "skip blanks, (2-3): ints",
			fs:   []Transformation{IgnoreBlankLines, Range(2, 3, Ints)},
			want: []interface{}{"10 two", 15, "165"},
		},
		{
			desc: "(-3): skip blanks, (0-2): split words, (2-3): ints",
			fs: []Transformation{
				Range(-1, 3, IgnoreBlankLines),
				Range(0, 2, SplitWords),
				Range(2, 3, Ints),
			},
			want: []interface{}{[]string{"10", "two"}, 15, "165"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := ReadTransformedFile(samplePath, tc.fs...)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Bad range:\n%s", diff)
			}
		})
	}
}

func TestBuffer(t *testing.T) {
	sample := strings.Join([]string{
		"a b c d",
		"   ",
		"aa ab ac",
		"ba bb bc",
		"",
		"ca cb cc",
		"da db dc",
		"",
	}, "\n")
	samplePath := filepath.Join(t.TempDir(), "sample.txt")
	err := ioutil.WriteFile(samplePath, []byte(sample), 0644)
	if err != nil {
		t.Fatalf("Can't write sample.txt: %v", err)
	}

	tests := []struct {
		desc string
		fs   []Transformation
		want interface{}
	}{
		{
			desc: "split words, then buffer",
			fs:   []Transformation{SplitWords, Block},
			want: []interface{}{
				[]interface{}{[]string{"a", "b", "c", "d"}},
				[]interface{}{[]string{"aa", "ab", "ac"}, []string{"ba", "bb", "bc"}},
				[]interface{}{[]string{"ca", "cb", "cc"}, []string{"da", "db", "dc"}},
			},
		},
		{
			desc: "just buffer",
			fs:   []Transformation{Block},
			want: []interface{}{
				[]interface{}{"a b c d", "   ", "aa ab ac", "ba bb bc"},
				[]interface{}{"ca cb cc", "da db dc"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := ReadTransformedFile(samplePath, tc.fs...)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Bad buffering:\n%s", diff)
			}
		})
	}

}
