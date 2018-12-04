package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "strings"
)

func readInput() ([]string, error) {
  in, err := ioutil.ReadFile("input")
  if err != nil {
    return nil, fmt.Errorf("couldn't read input: %v", err)
  }

  return strings.Split(string(in), "\n"), nil
}

func find(data map[rune]int, target int) int {
  for _, v := range data {
    if v == target {
      return 1
    }
  }
  return 0
}

func main() {
  in, err := readInput()
  if err != nil {
    log.Fatalf(err.Error())
  }

  var d, t int
  for _, id := range in {
    if len(id) == 0 {
      continue
    }
    count := make(map[rune]int)
    for _, c := range id {
      count[c]++
    }
    d += find(count, 2)
    t += find(count, 3)
  }

  fmt.Println(d*t)
}
