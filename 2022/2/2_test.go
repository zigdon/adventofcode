package main

import (
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestRead(t *testing.T) {
  got, err := readFile("sample.txt")
  if err != nil {
    t.Errorf("%v", err)
  }

  want := Game{
    Rounds: []Round{
      {rock, paper},
      {paper, rock},
      {scissors, scissors}},
  }

  if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf(diff)
  }

  if got.Score() != 15 {
    t.Errorf("Bad score: %d", got.Score())
  }
}

func TestRead2(t *testing.T) {
  got, err := readFile2("sample.txt")
  if err != nil {
    t.Errorf("%v", err)
  }

  t.Logf("%v", got)
  want := Game{
    Rounds: []Round{
      {rock, rock},
      {paper, rock},
      {scissors, rock}},
  }

  if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf(diff)
  }

  if got.Score() != 12 {
    t.Errorf("Bad score: %d", got.Score())
  }
}
