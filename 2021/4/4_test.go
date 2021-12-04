package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadInput(t *testing.T) {
	calls, boards := readInput("sample.txt")

	want := []int{
		7, 4, 9, 5, 11, 17, 23, 2, 0, 14, 21, 24, 10, 16, 13, 6, 15, 25, 12,
		22, 18, 20, 8, 19, 3, 26, 1}
	if diff := cmp.Diff(want, calls); diff != "" {
		t.Errorf("Bad calls:\n%s", diff)
	}

	if len(boards) != 3 {
		t.Errorf("Wrong number of boards, want 3, got %d", len(boards))
	}

	wantBoard := [][]int{
		{22, 13, 17, 11, 0},
		{8, 2, 23, 4, 24},
		{21, 9, 14, 16, 7},
		{6, 10, 3, 18, 5},
		{1, 12, 20, 15, 19},
	}

	if diff := cmp.Diff(wantBoard, boards[0].lines); diff != "" {
		e := ""
		for y, l := range boards[0].lines {
			e += fmt.Sprintf("%v  %v\n", wantBoard[y], l)
		}
		t.Errorf("Bad board:\n%s\n%s", e, diff)
	}
}

func TestPlay(t *testing.T) {
	calls, boards := readInput("sample.txt")

	last := -1
	for _, call := range calls {
		t.Logf("Playing %d", call)
		for i, board := range boards {
			if board.finished {
				continue
			}
			win := board.play(call)
			if win {
				t.Logf("Board %d finished", i)
				if !board.finished {
					t.Errorf("board %d not marked as finished!", i)
				}
				if last == -1 {
					if i != 2 {
						t.Errorf("Wrong board won")
					}

					sum := board.sum()
					if sum != 188 {
						t.Errorf("Bad sum, want 188, got %d", sum)
					}
				}
				last = i
			}
		}
	}

	if last != 1 {
		t.Errorf("Wrong board won last, expected 1, got %d", last)
	}
}
