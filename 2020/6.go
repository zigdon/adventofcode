package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type form struct {
	Answers  []string
	Anyone   int
	Everyone int
}

func parseInput(data string) []*form {
	res := []*form{}

	cur := &form{}
	cnt := make(map[rune]int)
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			cur.Anyone = len(cnt)
			for _, v := range cnt {
				if v == len(cur.Answers) {
					cur.Everyone = cur.Everyone + 1
				}
			}
			res = append(res, cur)
			cur = &form{}
			cnt = make(map[rune]int)
			continue
		}
		cur.Answers = append(cur.Answers, line)
		for _, r := range line {
			cnt[r] = cnt[r] + 1
		}
	}

	return res
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("Error reading %q: %v", input, err)
	}
	groups := parseInput(string(data))
	totalAny := 0
	totalEvery := 0
	for _, g := range groups {
		totalAny = totalAny + g.Anyone
		totalEvery = totalEvery + g.Everyone
	}
	fmt.Printf("Total anyone answered: %d", totalAny)
	fmt.Printf("Total everyone answered: %d", totalEvery)
}
