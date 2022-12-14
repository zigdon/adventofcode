package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func mkE() *Packet {
	return &Packet{Empty: true}
}

func mk(n ...interface{}) *Packet {
	p := &Packet{}
	for _, v := range n {
		if i, ok := v.(int); ok {
			p.Sub = append(p.Sub, &Packet{Literal: fmt.Sprintf("%d", i), Val: i})
		} else {
			p.Sub = append(p.Sub, v.(*Packet))
		}
	}
	return p
}

func TestReadFile(t *testing.T) {
	got := readFile("sample.txt")
	want := []Pair{
		{Left: mk(1, 1, 3, 1, 1), Right: mk(1, 1, 5, 1, 1)},
		{Left: mk(mk(1), mk(2, 3, 4)), Right: mk(mk(1), 4)},
		{Left: mk(9), Right: mk(mk(8, 7, 6))},
		{Left: mk(mk(4, 4), 4, 4), Right: mk(mk(4, 4), 4, 4, 4)},
		{Left: mk(7, 7, 7, 7), Right: mk(7, 7, 7)},
		{Left: mkE(), Right: mk(3)},
		{Left: mk(mk(mkE())), Right: mk(mkE())},
		{Left: mk(1, mk(2, mk(3, mk(4, mk(5, 6, 7)))), 8, 9),
			Right: mk(1, mk(2, mk(3, mk(4, mk(5, 6, 0)))), 8, 9)},
	}

	if diff := cmp.Diff(want, got,
		cmpopts.IgnoreFields(Pair{}, "Original"),
		cmpopts.IgnoreFields(Packet{}, "Literal")); diff != "" {
		t.Error(diff)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		in    string
		want  *Packet
		wantS string
	}{
		{
			in:    "[[8,[],[10]],[]]",
			want:  mk(mk(8, mkE(), mk(10)), mkE()),
			wantS: "[[8,E,[10]],E]",
		},
		{
			in:    "[[[[8,5,6,6,5],1,[10]],[]],[],[],[7],[2,2]]",
			want:  mk(mk(mk(mk(8, 5, 6, 6, 5), 1, mk(10)), mkE()), mkE(), mkE(), mk(7), mk(2, 2)),
			wantS: "[[[[8,5,6,6,5],1,[10]],E],E,E,[7],[2,2]]",
		},
		{
			in:    "[1,[],5]",
			want:  mk(1, mkE(), 5),
			wantS: "[1,E,5]",
		},
	}

	for _, tc := range tests {
		got := Pair{}.Parse(tc.in)
		if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(Packet{}, "Literal")); diff != "" {
			t.Errorf(diff)
		}
		if diff := cmp.Diff(tc.wantS, got.String()); diff != "" {
			t.Errorf(diff)
		}
	}

}

func TestLt(t *testing.T) {
	got := readFile("sample.txt")
	want := []bool{
		true, true, false, true, false, true, false, false,
	}

	for i, p := range got {
		t.Run(fmt.Sprintf("%s", p), func(t *testing.T) {
			if p.Ordered() != want[i] {
				t.Errorf("Pair #%d misordered: want %v, got %v", i, want[i], !want[i])
			}
		})
	}
}

func TestExtraOrdered(t *testing.T) {
	tests := []struct {
		inL, inR    string
		wantOrdered bool
	}{
		{"[[[9]]]", "[5]", false},
		{
			"[9,10,[1,0,6],[1,2]]",
			"[15]",
			true,
		},
		{
			"[9,10,[1,0,6],[1,2]]",
			"[5]",
			false,
		},
		{
			"[[9,10,[1,0,6],[1,2]],[0,[9,5,3],[7,3,10,1],0,0],5,10,0]",
			"[5,6,1,6]",
			false,
		},
		{
			"[[[9,10,[1,0,6],[1,2]],[0,[9,5,3],[7,3,10,1],0,0],5,10,0],0]",
			"[[5,6,1,6],[8,6,5,[7,0,5,[0,4,3,6,6],[8,7,8,5]]],[[[1,2],9,1],10,4,3,[5,[5,10,9,9],7,0,[2,3]]],[[[1,5],[0],0],3,[0]]]",
			false,
		},
		{
			"[[0,[[4,6,4,8]],3,[10]],0,[6,8,6],[8,[3,[7,7,10,6,5]],[0,8,[7],0,5]],[1,[3,[4,4,8,3,8],3,[1,9,8,3,10],[0,5,8]]]]",
			"[[0,0,10]]",
			false,
		},
		{"[[[]]]", "[[]]", false},
		{"[[]]", "[[[]]]", true},
        {"[[[]]]", "[0]", true},
        {"[0]", "[[[]]]", false},
        {"[9]", "[[9,9]]", true},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s / %s", tc.inL, tc.inR), func(t *testing.T) {
			p := NewPair(tc.inL, tc.inR)
			got := p.Ordered()
			if got != tc.wantOrdered {
				t.Errorf("Bad order, wanted %v, got %v", tc.wantOrdered, got)
			}
		})
	}
}

func TestOne(t *testing.T) {
	data := readFile("sample.txt")
	got := one(data)
	want := 13

	if diff := cmp.Diff(want, got); diff != "" {
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
