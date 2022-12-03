package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zigdon/adventofcode/common"
)

func one(data []int) (int, error) {
  return 0, nil
}

func two(data []int) (int, error) {
  return 0, nil
}

func readFile(path string) ([]int, error) {
  res := common.ReadTransformedFile(path)

  return common.AsInts(res), nil
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
