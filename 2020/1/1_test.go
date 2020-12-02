package main

import (
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sample() []int {
	return []int{1721, 979, 366, 299, 675, 1456}
}

func TestParseFile(t *testing.T) {
	data := strings.Join([]string{"1", "2", "3", ""}, "\n")
	got, err := parseFile(data)
	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}
	if diff := cmp.Diff([]int{1, 2, 3}, got); diff != "" {
		t.Errorf("Bad at reading: %s", diff)
	}
}

func TestPartA(t *testing.T) {
	a, b, err := findPair(2020, sample())
	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}

	res := []int{a, b}
	sort.Ints(res)
	if diff := cmp.Diff([]int{299, 1721}, res); diff != "" {
		t.Errorf("Wrong pair found: %s", diff)
	}
}

func TestPartB(t *testing.T) {
	a, b, c, err := findTriplet(2020, sample())

	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}

	res := []int{a, b, c}
	sort.Ints(res)
	if diff := cmp.Diff([]int{366, 675, 979}, res); diff != "" {
		t.Errorf("Wrong set found: %d, %d, %d", a, b, c)
	}
}
