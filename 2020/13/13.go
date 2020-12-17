package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func validateTime(ts int, sched map[int]int) bool {
	for line, off := range sched {
		if (ts+off)%line != 0 {
			return false
		}
	}
	return true
}

func findSequence(start, end int, sched map[int]int, verbose bool) int {
	var max, offset int
	for line, ord := range sched {
		fmt.Printf("%d:%d, ", ord, line)
		if line > max {
			max = line
			offset = ord
		}
	}
	fmt.Println()

	cur := start + (start % max) + offset
	fmt.Printf("starting at %d\n", cur)
	for {
		if validateTime(cur, sched) {
			fmt.Println()
			return cur
		}
		cur += max
		if end > 0 && cur > end {
			return -1
		}
		if verbose {
			fmt.Printf("%d\r", cur)
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

	fmt.Println("magic moment: ", findSequence(100000000000000, 0, sched, true))
}
