package main

import (
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestEval(t *testing.T) {
  _ = getParser()

  _ = []struct{
    exp string
    want int
  }{
    {
      exp: "1 + 2 * 3 + 4 * 5 + 6",
      want: 71,
    },
    {
      exp: "1 + (2 * 3) + (4 * (5 + 6))",
      want: 51,
    },
    {
      exp: "2 * 3 + (4 * 5)",
      want: 26,
    },
    {
      exp: "2 * 3 + (4 * 5)",
      want: 437,
    },
    {
      exp: "5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))",
      want: 12240,
    },
    {
      exp: "((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2",
      want: 13632,
    },
  }
}

func TestParser(t *testing.T) {
  p := getParser()
  got := &Expr{}
  err := p.ParseString("", "5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))", got)
  if err != nil {
    t.Fatalf("Can't parse example: %v", err)
  }

  want := &Expr{
    Left: &Value{5},
    Right: []*Ops{
      &Ops{Op: "*", Val: &Value{9}},
      &Ops{Op: "*", Val: &Sub{
        Left: &Value{7},
        Right: nil,
        },
      },
    },
  } 

  if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("Bad parsing: %s", diff)
  }
}
