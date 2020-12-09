package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func validNumber(x int, seq []int) bool {
	seen := make(map[int]bool)
	for _, n := range seq {
		need := x - n
		if seen[need] {
			return true
		}
		seen[n] = true
	}

	return false
}

func validateSequence(preamble int, seq []int) int {
	start, end := 0, preamble
	for end < len(seq) {
		if !validNumber(seq[end], seq[start:end]) {
			return seq[end]
		}
		start++
		end++
	}

	return -1
}

func main() {
	input := os.Args[1]

	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}
	var seq []int
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			log.Printf("bad line %q: %v", line, err)
			continue
		}
		seq = append(seq, n)
	}
	fmt.Println(validateSequence(25, seq))
}
