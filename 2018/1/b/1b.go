package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "strconv"
  "strings"
)

func main() {
  in, err := ioutil.ReadFile("input")
  if err != nil {
    log.Fatalf("couldn't read input: %v", err)
  }

  var tot int
  seen := make(map[int]bool)
  for {
    fmt.Println("starting...")
    for _, d := range strings.Split(string(in), "\n") {
      if len(d) == 0 {
        continue
      }
      n, err := strconv.ParseInt(d, 10, 64)
      if err != nil {
        fmt.Printf("couldn't parse %q: %v", d, err)
        continue
      }
      if seen[tot] {
        log.Fatalf("repeat freq: %d", tot)
      }
      seen[tot] = true
      tot += int(n)
    }
  }

  fmt.Println(tot)
}
