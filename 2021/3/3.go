package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/zigdon/adventofcode/common"
)

func countBits(data []string) []int {
	bitCounts := make([]int, len(data[0]))
	for _, n := range data {
		for i, b := range n {
			if b == '1' {
				bitCounts[i]++
			}
		}
	}

	return bitCounts
}

func getRates(bitCounts []int, total int) (int, int) {
	gamma := 0
	epsilon := 0
	p := 1
	for i := range bitCounts {
		if bitCounts[len(bitCounts)-i-1] > total/2 {
			gamma += p
		} else {
			epsilon += p
		}
		p = p * 2
	}

	return gamma, epsilon
}

func getRating(data []string, ge bool) int {
	var res []string = data

	pos := 0
	for {
		counts := countBits(res)
		next := []string{}
		for _, d := range res {
			var goal bool
			// Go won't let me do int > float, which I think is mean.
			if ge {
				goal = counts[pos]*10 >= len(res)*5
			} else {
				goal = counts[pos]*10 < len(res)*5
			}

			if (goal && d[pos] == '1') || (!goal && d[pos] == '0') {
				next = append(next, d)
			}
		}
		if len(res) == len(next) {
			log.Fatalf("No solution found: %v", next)
		}
		res = next
		if len(res) == 1 {
			break
		}
		pos++
	}

	n, err := strconv.ParseInt(res[0], 2, 32)
	if err != nil {
		log.Fatalf("Can't parse %s: %v", res[0], err)
	}

	return int(n)
}

func main() {
	data := common.AsStrings(common.ReadTransformedFile("input.txt", common.IgnoreBlankLines))
	bitCounts := countBits(data)
	g, e := getRates(bitCounts, len(bitCounts))
	fmt.Printf("%d(%b) * %d(%b) = %d\n", g, g, e, e, g*e)

	o2 := getRating(data, true)
	co2 := getRating(data, false)
	fmt.Printf("%d(%b) * %d(%b) = %d\n", o2, o2, co2, co2, o2*co2)
}
