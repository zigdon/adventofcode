package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	maxRows = 128
	maxCols = 8
)

type seat struct {
	BSP          string
	Row, Col, Id int
}

func (s *seat) parseBSP() int {
	var err error
	s.Row, err = bsp(s.BSP[:7], "F", "B", maxRows)
	if err != nil {
		log.Printf("BSP error: %v", err)
		return -1
	}
	s.Col, err = bsp(s.BSP[7:], "L", "R", maxCols)
	if err != nil {
		log.Printf("BSP error: %v", err)
		return -1
	}
	s.Id = s.Row*8 + s.Col
	return s.Id
}

func bsp(code string, low, high string, max int) (int, error) {
	start := 0
	end := max
	idx := 0
	for ; end-start > 1; idx++ {
		span := end - start
		cur := string(code[idx])
		if cur == low {
			end = end - span/2
		} else if cur == high {
			start = start + span/2
		} else {
			return -1, fmt.Errorf("bad rune in %q #%d: %q", code, idx, cur)
		}
	}

	return start, nil
}

func newSeat(id string) *seat {
	res := &seat{BSP: id}
	res.parseBSP()
	return res
}

func findOpenSeat(plane []bool) int {
	for i, full := range plane {
		if full || i == 0 {
			continue
		}
		if plane[i-1] && plane[i+1] {
			return i
		}
	}

	return -1
}

func main() {
	var plane [maxRows * maxCols]bool
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("Error reading input %q: %v", input, err)
	}
	seats := strings.Split(string(data), "\n")
	best := 0
	for _, s := range seats {
		bsp := strings.TrimSpace(s)
		if len(bsp) == 0 {
			continue
		}
		got := newSeat(bsp)
		plane[got.Id] = true
		if got.Id > best {
			best = got.Id
		}
	}
	fmt.Printf("Highest ID: %d", best)
	fmt.Printf("My seat: %d", findOpenSeat(plane[:]))
}
