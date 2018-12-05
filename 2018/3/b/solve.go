package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "strconv"
  "strings"
)

func readInput() ([]string, error) {
  var filename string
  if len(os.Args) > 1 {
    filename = os.Args[1]
  } else {
    filename = "input"
  }
  in, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, fmt.Errorf("couldn't read %s: %v", filename, err)
  }

  return strings.Split(string(in), "\n"), nil
}

func mark(f *[1000][1000][]int, x, y, w ,h, id int64) int {
  dup := 0
  for xp := x; xp < x+w; xp++ {
    for yp := y; yp < y+h; yp++ {
      if (*f)[xp][yp] == nil {
        (*f)[xp][yp] = []int{}
      } else {
        if len((*f)[xp][yp]) == 1 {
          dup++
        }
      }
      (*f)[xp][yp] = append((*f)[xp][yp], int(id))
    }
  }

  return dup
}

func main() {
  in, err := readInput()
  if err != nil {
    log.Fatalf(err.Error())
  }

  // x, y, list of claim ids
  var fabric [1000][1000][]int
  var dups int
  for _, l := range in {
    if len(l) == 0 {
      continue
    }
    bits := strings.Split(l, " ")
    if bits[1] != "@" {
      log.Fatalf("Invalid format: %q", l)
    }
    id, _ := strconv.ParseInt(bits[0][1:], 10, 32)
    sp := strings.SplitN(bits[2], ",", 2)
    x, _ := strconv.ParseInt(sp[0], 10, 32)
    y, _ := strconv.ParseInt(strings.TrimSuffix(sp[1], ":"), 10, 32)
    sp = strings.SplitN(bits[3], "x", 2)
    w, _ := strconv.ParseInt(sp[0], 10, 32)
    h, _ := strconv.ParseInt(sp[1], 10, 32)
    dups += mark(&fabric, x, y, w, h, id)
  }
  fmt.Println(dups)

  // Find untarnished samples

  candidates := make(map[int]bool)
  for xp := 0; xp < 1000; xp++ {
    for yp := 0; yp < 1000; yp++ {
      valid := (len(fabric[xp][yp]) == 1)
      for _, id := range(fabric[xp][yp]) {
        c, ok := candidates[id]
        if ok && !c {
          continue
        }
        candidates[id] = valid
      }
    }
  }

  for k, v := range candidates {
    if !v {
      continue
    }
    fmt.Println(k)
  }

}
