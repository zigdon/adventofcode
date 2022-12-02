package main

import (
	"fmt"
	"log"
	"os"
    "sort"

	"github.com/zigdon/adventofcode/common"
)

func one(data [][]int) (int, error) {
  best := 0
  for _, d := range data {
    a := 0
    for _, f := range d {
      a += f
    }
    if a > best {
      best = a
    }
  }
  return best, nil
}

func two(data [][]int) (int, error) {
  totals := []int{}
  for _, d := range data {
    a := 0
    for _, f := range d {
      a += f
    }
    totals = append(totals, a)
  }
  sort.Sort(sort.Reverse(sort.IntSlice(totals)))
  return totals[0] + totals[1] + totals[2], nil
}

func readFile(path string) ([][]int, error) {
  res := common.ReadTransformedFile(path, common.Block)

  return common.AsIntGrid(res), nil
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
