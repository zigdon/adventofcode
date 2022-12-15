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
	want.ParseMap(Point{-8, -10}, strings.Join([]string{
		"..........#..........................",
		".........###.........................",
		"........#####........................",
		".......#######.......................",
		"......#########.............#........",
		".....###########...........###.......",
		"....#############.........#####......",
		"...###############.......#######.....",
		"..#################.....#########....",
		".###################.#.###########...",
		"##########S########################..",
		".###########################S#######.",
		"..###################S#############..",
		"...###################SB##########...",
		"....#############################....",
		".....###########################.....",
		"......#########################......",
		".......#########S#######S#####.......",
		"........#######################......",
		".......#########################.....",
		"......####B######################....",
		".....###S#############.###########...",
		"......#############################..",
		".......#############################.",
		".......#############S#######S########",
		"......B#############################.",
		".....############SB################..",
		"....##################S##########B...",
		"...#######S######################....",
		"....############################.....",
		".....#############S######S######.....",
		"......#########################......",
		".......#######..#############B.......",
		"........#####....###..#######........",
		".........###......#....#####.........",
		"..........#.............###..........",
		".........................#...........",
	}, "\n"))

	if *want.Min != *got.Min || *want.Max != *got.Max {
		t.Errorf("wrong size: want: %s-%s, got %s-%s", want.Min, want.Max, got.Min, got.Max)
	}
	got.Fill()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
		sides := []string{}
		ws := strings.Split(want.String(), "\n")
		for y, l := range strings.Split(got.String(), "\n") {
			sides = append(sides, l+" || "+ws[y])
		}
		t.Logf("\n%s", strings.Join(sides, "\n"))
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
				Min: &Point{-1, -1},
				Max: &Point{1, 1},
				Objs: map[Point]Object{
					{1, 0}:  Air,
					{-1, 0}: Air,
					{0, 1}:  Air,
					{0, -1}: Air,
				},
			},
		},
		{
			p: Point{10, 5},
			r: 2,
			want: &Field{
				Min: &Point{8, 3},
				Max: &Point{12, 7},
				Objs: map[Point]Object{
					{8, 5}:  Air,
					{9, 4}:  Air,
					{9, 5}:  Air,
					{9, 6}:  Air,
					{10, 3}: Air,
					{10, 4}: Air,
					{10, 6}: Air,
					{10, 7}: Air,
					{11, 4}: Air,
					{11, 5}: Air,
					{11, 6}: Air,
					{12, 5}: Air,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d@%s", tc.r, tc.p), func(t *testing.T) {
			f := NewField()
			f.FillAir(tc.p, tc.r)
			if diff := cmp.Diff(tc.want, f); diff != "" {
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
		t.Logf("\n%s", strings.Join(s[10:12], "\n"))
		t.Errorf(diff)
	}
}

func TestTwo(t *testing.T) {
	data := readFile("sample.txt")
	got := two(data)
	want := 0

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}
