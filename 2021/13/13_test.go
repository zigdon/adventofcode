package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	layout, insts := readFile("sample.txt")
	wantLayout := &paper{
		Marks: map[coord]bool{
			{6, 10}:  true,
			{0, 14}:  true,
			{9, 10}:  true,
			{0, 3}:   true,
			{10, 4}:  true,
			{4, 11}:  true,
			{6, 0}:   true,
			{6, 12}:  true,
			{4, 1}:   true,
			{0, 13}:  true,
			{10, 12}: true,
			{3, 4}:   true,
			{3, 0}:   true,
			{8, 4}:   true,
			{1, 10}:  true,
			{2, 14}:  true,
			{8, 10}:  true,
			{9, 0}:   true,
		},
		Size: coord{10, 14},
	}
	wantInst := []inst{
		{"y", 7},
		{"x", 5},
	}

	if diff := cmp.Diff(wantLayout, layout); diff != "" {
		t.Errorf("bad layout:\n%s", diff)
	}
	if diff := cmp.Diff(wantInst, insts); diff != "" {
		t.Errorf("bad instructions:\n%s", diff)
	}
}

func TestFold(t *testing.T) {
	tests := []struct {
		input *paper
		folds []inst
		want  *paper
	}{
		{
			input: &paper{
				Marks: map[coord]bool{
					{0, 0}: true,
					{2, 0}: true,
					{3, 4}: true,
					{3, 3}: true,
				},
				Size: coord{3, 4},
			},
			folds: []inst{
				{"y", 2},
			},
			want: &paper{
				Marks: map[coord]bool{
					{0, 0}: true,
					{2, 0}: true,
					{3, 0}: true,
					{3, 1}: true,
				},
				Size: coord{3, 2},
			},
		},
		{
			input: &paper{
				Marks: map[coord]bool{
					{0, 0}: true,
					{3, 0}: true,
					{3, 4}: true,
					{3, 3}: true,
				},
				Size: coord{3, 4},
			},
			folds: []inst{
				{"y", 2},
				{"x", 2},
			},
			want: &paper{
				Marks: map[coord]bool{
					{0, 0}: true,
					{1, 0}: true,
					{1, 1}: true,
				},
				Size: coord{2, 2},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%#v", tc.folds), func(t *testing.T) {
			tc.input.fold(tc.folds)
			if diff := cmp.Diff(tc.want, tc.input); diff != "" {
				t.Errorf("bad folding:\n%s", diff)
			}
		})
	}
}
