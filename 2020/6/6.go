package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type form struct {
	Answers []string
	Anyone  int
}

func parseInput(data string) []*form {
	res := []*form{}

	cur := &form{}
	cnt := make(map[rune]int)
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			res = append(res, cur)
			cur = &form{}
			cnt = make(map[rune]int)
			continue
		}
		cur.Answers = append(cur.Answers, line)
		for _, r := range line {
			cnt[r] = cnt[r] + 1
		}
		cur.Anyone = len(cnt)
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
	total := 0
	for _, g := range groups {
		total = total + g.Anyone
	}
	fmt.Printf("Total uniq answers: %d", total)
}
