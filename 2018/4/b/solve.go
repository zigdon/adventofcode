package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "sort"
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

type hour [60]int

type guard struct {
  id int
  schedule hour
  log map[int]hour
  totalSleep int
  mostWhen int
  mostCount int
}

func (g guard) LogSleep(day, start, end int) guard {
  l := g.log[day]
  for m := start; m < end; m++ {
    l[m]++
    g.schedule[m]++
    if g.schedule[m] > g.mostCount {
      g.mostWhen = m
      g.mostCount = g.schedule[m]
    }
  }
  g.log[day] = l
  g.totalSleep += (end - start)
  return g
}

type guards map[int]*guard

func (g guards) New(id int) {
  g[id] = &guard{id: id, log: make(map[int]hour), schedule: hour{}}
}

func parseLog(lines []string) (guards, error) {
  var current int
  guards := make(guards)
  sleeping := -1
  pi := func(s string) int {
    i, err := strconv.ParseInt(s, 10, 32)
    if err != nil {
      log.Fatalf("Error parsing %q: %v", s, err)
    }
    return int(i)
  }
  for _, line := range lines {
    if line == "" {
      continue
    }
    bits := strings.Fields(line)
    date := strings.Split(bits[0], "-")
    day := pi(date[2])
    time := strings.Split(bits[1], ":")
    hr := pi(time[0])
    min := pi(strings.TrimRight(time[1], "]"))
    switch bits[2] {
    case "Guard":
      current = pi(strings.TrimLeft(bits[3], "#"))
      if _, ok := guards[current]; !ok {
        guards.New(current)
      }
      sleeping = -1
    case "falls":
      if hr == 0 {
        sleeping = min
      } else {
        sleeping = 0
      }
    case "wakes":
      if hr == 0 {
        if min <= sleeping {
          return guards, fmt.Errorf("%d wakes before sleeping!", current)
        }
        g := guards[current].LogSleep(day, sleeping, min)
        guards[current] = &g
      } else {
        sleeping = -1
      }
    default:
      return guards, fmt.Errorf("can't parse line %q", line)
    }
  }

  return guards, nil
}

func main() {
  in, err := readInput()
  if err != nil {
    log.Fatalf(err.Error())
  }

  sort.Strings(in)
  sched, err := parseLog(in)
  if err != nil {
    log.Fatalf(err.Error())
  }

  var sleepy *guard
  for _, g := range sched {
    if sleepy == nil {
      sleepy = g
      continue
    }

    if g.mostCount > sleepy.mostCount {
      sleepy = g
    }
  }

  log.Printf("Sleepingest guard: %d (%d minutes at :%d)", sleepy.id, sleepy.mostCount, sleepy.mostWhen)
  log.Printf("Answer: %d", sleepy.mostWhen * sleepy.id)
}
