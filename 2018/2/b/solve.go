package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "sort"
  "strings"
)

func readInput() ([]string, error) {
  in, err := ioutil.ReadFile("input")
  if err != nil {
    return nil, fmt.Errorf("couldn't read input: %v", err)
  }

  return strings.Split(string(in), "\n"), nil
}

func common(a, b string) string {
  if len(a) != len(b) {
    return ""
  }

  var c []string
  var diff int
  for i, la := range a {
    lb := rune(b[i])
    if lb == la {
      c = append(c, string(la))
    } else {
      diff++
    }
  }

  common := strings.Join(c, "")
  if diff == 1 {
    return common
  } else {
    return ""
  }
}

func main() {
  in, err := readInput()
  if err != nil {
    log.Fatalf(err.Error())
  }
  sort.Strings(in)

  prev := in[0]
  for _, id := range in[1:] {
    if len(id) == 0 {
      continue
    }
    if c := common(id, prev); c != "" {
      fmt.Printf("%s || %s = %s\n", id, prev, c)
      break
    }
    prev = id
  }
}
