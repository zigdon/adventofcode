package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func countIncrease(data []int) int {
	last := data[0]
	count := 0
	for _, n := range data[1:] {
		if n > last {
			count++
		}
		last = n
	}

	return count
}

func countWindowIncrease(data []int, window int) int {
	last := 0
	cur := 0
	for _, n := range data[0:window] {
		cur += n
	}
	count := 0
	for i, n := range data[window:] {
		last = cur
		cur = cur + n - data[i]
		if cur > last {
			count++
		}
	}

	return count
}

func readFile(path string) []int {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Can't read input: %v", err)
	}
	var data []int
	for _, l := range strings.Split(string(text), "\n") {
		if len(strings.TrimSpace(l)) == 0 {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(l))
		if err != nil {
			log.Fatalf("Can't convert %q to number: %v", l, err)
			continue
		}
		data = append(data, n)
	}

	return data
}

func main() {
	input := os.Args[1]
	data := readFile(input)
	fmt.Printf("Increase: %d\n", countIncrease(data))
	fmt.Printf("Increase (3): %d\n", countWindowIncrease(data, 3))
}
