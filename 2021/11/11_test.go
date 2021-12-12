package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	got := readFile("sample2.txt")
	want := newDumbo()
	want.init(0, "11111", "19991", "19191", "19991", "11111")
	want.init(1, "34543", "40004", "50005", "40004", "34543")
	want.init(2, "45654", "51115", "61116", "51115", "45654")

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad reading:\n%s", diff)
	}
}

func TestFlash(t *testing.T) {
	want := newDumbo()
	want.init(0, "11111", "19991", "19191", "19991", "11111")
	want.init(1, "34543", "40004", "50005", "40004", "34543")
	want.init(2, "45654", "51115", "61116", "51115", "45654")

	data := readFile("sample2.txt")
	got := data.step(1)
	if got != 9 {
		t.Errorf("wrong number of flashes at step 1: want 9, got %d", got)
	}
	if diff := cmp.Diff(want.Layout[1], data.Layout[1]); diff != "" {
		t.Logf("want:\n%s", want.string(1))
		t.Logf("got:\n%s", data.string(1))
		t.Errorf("bad board after step 1:\n%s", diff)
	}
	got = data.step(2) - got
	if got != 0 {
		t.Errorf("wrong number of flashes at step 2: want 0, got %d", got)
	}
	if diff := cmp.Diff(want.Layout[2], data.Layout[2]); diff != "" {
		t.Logf("want:\n%s", want.string(2))
		t.Logf("got:\n%s", data.string(2))
		t.Errorf("bad board after step 2:\n%s", diff)
	}

}

func TestMissingSteps(t *testing.T) {
	data := readFile("sample.txt")
	d := newDumbo()
	in := []string{}
	for _, l := range *data.Layout[0] {
		line := ""
		for _, n := range l {
			line += fmt.Sprintf("%d", n)
		}
		in = append(in, line)
	}
	d.init(0, in...)

	want := map[int]int{
		10:  204,
		100: 1656,
	}

	for _, it := range []int{10, 100} {
		if it == 0 {
			continue
		}
		got := d.step(it)
		w, ok := want[it]
		if ok && got != w {
			t.Errorf("bad number of flashes at %d: want %d, got %d", it, w, got)
		}
		if diff := cmp.Diff(data.Layout[it], d.Layout[it]); diff != "" {
			t.Logf("want:\n%s", data.string(it))
			t.Logf("got:\n%s", d.string(it))
			t.Errorf("bad board after step %d:\n%s", it, diff)
		}

	}
}

func TestSync(t *testing.T) {
	data := readFile("sample.txt")
	if !data.isSynced(195) {
		t.Errorf("step 195 was supposed to be in sync!")
	}
}
