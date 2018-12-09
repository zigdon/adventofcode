package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "os"
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

  var lines []string
  for i, l := range strings.Split(string(in), "\n"){
    if len(l) == 0 {
      log.Printf("skipping blank line at #%d", i+1)
      continue
    }
    lines = append(lines, l)
  }

  return lines, nil
}

func reduce(p string, r *strings.Replacer) int {
  for {
    // log.Println(polymer)
    before := len(p)
    p = r.Replace(p)
    after := len(p)
    if before != after {
      // log.Printf("%d -> %d", before, after)
      continue
    } else {
      break
    }
  }

  return len(p)
}

func main() {
  in, err := readInput()
  if err != nil {
    log.Fatalf(err.Error())
  }

  var chars []string
  for i := 'a'; i <= 'z'; i++ {
    c := string(i)
    C := strings.ToUpper(c)
    chars = append(chars, c+C, "", C+c, "")
  }
  replacer := strings.NewReplacer(chars...)
  log.Println(chars)

  orig := in[0]
  best := 0
  target := ""
  for i := 'a'; i <= 'z'; i++ {
    c := string(i)
    C := strings.ToUpper(c)
    p := strings.Replace(orig, c, "", -1)
    p = strings.Replace(p, C, "", -1)
    reduced := reduce(p, replacer)
    log.Printf("%s: %d", c, reduced)
    if best == 0 || reduced < best {
      best = reduced
      target = c
    }
  }

  log.Printf("===== %s: %d", target, best)
}
