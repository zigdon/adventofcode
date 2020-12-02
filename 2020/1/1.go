package main

import "fmt"
import "io/ioutil"
import "log"
import "os"
import "strconv"
import "strings"

func findPair(target int, path string) (int, int, error) {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return 0, 0, err
  }

  need := make(map[int]int)
  for _, line := range(strings.Split(string(data), "\n")) {
    got, err := strconv.Atoi(strings.Trim(line, "\n"))
    if err != nil {
      return 0, 0, err
    }
    if want, ok  := need[got]; ok {
      return want, got, nil
    }
    need[target - got] = got
  }

  return 0, 0, fmt.Errorf("no luck")
}

func main() {
  input := os.Args[1]
  target := 2020

  a, b, err := findPair(target, input)
  if err != nil {
    log.Fatalf("something wrong wrong: %v", err)
  }
  fmt.Printf("Found a pair: %d, %d -> %d", a, b, a*b)
}
