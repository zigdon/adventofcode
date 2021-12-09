package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/zigdon/adventofcode/common"
)

type coord struct {
	x, y int
}

type hMap map[coord]int

func (h hMap) print() {
	for y := 0; ; y++ {
		if h.get(0, y) == -1 {
			break
		}
		l := fmt.Sprintf("%d: ", y)
		for x := 0; ; x++ {
			v := h.get(x, y)
			if v == -1 {
				break
			}
			l += fmt.Sprintf("%d ", v)
		}
		log.Print(l)
	}
}

func (h hMap) get(x, y int) int {
	p, ok := h[coord{x, y}]
	if !ok {
		return -1
	}
	return p
}

func (h hMap) isLow(c coord) bool {
	v := h[c]
	for _, delta := range []coord{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} {
		n := h.get(c.x+delta.x, c.y+delta.y)
		if n == -1 {
			continue
		}
		if n <= v {
			return false
		}
	}

	return true
}

// lets assume all the basins are bordered by 9s
func (h hMap) basinSize(c coord) int {
	count := 0
	next := map[coord]bool{c: true}
	queue := map[coord]bool{}
	seen := map[coord]bool{}
	for len(next) > 0 {
		queue = next
		next = map[coord]bool{}
		for p := range queue {
			if seen[p] {
				continue
			}
			seen[p] = true
			v := h.get(p.x, p.y)
			count++

			for _, delta := range []coord{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} {
				nc := coord{p.x + delta.x, p.y + delta.y}
				n := h.get(nc.x, nc.y)
				if n == -1 {
					continue
				}
				if n > v && n != 9 {
					next[nc] = true
				}
			}
		}
	}

	return count
}

func readFile(path string) hMap {
	data := common.ReadTransformedFile(
		path,
		common.IgnoreBlankLines,
		common.Split(""),
	)

	hm := hMap{}
	for y, dl := range data {
		l := dl.([]string)
		for x, h := range l {
			hn, err := strconv.Atoi(h)
			if err != nil {
				log.Fatalf("can't convert %q to int: %v", h, err)
			}
			hm[coord{x: x, y: y}] = hn
		}
	}

	return hm
}

func findLocalLow(m hMap) ([]int, []coord) {
	res := []int{}
	lows := []coord{}

	for c, v := range m {
		if m.isLow(c) {
			res = append(res, v)
			lows = append(lows, c)
		}
	}

	return res, lows
}

func findRisk(m hMap) int {
	res := 0
	lows, _ := findLocalLow(m)
	for _, h := range lows {
		res += 1 + h
	}

	return res
}

func main() {
	data := readFile("input.txt")
	risk := findRisk(data)
	log.Printf("risk: %d", risk)

	sizes := []int{}
	_, lows := findLocalLow(data)
	for _, l := range lows {
		sizes = append(sizes, data.basinSize(l))
	}

	sort.Ints(sizes)
	sn := len(sizes) - 1
	log.Printf("basins: %d * %d * %d = %d", sizes[sn], sizes[sn-1], sizes[sn-2], sizes[sn]*sizes[sn-1]*sizes[sn-2])
}
