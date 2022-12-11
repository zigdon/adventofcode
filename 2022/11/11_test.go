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

func TestOne(t *testing.T) {
	data, err := readFile("sample.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	got := one(data)
	want := 0

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
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
