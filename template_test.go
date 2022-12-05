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

  want := []int{}

  if diff := cmp.Diff(want, got); diff != "" {
    t.Error(diff)
  }
}
