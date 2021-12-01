package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseFile(data string) ([]int, error) {
	var res []int
	for _, line := range strings.Split(data, "\n") {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		num, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			return nil, err
		}
		res = append(res, num)
	}
	log.Printf("Read data: %d bytes, %d numbers", len(data), len(res))

	return res, nil
}

func findPair(target int, data []int) (int, int, error) {
	need := make(map[int]int)
	for _, got := range data {
		// Do we know this number works?
		if want, ok := need[got]; ok {
			return want, got, nil
		}

		// Record what we would need for success
		need[target-got] = got
	}

	return 0, 0, fmt.Errorf("no luck")
}

func findTriplet(target int, data []int) (int, int, int, error) {
	seen := []int{}
	want := make(map[int]struct{ a, b int })

	for _, got := range data {
		// Do we know this works?
		if pair, ok := want[got]; ok {
			return pair.a, pair.b, got, nil
		}

		// Figure out what we would need to get a match with all the previous numbers
		for _, i := range seen {
			if got+i >= target {
				continue
			}
			missing := target - got - i
			want[missing] = struct{ a, b int }{got, i}
		}
		seen = append(seen, got)
	}

	return 0, 0, 0, fmt.Errorf("never found the set!")
}

func main() {
	input := os.Args[1]
	target := 2020

	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	nums, err := parseFile(string(data))
	if err != nil {
		log.Fatalf("Error parsing data: %v", err)
	}
	a, b, err := findPair(target, nums)
	if err != nil {
		log.Fatalf("something wrong wrong: %v", err)
	}
	fmt.Printf("Found a pair: %d, %d -> %d\n", a, b, a*b)

	a, b, c, err := findTriplet(target, nums)
	if err != nil {
		log.Fatalf("didn't work out: %v", err)
	}
	fmt.Printf("Found a triplet: %d, %d, %d -> %d\n", a, b, c, a*b*c)
}
