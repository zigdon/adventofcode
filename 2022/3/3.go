package main

import (
	"fmt"
	"log"
	"os"
    "strings"

	"github.com/zigdon/adventofcode/common"
)

type backpack struct {
  C1, C2 string
  BadPriority int
  invalid string
}
func (b *backpack) String() string {
  return fmt.Sprintf("[%s,%s]", b.C1, b.C2)
}
func (b *backpack) Invalid() string {
  if b.invalid != "" {
    return b.invalid
  }
  res := []rune{}
  seen := map[rune]bool{}
  for _, c := range b.C1 {
    if seen[c] { continue }
    seen[c] = true
    if strings.ContainsRune(b.C2, c) {
      res = append(res, c)
      b.BadPriority += GetPriority(c)
    }
  }
  b.invalid = string(res)

  return b.invalid
}

func GetPriority(r rune) int {
  if r >= 'a' {
    return int(r - 'a') + 1
  }
  return int(r - 'A') + 27
}

func findBadge(bs []*backpack) rune {
  count := map[rune]int{}
  for _, b := range bs {
    seen := map[rune]bool{}
    for _, c := range b.C1 {
      if seen[c] { continue }
      seen[c] = true
      count[c] += 1
    }
    for _, c := range b.C2 {
      if seen[c] { continue }
      seen[c] = true
      count[c] += 1
    }
  }

  for k, v := range count {
    if v == len(bs) {
      bstr := []string{}
      for _, b := range bs {
        bstr = append(bstr, b.String())
      }
      return k
    }
  }

  return 0
}

func one(data []*backpack) (int, error) {
  res := 0
  for _, b := range data {
    if b.Invalid() != "" {
      res += b.BadPriority
    }
  }
  return res, nil
}

func two(data []*backpack) (int, error) {
  c := 0;
  res := 0
  for c < len(data) {
    res += GetPriority(findBadge(data[c:c+3]))
    c += 3
  }

  return res, nil
}

func readFile(path string) ([]*backpack, error) {
  backpacks := common.AsStrings(common.ReadTransformedFile(path, common.IgnoreBlankLines))
  res := []*backpack{}
  for _, b := range backpacks {
    mid := len(b)/2
    backpack := backpack{
      C1: b[:mid],
      C2: b[mid:],
    }
    res = append(res, &backpack)
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
