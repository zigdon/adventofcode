package common

import (
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
		desc    string
		in      interface{}
		want    interface{}
		wantErr bool
	}{
		{
			desc: "fine line",
			in:   "one",
			want: "one",
		},
		{
			desc:    "blank line",
			in:      "    ",
			want:    "",
			wantErr: true,
		},
		{
			desc:    "not string",
			in:      []int{1, 2},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := IgnoreBlankLines(0, tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad blanks:\n%s", diff)
			}
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: want %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestSplitWords(t *testing.T) {
	tests := []struct {
		desc    string
		in      interface{}
		want    interface{}
		wantErr bool
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
			desc:    "not string",
			in:      []int{1, 2},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := SplitWords(0, tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad words:\n%s", diff)
			}
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: want %v, got %v", tc.wantErr, err)
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
			desc:    "not strings",
			in:      20,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := Ints(0, tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bad ints:\n%s", diff)
			}
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: want %v, got %v", tc.wantErr, err)
			}
		})
	}
}
