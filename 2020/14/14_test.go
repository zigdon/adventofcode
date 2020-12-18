package main

import (
	"fmt"
	"testing"
)

func sample() []string {
	return []string{
		"mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X",
		"mem[8] = 11",
		"mem[7] = 101",
		"mem[8] = 0    ",
	}
}

func setMask(t *testing.T) {
	tests := []struct {
		input   string
		wantOn  string
		wantOff string
	}{
		{
			input:   "XXX100X",
			wantOn:  "1000",
			wantOff: "1111001",
		},
	}

	c := newCompy()
	for _, tc := range tests {
		c.setMask(tc.input)
		got := fmt.Sprintf("%b", c.maskOn)
		if got != tc.wantOn {
			t.Errorf("bad on mask\nwant: %s\n got: %s", tc.wantOn, got)
		}
		got = fmt.Sprintf("%b", c.maskOff)
		if got != tc.wantOff {
			t.Errorf("bad off mask\nwant: %s\n got: %s", tc.wantOff, got)
		}
	}
}

func TestRunLines(t *testing.T) {
	c := newCompy()
	c.runLines(sample())
	got := c.sumMem()
	if got != 165 {
		t.Errorf("bad memory sum: want 165, got %d", got)
	}
}
