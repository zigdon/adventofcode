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

}
