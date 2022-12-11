package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	got, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	want := []struct {
		id         int
		items      []int
		test, T, F int
		opOut      int
	}{
		{
			id:    0,
			items: []int{79, 98},
			test:  23,
			T:     2,
			F:     3,
			opOut: 95,
		},
		{
			id:    1,
			items: []int{54, 65, 75, 74},
			test:  19,
			T:     2,
			F:     0,
			opOut: 11,
		},
		{
			id:    2,
			items: []int{79, 60, 97},
			test:  13,
			T:     1,
			F:     3,
			opOut: 25,
		},
		{
			id:    3,
			items: []int{74},
			test:  17,
			T:     0,
			F:     1,
			opOut: 8,
		},
	}

	for i, m := range got {
		t.Run(fmt.Sprintf("test #%d", i), func(t *testing.T) {
			w := want[i]
			if m.ID != w.id {
				t.Errorf("Bad ID: want %d, got %d", m.ID, w.id)
			}
			if diff := cmp.Diff(w.items, m.Items); diff != "" {
				t.Errorf("Bad items:\n%s", diff)
			}
			if m.Test != w.test {
				t.Errorf("Bad test: want %d, got %d", w.test, m.Test)
			}
			if m.True != w.T {
				t.Errorf("Bad true: want %d, got %d", w.T, m.True)
			}
			if m.False != w.F {
				t.Errorf("Bad false: want %d, got %d", w.F, m.False)
			}
			op := m.Op(5)
			if op != w.opOut {
				t.Errorf("Bad op: want %d, got %d", w.opOut, op)
			}
		})
	}
}

func TestTurn(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	tr := NewTroop(data, 3)
	tr.Turn(0)

	if !tr.Monkeys[3].Has(500) {
		t.Errorf("M3 doesn't have 500: %v", tr.Monkeys[3].Items)
	}
	if !tr.Monkeys[3].Has(620) {
		t.Errorf("M3 doesn't have 620: %v", tr.Monkeys[3].Items)
	}

}

func TestRound(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	tr := NewTroop(data, 3)

	tests := [][][]int{
		{{20, 23, 27, 26}, {2080, 25, 167, 207, 401, 1046}, {}, {}},
		{{695, 10, 71, 135, 350}, {43, 49, 58, 55, 362}, {}, {}},
		{{16, 18, 21, 20, 122}, {1468, 22, 150, 286, 739}, {}, {}},
		{{491, 9, 52, 97, 248, 34}, {39, 45, 43, 258}, {}, {}},
		{{15, 17, 16, 88, 1037}, {20, 110, 205, 524, 72}, {}, {}},
		{{8, 70, 176, 26, 34}, {481, 32, 36, 186, 2190}, {}, {}},
		{{162, 12, 14, 64, 732, 17}, {148, 372, 55, 72}, {}, {}},
		{{51, 126, 20, 26, 136}, {343, 26, 30, 1546, 36}, {}, {}},
		{{116, 10, 12, 517, 14}, {108, 267, 43, 55, 288}, {}, {}},
		{{91, 16, 20, 98}, {481, 245, 22, 26, 1092, 30}, {}, {}}, // 10
		{},
		{},
		{},
		{},
		{{83, 44, 8, 184, 9, 20, 26, 102}, {110, 36}, {}, {}}, // 15
		{},
		{},
		{},
		{},
		{{10, 12, 14, 26, 34}, {245, 93, 53, 199, 115}, {}, {}}, //20
	}

	for n, tc := range tests {
		tr.Round()
		if tc == nil {
			t.Logf("Skipping round #%d", n)
			continue
		}
		t.Logf("Playing round #%d", n)
		for id, items := range tc {
			if diff := cmp.Diff(items, tr.Monkeys[id].Items); diff != "" {
				t.Errorf("Round #%d: Wrong items for M%d: want %v, got %v",
					n, id, items, tr.Monkeys[id].Items)
			}
		}
	}

}

func TestOne(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

    got := one(data)
	want := 10605

	if got != want {
		t.Errorf("bad MB: want %d, got %d", want, got)
	}
}

func TestTwo(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := two(data)
	want := 0

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}
