package main

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sample1() []int {
	return []int{16, 10, 15, 5, 1, 11, 7, 19, 6, 12, 4}
}

func sample2() []int {
	return []int{28, 33, 18, 42, 31, 14, 46, 20, 48, 47, 24, 23, 49,
		45, 19, 38, 39, 11, 1, 32, 25, 35, 8, 17, 7, 9, 4, 2, 34, 10, 3}
}

func TestSectionAdapters(t *testing.T) {
	tests := []struct {
		data []int
		want [][]int
	}{
		{
			data: sample1(),
			want: [][]int{
				{4, 5, 6, 7},
				{10, 11, 12},
				{15, 16},
			},
		},
		{
			data: sample2(),
			want: [][]int{
				{1, 2, 3, 4},
				{7, 8, 9, 10, 11},
				{17, 18, 19, 20},
				{23, 24, 25},
				{31, 32, 33, 34, 35},
				{38, 39},
				{45, 46, 47, 48, 49},
			},
		},
	}

	for _, tc := range tests {
		got := sectionAdapters(tc.data)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("Bad sections:\n%s", diff)
		}
	}
}

func TestCountChains(t *testing.T) {
	tests := []struct {
		data []int
		want int
	}{
		{data: sample1(), want: 8},
		{data: sample2(), want: 19208},
	}

	for _, tc := range tests {
		tc.data = append(tc.data, 0)
		sort.Ints(tc.data)
		got := countChains(tc.data)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("bad count: %s", diff)
		}
	}
}

func TestBuildChains(t *testing.T) {
	tests := []struct {
		data []int
		want int
	}{
		{data: sample1(), want: 8},
		{data: sample2(), want: 19208},
	}

	for _, tc := range tests {
		tc.data = append(tc.data, 0)
		sort.Ints(tc.data)
		got := buildChain(tc.data)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("bad build: %s", diff)
		}
	}
}

func TestGetDist(t *testing.T) {
	tests := []struct {
		data []int
		want map[int]int
	}{
		{data: sample1(), want: map[int]int{1: 7, 3: 5}},
		{data: sample2(), want: map[int]int{1: 22, 3: 10}},
	}

	for _, tc := range tests {
		got := getDist(tc.data)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("bad dist: %s", diff)
		}
	}
}
