package main

import (
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
  got, err := readFile("sample.txt")
  if err != nil {
    t.Error(err)
  }

  want := [][]int{{1000, 2000, 3000}, {4000}, {5000, 6000}, {7000, 8000, 9000}, {10000}}

  if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf(diff)
  }
}

func TestOne(t *testing.T) {
  data, err := readFile("sample.txt")
  if err != nil {
    t.Error(err)
  }

  got, err := one(data)
  if err != nil {
    t.Error(err)
  }

  if got != 24000 {
    t.Fail()
  }
}

func TestTwo(t *testing.T) {
  data, err := readFile("sample.txt")
  if err != nil {
    t.Error(err)
  }

  got, err := two(data)
  if err != nil {
    t.Error(err)
  }

  if got != 45000 {
    t.Fail()
  }
}
