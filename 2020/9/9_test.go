package main

import "testing"

func sample() []int {
	return []int{
		35, 20, 15, 25, 47, 40, 62,
		55, 65, 95, 102, 117, 150, 182,
		127, 219, 299, 277, 309, 576}
}

func TestValidateSequence(t *testing.T) {
	got := validateSequence(5, sample())
	if got != 127 {
		t.Errorf("sample validation wrong: want 127, got %d", got)
	}
}

func TestValidNumber(t *testing.T) {
	tests := []struct {
		seq  []int
		n    int
		want bool
	}{
		{[]int{1, 2, 3, 4}, 6, true},
		{[]int{1, 2, 3, 4}, 3, true},
		{[]int{1, 2, 3, 4}, 10, false},
	}

	for i, tc := range tests {
		got := validNumber(tc.n, tc.seq)
		if got != tc.want {
			t.Errorf("bad valid number in test #%d", i)
		}
	}
}
