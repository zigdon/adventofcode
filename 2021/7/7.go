package main

import (
	"log"

	"github.com/zigdon/adventofcode/common"
)

func readFile(path string) []int {
	return common.AsInts(common.ReadTransformedFile(
		path,
		common.Split(","),
	))
}

func align(fleet []int, dest int) int {
	cost := 0
	for _, crab := range fleet {
		c := crab - dest
		if c > 0 {
			cost += c
		} else {
			cost -= c
		}
	}

	return cost
}

// feeling silly to have to write this, but i can be good
// return 1+2+...+n, ideally without looping
func sum(dist int) int {
	if dist == 1 {
		return 1
	}
	return (dist + 1) * dist / 2
}

func align2(fleet []int, dest int) int {
	cost := 0
	for _, crab := range fleet {
		c := crab - dest
		if c < 0 {
			c = -c
		}

		cost += sum(c)
	}

	return cost
}

func max(nums []int) int {
	m := 0
	for _, n := range nums {
		if n <= m {
			continue
		}
		m = n
	}

	return m
}

// maybe something like this? https://en.wikipedia.org/wiki/Ternary_search
func plan(fleet []int, f func([]int, int) int) (int, int) {
	l := 0
	r := max(fleet)

	var lCost, rCost int
	for r-l > 1 {
		var n1, n2 int
		n2 = l + (r-l)/2
		if n2 == l {
			n2++
		}
		lCost = f(fleet, l)
		mCost := f(fleet, n2)
		rCost = f(fleet, r)
		log.Printf("%d(%d) - %d(%d) - %d(%d)", l, lCost, n2, mCost, r, rCost)

		var nCost int
		if lCost < rCost {
			n1 = l
			nCost = lCost
		} else {
			n1 = r
			nCost = rCost
		}

		if n1 < n2 {
			l, r = n1, n2
			lCost, rCost = nCost, mCost
		} else {
			l, r = n2, n1
			lCost, rCost = mCost, nCost
		}
	}

	if lCost < rCost {
		return l, lCost
	} else {
		return r, rCost
	}

}

func main() {
	fleet := readFile("input.txt")
	dest, cost := plan(fleet, align)
	log.Printf("=1= %d(%d)", dest, cost)
	dest, cost = plan(fleet, align2)
	log.Printf("=2= %d(%d)", dest, cost)
}
