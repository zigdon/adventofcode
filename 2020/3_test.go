package main

import (
	"fmt"
	"testing"
)

const (
	samplePath = "sample.txt"
)

func TestReadInput(t *testing.T) {
	terrain := &trees{}
	err := terrain.readInput(samplePath)
	if err != nil {
		t.Fatalf("Error reading input: %v", err)
	}

	var trees = map[int]map[int]bool{
		0:  {2: true, 3: true},
		1:  {0: true, 4: true, 8: true},
		2:  {1: true, 6: true, 9: true},
		3:  {2: true, 4: true, 8: true, 10: true},
		4:  {1: true, 5: true, 6: true, 9: true},
		5:  {2: true, 4: true, 5: true},
		6:  {1: true, 3: true, 5: true, 10: true},
		7:  {1: true, 10: true},
		8:  {0: true, 2: true, 3: true, 7: true},
		9:  {0: true, 4: true, 5: true, 10: true},
		10: {1: true, 4: true, 8: true, 10: true},
	}
	for y, line := range terrain.data {
		for x, tree := range line {
			if tree != trees[y][x] {
				t.Errorf("Bad tree at %d, %d: want %v got %v", x, y, trees[y][x], tree)
			}
		}
	}
}

func TestGet(t *testing.T) {
	data := trees{}
	data.init("#.#\n..#\n")
	tests := []struct {
		desc string
		x, y int
		want bool
	}{
		{"simple tree", 0, 0, true},
		{"simple open", 0, 1, false},
		{"past", 1, 2, false},
		{"repeat tree", 3, 0, true},
		{"repeat open", 3, 1, false},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := data.get(tc.x, tc.y)
			if got != tc.want {
				t.Errorf("Tree mismatch at %d,%d: want %v got %v", tc.x, tc.y, tc.want, got)
			}
		})
	}
}

func TestCount(t *testing.T) {
	terrain := &trees{}
	err := terrain.readInput(samplePath)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		right, down, trees int
	}{
		{3, 1, 7},
		{1, 1, 2},
		{5, 1, 3},
		{7, 1, 4},
		{1, 2, 2},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d,%d", tc.right, tc.down), func(t *testing.T) {
			found := terrain.countTrees(tc.right, tc.down)
			if found != tc.trees {
				t.Errorf("Ran into the wrong number of trees: expacted %d, found %d", tc.trees, found)
			}
		})
	}
}
