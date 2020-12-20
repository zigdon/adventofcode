package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ent struct {
	last int
	init bool
}

func getSequence(init []int, pos int) int {
	mem := make(map[int]int)
	var seq []int
	var last, next int
	for step := 0; step < len(init); step++ {
		next = init[step]
		mem[last] = step
		last = next
		seq = append(seq, next)
	}
	for step := len(init); step < pos; step++ {
		// log.Printf("step=%d", step)
		if seen, ok := mem[last]; ok {
			next = step - seen
			// log.Printf("%d was last seen on %d, so next=%d", last, seen, next)
		} else {
			next = 0
		}
		mem[last] = step
		last = next
		seq = append(seq, next)
		// log.Printf("mem: %v", mem)
	}
	// log.Printf("%v", seq)

	return last
}

func main() {
	input := []int{}
	count, _ := strconv.Atoi(os.Args[2])
	for _, s := range strings.Split(os.Args[1], ",") {
		n, _ := strconv.Atoi(s)
		input = append(input, n)
	}

	fmt.Println(input)
	fmt.Println(getSequence(input, count))
}
