package main

import (
	"fmt"
	"log"
	"os"
    "strings"

	"github.com/zigdon/adventofcode/common"
)

type Play int
func(p Play) String() string {
  return []string{"rock", "paper", "scissors"}[int(p)]
}

const (
  rock Play = iota
  paper
  scissors
)
var conv map[string]Play = map[string]Play{
  "A": rock,
  "B": paper,
  "C": scissors,
  "X": rock,
  "Y": paper,
  "Z": scissors,
}
type Round struct {
  P1, P2 Play
}
func(r Round) Judge() int {
  switch r.P2 {
    case rock:
      return []int{0, -1, 1}[int(r.P1)]
    case paper:
      return []int{1, 0, -1}[int(r.P1)]
    case scissors:
      return []int{-1, 1, 0}[int(r.P1)]
    default:
      log.Fatalf("Bad play %s", r)
  }
  return 99
}
func(r Round) Score() int {
  return int(r.P2)+1 + 3*(r.Judge()+1)
}

type Game struct {
  Rounds []Round
  Win bool  // player 2 wins?
}

func(g Game) Score() int {
  res := 0
  for _, r := range g.Rounds {
    res += r.Score()
  }

  return res
}

func(g Game) String() string {
  res := []string{}
  for _, r := range g.Rounds {
    res = append(res, fmt.Sprintf("%s:%s", r.P1, r.P2))
  }
  return strings.Join(res, "\n")
}

// x lose, y draw, z win
func readFile2(path string) (Game, error) {
  rounds := common.ReadTransformedFile(path, common.IgnoreBlankLines, common.SplitWords)

  res := Game{Rounds: []Round{}}
  for _, r := range rounds {
    p1 := conv[r.([]string)[0]]
    outcome := int(r.([]string)[1][0] - "X"[0])-1
    p2 := Play(((int(p1) + int(outcome) + 3) % 3))
    rd := Round{p1, p2}
    res.Rounds = append(res.Rounds, rd)
  }

  return res, nil
}

func readFile(path string) (Game, error) {
  rounds := common.ReadTransformedFile(path, common.IgnoreBlankLines, common.SplitWords)

  res := Game{Rounds: []Round{}}
  for _, r := range rounds {
    rd := Round{conv[r.([]string)[0]], conv[r.([]string)[1]]}
    res.Rounds = append(res.Rounds, rd)
  }

  return res, nil
}

func main() {
  data, err := readFile(os.Args[1])
  if err != nil {
    log.Fatalf("%v", err)
  }

  fmt.Printf("%v\n", data.Score())

  data, err = readFile2(os.Args[1])
  if err != nil {
    log.Fatalf("%v", err)
  }
  fmt.Printf("%v\n", data.Score())
}
