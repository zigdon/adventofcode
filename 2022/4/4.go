package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zigdon/adventofcode/common"
)

type pair struct {
  P1, P2 rng
}
type rng struct {
  Start, End int
}
func(r1 rng) FullOverlap(r2 rng) bool {
  return (r1.Start <= r2.Start && r1.End >= r2.End) || (r1.Start >= r2.Start && r1.End <= r2.End)
}
func(r1 rng) AnyOverlap(r2 rng) bool {
  return !(r1.End < r2.Start || r2.End < r1.Start)
}

func one(data []pair) (int, error) {
  res := 0
  for _, d := range data {
    if d.P1.FullOverlap(d.P2) {
      res += 1
    }
  }
  return res, nil
}

func two(data []pair) (int, error) {
  res := 0
  for _, d := range data {
    if d.P1.AnyOverlap(d.P2) {
      res += 1
    }
  }
  return res, nil
}

func readFile(path string) ([]pair, error) {
  data := common.ReadTransformedFile(
    path,
    common.IgnoreBlankLines,
    common.Split(","),
    common.Split("-"))
  res := []pair{}
  for _, d := range data {
    r1 := rng{
      common.MustInt(d.([][]string)[0][0]),
      common.MustInt(d.([][]string)[0][1])}
    r2 := rng{
      common.MustInt(d.([][]string)[1][0]),
      common.MustInt(d.([][]string)[1][1])}
    res = append(res, pair{r1, r2})
  }

  return res, nil
}

func main() {
  data, err := readFile(os.Args[1])
  if err != nil {
    log.Fatalf("%v", err)
  }

  res, err := one(data)
  if err != nil {
    log.Fatalf("%v", err)
  }
  fmt.Printf("%v\n", res)

  res, err = two(data)
  if err != nil {
    log.Fatalf("%v", err)
  }
  fmt.Printf("%v\n", res)
}
