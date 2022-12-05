package main

import (
    "testing"
)

func TestReadFile(t *testing.T) {
  got, err := readFile("sample.txt")
  if err != nil {
    t.Errorf("%v", err)
  }

  want := []string{"p","L","P","v","t","s"}

  for i, b := range got {
    if want[i] != b.Invalid() {
      t.Errorf("Bad backpack %v: want %s, got %s", b, want, b.Invalid())
    }
  }
}

func TestOne(t *testing.T) {
  got, err := readFile("sample.txt")
  if err != nil {
    t.Errorf("%v", err)
  }
  wantP := 157
  gotP, _ := one(got)
  
  if gotP != wantP {
    t.Errorf("Bad total priority: want %d, got %d", wantP, gotP)
  }
}

func TestTwo(t *testing.T) {
  got, err := readFile("sample.txt")
  if err != nil {
    t.Errorf("%v", err)
  }
  wantB := 70
  gotB, _ := two(got)

  if gotB != wantB {
    t.Errorf("Bad total badge: want %d, got %d", wantB, gotB)
  }
}
