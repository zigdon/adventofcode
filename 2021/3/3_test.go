package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zigdon/adventofcode/common"
)

func TestCountBits(t *testing.T) {
	data := common.AsStrings(common.ReadTransformedFile("sample.txt", common.IgnoreBlankLines))
	got := countBits(data)
	want := []int{7, 5, 8, 7, 5}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Bad at counting: -want +got\n%s", diff)
	}
}

func TestGetRates(t *testing.T) {
	g, e := getRates([]int{7, 5, 8, 7, 5}, 12)
	if g != 22 {
		t.Errorf("Bad gamma: want 22, got %d (%b)", g, g)
	}

	if e != 9 {
		t.Errorf("Bad epsilon: want 9, got %d (%b)", e, e)
	}
}

func TestGetRatings(t *testing.T) {
	data := common.AsStrings(common.ReadTransformedFile("sample.txt", common.IgnoreBlankLines))
	tests := []struct {
		desc string
		ge   bool
		want int
	}{
		{"o2", true, 23},
		{"co2", false, 10},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := getRating(data, tc.ge)
			if got != tc.want {
				t.Errorf("Bad rating: want %d(%b), got %d (%b)", tc.want, tc.want, got, got)
			}
		})
	}
}
