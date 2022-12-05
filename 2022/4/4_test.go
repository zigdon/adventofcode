package main

import (
    "testing"

    "github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
  got, err := readFile("sample.txt")
  if err != nil {
    t.Errorf("%v", err)
  }

  want := []pair{
    {rng{2, 4}, rng{6, 8}},
    {rng{2, 3}, rng{4, 5}},
    {rng{5, 7}, rng{7, 9}},
    {rng{2, 8}, rng{3, 7}},
    {rng{6, 6}, rng{4, 6}},
    {rng{2, 6}, rng{4, 8}},
  }

  if diff := cmp.Diff(want, got); diff != "" {
    t.Error(diff)
  }
}

func TestOne(t *testing.T) {
  data, err := readFile("sample.txt")
  if err != nil {
    t.Errorf("%v", err)
  }

  got, _ := one(data)
  if got != 2 {
    t.Errorf("Bad overlaps: want %d, got %d", 2, got)
  }
}

func TestTwo(t *testing.T) {
  data, err := readFile("sample.txt")
  if err != nil {
    t.Errorf("%v", err)
  }

  got, _ := two(data)
  if got != 4 {
    t.Errorf("Bad overlaps: want %d, got %d", 4, got)
  }
}
