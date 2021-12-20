package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMakeSnail(t *testing.T) {
	tests := []struct {
		in       string
		want     *snail
		wantStr  string
		leftover string
	}{
		{
			in:      "[4,5]",
			want:    &snail{LeftVal: 4, RightVal: 5},
			wantStr: "[4,5]",
		},
		{
			in: "[[4,5],[6,7]]",
			want: &snail{
				LeftSnail:  &snail{LeftVal: 4, RightVal: 5},
				RightSnail: &snail{LeftVal: 6, RightVal: 7}},
			wantStr: "[[4,5],[6,7]]",
		},
		{
			in:       "[4,5][6,7]",
			want:     &snail{LeftVal: 4, RightVal: 5},
			leftover: "[6,7]",
			wantStr:  "[4,5]",
		},
		{
			in: "[[6,7],0]",
			want: &snail{
				LeftSnail: &snail{LeftVal: 6, RightVal: 7},
				RightVal:  0},
			wantStr: "[[6,7],0]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, leftover := makeSnail(tc.in, nil)
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(snail{}, "Parent")); diff != "" {
				t.Errorf("misbehaving snail:\n%s", diff)
			}
			if got.String() != tc.wantStr {
				t.Errorf("bad string rep: want %q, got %q", tc.wantStr, got.String())
			}
			if leftover != tc.leftover {
				t.Errorf("bad leftover, want %q, got %q", tc.leftover, leftover)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	ss := readFile("sample.txt")
	want := []*snail{
		{
			LeftSnail: &snail{
				LeftSnail: &snail{
					LeftVal:    0,
					RightSnail: &snail{LeftVal: 4, RightVal: 5}},
				RightSnail: &snail{}},
			RightSnail: &snail{
				LeftSnail: &snail{
					LeftSnail:  &snail{LeftVal: 4, RightVal: 5},
					RightSnail: &snail{LeftVal: 2, RightVal: 6}},
				RightSnail: &snail{LeftVal: 9, RightVal: 5},
			},
		},
		{
			LeftVal: 7,
			RightSnail: &snail{
				LeftSnail: &snail{
					LeftSnail:  &snail{LeftVal: 3, RightVal: 7},
					RightSnail: &snail{LeftVal: 4, RightVal: 3}},
				RightSnail: &snail{
					LeftSnail:  &snail{LeftVal: 6, RightVal: 3},
					RightSnail: &snail{LeftVal: 8, RightVal: 8}},
			},
		},
		{
			LeftSnail: &snail{
				LeftVal: 2,
				RightSnail: &snail{
					LeftSnail:  &snail{RightVal: 8},
					RightSnail: &snail{LeftVal: 3, RightVal: 4}},
			},
			RightSnail: &snail{
				LeftSnail:  &snail{RightVal: 1, LeftSnail: &snail{LeftVal: 6, RightVal: 7}},
				RightSnail: &snail{LeftVal: 7, RightSnail: &snail{LeftVal: 1, RightVal: 6}},
			},
		},
		{
			LeftSnail: &snail{
				LeftSnail:  &snail{RightVal: 7, LeftSnail: &snail{LeftVal: 2, RightVal: 4}},
				RightSnail: &snail{LeftVal: 6, RightSnail: &snail{RightVal: 5}},
			},
			RightSnail: &snail{
				LeftSnail: &snail{
					LeftSnail:  &snail{LeftVal: 6, RightVal: 8},
					RightSnail: &snail{LeftVal: 2, RightVal: 8}},
				RightSnail: &snail{
					LeftSnail:  &snail{LeftVal: 2, RightVal: 1},
					RightSnail: &snail{LeftVal: 4, RightVal: 5}},
			},
		},
		{
			LeftVal: 7,
			RightSnail: &snail{
				LeftVal: 5,
				RightSnail: &snail{
					LeftSnail:  &snail{LeftVal: 3, RightVal: 8},
					RightSnail: &snail{LeftVal: 1, RightVal: 4}},
			},
		},
		{
			LeftSnail:  &snail{LeftVal: 2, RightSnail: &snail{LeftVal: 2, RightVal: 2}},
			RightSnail: &snail{LeftVal: 8, RightSnail: &snail{LeftVal: 8, RightVal: 1}},
		},
		{LeftVal: 2, RightVal: 9},
		{
			LeftVal: 1,
			RightSnail: &snail{
				LeftSnail:  &snail{RightVal: 9, LeftSnail: &snail{LeftVal: 9, RightVal: 3}},
				RightSnail: &snail{LeftSnail: &snail{LeftVal: 9}, RightSnail: &snail{RightVal: 7}},
			},
		},
		{
			RightVal: 1,
			LeftSnail: &snail{
				RightVal:  7,
				LeftSnail: &snail{LeftVal: 5, RightSnail: &snail{LeftVal: 7, RightVal: 4}}},
		},
		{
			LeftSnail: &snail{
				RightVal:  6,
				LeftSnail: &snail{RightVal: 2, LeftSnail: &snail{LeftVal: 4, RightVal: 2}}},
			RightSnail: &snail{LeftVal: 8, RightVal: 7},
		},
	}

	for i, s := range ss {
		if diff := cmp.Diff(want[i], s, cmpopts.IgnoreFields(snail{}, "Parent")); diff != "" {
			t.Errorf("bad #%d snails!\n%s", i, diff)
		}
	}
}

func TestAddAdjecent(t *testing.T) {
	init := func() []*snail {
		p, _ := makeSnail("[[1,2],[3,[4,5]]]", nil)
		s1 := p.LeftSnail
		s2 := p.RightSnail
		s3 := p.RightSnail.RightSnail
		return []*snail{p, s1, s2, s3}
	}
	init()
	tests := []struct {
		in   int
		d    direction
		want string
	}{
		{
			in:   1,
			d:    LEFT,
			want: "[[1,2],[3,[4,5]]]",
		},
		{
			in:   1,
			d:    RIGHT,
			want: "[[1,2],[5,[4,5]]]",
		},
		{
			in:   2,
			d:    LEFT,
			want: "[[1,5],[3,[4,5]]]",
		},
		{
			in:   3,
			d:    LEFT,
			want: "[[1,2],[7,[4,5]]]",
		},
		{
			in:   3,
			d:    RIGHT,
			want: "[[1,2],[3,[4,5]]]",
		},
	}

	for _, tc := range tests {
		ss := init()
		t.Run(fmt.Sprintf("%s (%s): %s", ss[0], ss[tc.in], tc.d), func(t *testing.T) {
			ss[tc.in].addAdjecent(tc.d)
			want, _ := makeSnail(tc.want, nil)

			if diff := cmp.Diff(want, ss[0]); diff != "" {
				t.Errorf("bad addAdj(%s):\n%s", tc.d, diff)
			}
		})
	}
}

func TestExplode(t *testing.T) {
	tests := []struct {
		in, out string
		want    bool
	}{
		{
			in:  "[1,2]",
			out: "[1,2]",
		},
		{
			in:   "[[[[[1,2],3],4],5],6]",
			out:  "[[[[0,5],4],5],6]",
			want: true,
		},
		{
			in:   "[[[[1,3],[[1,2],4]],5],6]",
			out:  "[[[[1,4],[0,6]],5],6]",
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			s, _ := makeSnail(tc.in, nil)
			got := s.explode()
			if s.String() != tc.out {
				t.Errorf("bad explosion: want %s, got %s", tc.out, s)
			}
			if got != tc.want {
				t.Errorf("bad explosion status: want %v, got %v", tc.want, got)
			}

		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		ns   []string
		want string
	}{
		{
			ns:   []string{"[[[[4,3],4],4],[7,[[8,4],9]]]", "[1,1]"},
			want: "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
		},
		{
			ns: []string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
			},
			want: "[[[[1,1],[2,2]],[3,3]],[4,4]]",
		},
		{
			ns: []string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
				"[5,5]",
			},
			want: "[[[[3,0],[5,3]],[4,4]],[5,5]]",
		},
		{
			ns: []string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
				"[5,5]",
				"[6,6]",
			},
			want: "[[[[5,0],[7,4]],[5,5]],[6,6]]",
		},
		{
			ns: []string{
				"[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]",
				"[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]",
				"[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]",
				"[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]",
				"[7,[5,[[3,8],[1,4]]]]",
				"[[2,[2,2]],[8,[8,1]]]",
				"[2,9]",
				"[1,[[[9,3],9],[[9,0],[0,7]]]]",
				"[[[5,[7,4]],7],1]",
				"[[[[4,2],2],6],[8,7]]",
			},
			want: "[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]",
		},
		{
			ns: []string{
				"[[[0,[5,8]],[[1,7],[9,6]]],[[4,[1,2]],[[1,4],2]]]",
				"[[[5,[2,8]],4],[5,[[9,9],0]]]",
				"[6,[[[6,2],[5,6]],[[7,6],[4,7]]]]",
				"[[[6,[0,7]],[0,9]],[4,[9,[9,0]]]]",
				"[[[7,[6,4]],[3,[1,3]]],[[[5,5],1],9]]",
				"[[6,[[7,3],[3,2]]],[[[3,8],[5,7]],4]]",
				"[[[[5,4],[7,7]],8],[[8,3],8]]",
				"[[9,3],[[9,9],[6,[4,9]]]]",
				"[[2,[[7,7],7]],[[5,8],[[9,3],[0,2]]]]",
				"[[[[5,2],5],[8,[3,7]]],[[5,[7,5]],[4,4]]]",
			},
			want: "[[[[6,6],[7,6]],[[7,7],[7,0]]],[[[7,7],[7,7]],[[7,8],[9,9]]]]",
		},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			t.Log(tc.ns[0])
			a, _ := makeSnail(tc.ns[0], nil)
			for _, sb := range tc.ns[1:] {
				t.Log(sb)
				b, _ := makeSnail(sb, nil)
				a = a.add(b)
			}

			if a.String() != tc.want {
				t.Errorf("bad addition:\nwant: %s\n got: %s", tc.want, a)
			}
		})
	}
}

func TestMagnitude(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{
			in:   "[9,1]",
			want: 29,
		},
		{
			in:   "[1,9]",
			want: 21,
		},
		{
			in:   "[[9,1],[1,9]]",
			want: 129,
		},
		{
			in:   "[[1,2],[[3,4],5]]",
			want: 143,
		},
		{
			in:   "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
			want: 1384,
		},
		{
			in:   "[[[[1,1],[2,2]],[3,3]],[4,4]]",
			want: 445,
		},
		{
			in:   "[[[[3,0],[5,3]],[4,4]],[5,5]]",
			want: 791,
		},
		{
			in:   "[[[[5,0],[7,4]],[5,5]],[6,6]]",
			want: 1137,
		},
		{
			in:   "[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]",
			want: 3488,
		},
		{
			in:   "[[[[6,6],[7,6]],[[7,7],[7,0]]],[[[7,7],[7,7]],[[7,8],[9,9]]]]",
			want: 4140,
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			s, _ := makeSnail(tc.in, nil)
			got := s.magnitude()
			if got != tc.want {
				t.Errorf("bad magnitude: want %d, got %d", tc.want, got)
			}
		})
	}
}

func TestE2E(t *testing.T) {
	ss := readFile("sample2.txt")
	s := ss[0]
	t.Logf("  %s", s)
	for _, sb := range ss[1:] {
		t.Logf("+ %s", sb)
		s = s.add(sb)
		t.Logf("= %s", s)
	}
	t.Logf("magnitude: %d", s.magnitude())
	if s.magnitude() != 4140 {
		t.Errorf("bad e2e magnitude: want 4140, got %d", s.magnitude())
	}
}
