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

func parseSchedule(lines []string) (int, map[int]int) {
	var start int
	buses := make(map[int]int)
	start, err := strconv.Atoi(lines[0])
	if err != nil {
		log.Fatalf("parse error for %q: %v", lines[0], err)
	}
	for n, bus := range strings.Split(lines[1], ",") {
		if bus == "x" {
			continue
		}
		no, err := strconv.Atoi(bus)
		if err != nil {
			log.Fatalf("skipping bus %q: %v", bus, err)
			continue
		}
		buses[no] = n
	}

	return start, buses
}

func nextBus(time int, buses map[int]int) (int, int) {
	var line int
	wait := -1
	for b := range buses {
		delay := b - time%b
		if delay < wait || wait == -1 {
			wait = delay
			line = b
		}
	}

	return wait, line
}

// {lineno: offset}
func findSequence(start int, sched map[int]int) int {
	fmt.Printf("sched: %v\n", sched)
	var keys []int
	pairs := make(map[int]int)
	for line, offset := range sched {
		keys = append(keys, line)
		if offset > 0 {
			pairs[line] = line - offset
			for pairs[line] < 0 {
				pairs[line] += line
			}
		} else {
			pairs[line] = 0
		}
	}
	sort.Ints(keys)
	var rev []int // line numbers, descending order
	for _, line := range keys {
		rev = append([]int{line}, rev...)
	}
	fmt.Printf("pairs: %v\n", pairs)

	// https://en.wikipedia.org/wiki/Chinese_remainder_theorem
	idx := 1
	var cur int
	if start == 0 {
		cur = pairs[rev[0]]
	} else {
		cur = start - (start % rev[0]) + pairs[rev[0]]
		fmt.Printf("starting at %d\n", cur)
	}

	step := rev[0]
	for {
		fmt.Printf("%d  %d\r", step, cur)
		if cur%rev[idx] != pairs[rev[idx]] {
			cur += step
			continue
		}
		step *= rev[idx]
		fmt.Printf("%d > %v %d %v\n", cur, rev[:idx], rev[idx], rev[idx+1:])
		idx++
		if idx >= len(rev) {
			return cur
		}
	}
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}

	start, sched := parseSchedule(strings.Split(string(data), "\n"))
	wait, line := nextBus(start, sched)
	fmt.Printf("%d minues until %d arrives => %d\n", wait, line, wait*line)

	fmt.Println("magic moment: ", findSequence(100000000000000, sched))
}
