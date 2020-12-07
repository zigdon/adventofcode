package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type rule struct {
	CanContain   map[string]int
	CanBeFoundIn map[string]bool
}

func newRule() rule {
	return rule{
		CanContain:   make(map[string]int),
		CanBeFoundIn: make(map[string]bool),
	}
}

func parseRules(data string) map[string]rule {
	res := make(map[string]rule)

	ts := strings.TrimSuffix
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.Contains(line, "no other bags") {
			bits := strings.Split(line, " bags ")
			res[bits[0]] = newRule()
			continue
		}
		bits := strings.Split(line, " contain ")
		color := ts(ts(bits[0], "s"), " bag")
		contents := strings.Split(bits[1], ", ")
		cur := newRule()

		for _, bag := range contents {
			innerBits := strings.SplitN(bag, " ", 2)
			count, err := strconv.Atoi(innerBits[0])
			if err != nil {
				log.Fatalf("Number expected in %q: %v", bag, err)
			}
			innerColor := ts(ts(ts(innerBits[1], "."), "s"), " bag")
			if cur.CanContain == nil {
				cur.CanContain = make(map[string]int)
			}
			cur.CanContain[innerColor] = count
		}

		res[color] = cur
	}

	for color, rule := range res {
		if rule.CanContain != nil {
			for innerColor := range rule.CanContain {
				res[innerColor].CanBeFoundIn[color] = true
			}
		}
	}

	return res
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("can't read input: %v", err)
	}
	parseRules(string(data))
}
