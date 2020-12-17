package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseSchedule(lines []string) (int, map[int]bool) {
	var start int
	buses := make(map[int]bool)
	start, err := strconv.Atoi(lines[0])
	if err != nil {
		log.Fatalf("parse error for %q: %v", lines[0], err)
	}
	for _, bus := range strings.Split(lines[1], ",") {
		if bus == "x" {
			continue
		}
		no, err := strconv.Atoi(bus)
		if err != nil {
			log.Fatalf("skipping bus %q: %v", bus, err)
			continue
		}
		buses[no] = true
	}

	return start, buses
}

func nextBus(time int, buses map[int]bool) (int, int) {
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

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}

	start, sched := parseSchedule(strings.Split(string(data), "\n"))
	wait, line := nextBus(start, sched)
	fmt.Printf("%d minues until %d arrives => %d", wait, line, wait*line)
}
