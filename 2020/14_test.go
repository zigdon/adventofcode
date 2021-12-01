package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sample() []string {
	return []string{
		"mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X",
		"mem[8] = 11",
		"mem[7] = 101",
		"mem[8] = 0    ",
	}
}

func sample2() []string {
	return []string{
		"mask = 000000000000000000000000000000X1001X",
		"mem[42] = 100",
		"mask = 00000000000000000000000000000000X0XX",
		"mem[26] = 1    ",
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

func TestSetMaskedAddr(t *testing.T) {
	tests := []struct {
		addr int
		mask string
		want []int
	}{
		{
			addr: 42,
			mask: "000000000000000000000000000000X1101X",
			want: []int{26, 27, 58, 59},
		},
		{
			addr: 26,
			mask: "00000000000000000000000000000000X0XX",
			want: []int{16, 17, 18, 19, 24, 25, 26, 27},
		},
	}

	c := newCompy()
	for i, tc := range tests {
		c.runLines([]string{"mask = " + tc.mask}, true)
		got := c.setMaskedAddr(tc.addr, 0)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("wrong addresses set for test #%d:\n%s", i, diff)
		}
	}
}

func TestRunLines(t *testing.T) {
	c := newCompy()
	c.runLines(sample(), false)
	got := c.sumMem()
	if got != 165 {
		t.Errorf("bad memory sum: want 165, got %d", got)
	}
}
