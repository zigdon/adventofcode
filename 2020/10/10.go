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

// Given a sorted chain of adapters, how ways could they be combined
func countChains(adapters []int) int {
	if len(adapters) <= 2 {
		return 1
	}
	start := adapters[0]
	count := 0
	count = count + countChains(adapters[1:])

	if len(adapters) >= 3 && adapters[2]-start <= 3 {
		count = count + countChains(adapters[2:])
	}

	if len(adapters) >= 4 && adapters[3]-start <= 3 {
		count = count + countChains(adapters[3:])
	}

	return count
}

func sectionAdapters(adapters []int) [][]int {
	res := [][]int{}
	sort.Ints(adapters)
	section := []int{}
	last := adapters[0]
	for _, n := range adapters[1:] {
		if n-last < 3 {
			if len(section) == 0 {
				section = append(section, last)
			}
			section = append(section, n)
		} else {
			if len(section) > 1 {
				res = append(res, section)
			}
			section = []int{}
		}
		last = n
	}

	if len(section) > 0 {
		res = append(res, section)
	}

	return res
}

func buildChain(adapters []int) int {
	count := 1
	sort.Ints(adapters)

	sections := sectionAdapters(adapters)
	for _, s := range sections {
		c := countChains(s)
		count = count * c
	}

	return count
}

func getDist(adapters []int) map[int]int {
	sort.Ints(adapters)
	dist := make(map[int]int)
	var last int
	for _, n := range adapters {
		dist[n-last]++
		last = n
	}
	// And one for the phone
	dist[3]++

	return dist
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}

	adapterStrs := strings.Split(string(data), "\n")
	var adapters []int
	for _, a := range adapterStrs {
		n, err := strconv.Atoi(a)
		if err != nil {
			log.Printf("Skipping %q: %v", a, err)
			continue
		}
		adapters = append(adapters, n)
	}

	dist := getDist(adapters)
	fmt.Printf("%v\n", dist)
	fmt.Println(dist[1] * dist[3])
	adapters = append(adapters, 0)
	fmt.Printf("chains: %d", buildChain(adapters))
}
