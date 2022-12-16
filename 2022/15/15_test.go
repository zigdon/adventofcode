package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	got := readFile("sample.txt")
	want := NewField()
	want.Min = &Point{-8, -10}
	want.Max = &Point{28, 26}
	want.Objs = map[Point]Object{
		{-2, 15}: Beacon,
		{0, 11}:  Sensor,
		{2, 0}:   Sensor,
		{2, 10}:  Beacon,
		{2, 18}:  Sensor,
		{8, 7}:   Sensor,
		{9, 16}:  Sensor,
		{10, 16}: Beacon,
		{10, 20}: Sensor,
		{12, 14}: Sensor,
		{13, 2}:  Sensor,
		{14, 3}:  Sensor,
		{14, 17}: Sensor,
		{15, 3}:  Beacon,
		{16, 7}:  Sensor,
		{17, 20}: Sensor,
		{20, 1}:  Sensor,
		{20, 14}: Sensor,
		{21, 22}: Beacon,
		{25, 17}: Beacon,
	}
	want.Gaps = map[int][]*Range{
		-10: {{2, 2}},
		-9:  {{1, 3}},
		-8:  {{0, 4}},
		-7:  {{-1, 5}},
		-6:  {{-2, 6}, {20, 20}},
		-5:  {{-3, 7}, {19, 21}},
		-4:  {{-4, 8}, {18, 22}},
		-3:  {{-5, 9}, {17, 23}},
		-2:  {{-6, 10}, {16, 24}},
		-1:  {{-7, 11}, {13, 13}, {15, 25}},
		0:   {{-8, 26}},
		1:   {{-7, 27}},
		2:   {{-6, 26}},
		3:   {{-5, 25}},
		4:   {{-4, 24}},
		5:   {{-3, 23}},
		6:   {{-2, 22}},
		7:   {{-1, 21}},
		8:   {{0, 22}},
		9:   {{-1, 23}},
		10:  {{-2, 24}},
		11:  {{-3, 13}, {15, 25}},
		12:  {{-2, 26}},
		13:  {{-1, 27}},
		14:  {{-1, 28}},
		15:  {{-2, 27}},
		16:  {{-3, 26}},
		17:  {{-4, 25}},
		18:  {{-5, 24}},
		19:  {{-4, 23}},
		20:  {{-3, 23}},
		21:  {{-2, 22}},
		22:  {{-1, 5}, {8, 21}},
		23:  {{0, 4}, {9, 11}, {14, 20}},
		24:  {{1, 3}, {10, 10}, {15, 19}},
		25:  {{2, 2}, {16, 18}},
		26:  {{17, 17}},
	}

	if *want.Min != *got.Min || *want.Max != *got.Max {
		t.Errorf("wrong size: want: %s-%s, got %s-%s", want.Min, want.Max, got.Min, got.Max)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
		t.Logf("\n%s", got)
	}
}

func TestFill(t *testing.T) {
	tests := []struct {
		p    Point
		r    int
		want *Field
	}{
		{
			p: Point{0, 0},
			r: 1,
			want: &Field{
				Objs: map[Point]Object{},
				Min:  &Point{-1, -1},
				Max:  &Point{1, 1},
				Gaps: map[int][]*Range{
					-1: {{}},
					1:  {{}},
					0:  {{-1, 1}},
				},
			},
		},
		{
			p: Point{10, 5},
			r: 2,
			want: &Field{
				Objs: map[Point]Object{},
				Min:  &Point{8, 3},
				Max:  &Point{12, 7},
				Gaps: map[int][]*Range{
					3: {{10, 10}},
					4: {{9, 11}},
					5: {{8, 12}},
					6: {{9, 11}},
					7: {{10, 10}},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d@%s", tc.r, tc.p), func(t *testing.T) {
			f := NewField()
			f.FillSparseAir(tc.p, tc.r)
			if diff := cmp.Diff(tc.want, f); diff != "" {
				t.Errorf(diff)
			}
		})
	}

}

func TestMergeRanges(t *testing.T) {
	tests := []struct {
		in, out []*Range
	}{
		{
			in:  []*Range{{1, 5}, {2, 3}, {7, 8}},
			out: []*Range{{1, 5}, {7, 8}},
		},
		{
			in: []*Range{
				{-1, 28}, {0, 5}, {8, 16}, {8, 16}, {12, 16}, {6, 10},
				{0, 0}, {12, 28}, {12, 28}, {17, 17}},
			out: []*Range{{-1, 28}},
		},
		{
			in:  []*Range{{13, 13}, {7, 9}, {-7, 11}, {15, 25}},
			out: []*Range{{-7, 11}, {13, 13}, {15, 25}},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			got := MergeRanges(tc.in)
			if diff := cmp.Diff(tc.out, got); diff != "" {
				t.Log(tc.in)
				t.Log(got)
				t.Errorf(diff)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		a, b, out *Range
		want      bool
	}{
		{
			a:    NewRange(7, 9),
			b:    NewRange(-7, 11),
			out:  NewRange(-7, 11),
			want: true,
		},
		{
			a:    NewRange(1, 5),
			b:    NewRange(3, 5),
			out:  NewRange(1, 5),
			want: true,
		},
		{
			a:    NewRange(1, 2),
			b:    NewRange(3, 5),
			out:  NewRange(1, 5),
			want: true,
		},
		{
			a:    NewRange(1, 2),
			b:    NewRange(4, 5),
			out:  NewRange(1, 2),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s+%s", tc.a, tc.b), func(t *testing.T) {
			got := tc.a.Merge(tc.b)
			if got != tc.want {
				t.Errorf("bad merge")
			}
			if diff := cmp.Diff(tc.out, tc.a); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestOne(t *testing.T) {
	data := readFile("sample.txt")
	got := one(data, 10)
	want := 26

	if diff := cmp.Diff(want, got); diff != "" {
		s := strings.Split(data.String(), "\n")
		t.Logf("\n%s", strings.Join(s[10:13], "\n"))
		t.Errorf(diff)
	}
}

func TestTwo(t *testing.T) {
	data := readFile("sample.txt")
	t.Logf("\n%s", data)
	got := two(data, 20)
	want := 56000011

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}
