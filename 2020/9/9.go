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

func findSeq(target int, seq []int) []int {
	sum := seq[0]
	start := 0
	end := 1
	for sum != target {
		if sum < target {
			end++
			sum = sum + seq[end-1]
		} else {
			sum = sum - seq[start]
			start++
		}

		if end >= len(seq) {
			return []int{}
		}
	}

	return seq[start:end]
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
	invalid := validateSequence(25, seq)
	fmt.Printf("Invalid number: %d\n", invalid)
	sum := findSeq(invalid, seq)
	if len(sum) > 0 {
		fmt.Printf("seq: %v\n", sum)
		sort.Ints(sum)
		fmt.Println(sum[0] + sum[len(sum)-1])
	}
}
