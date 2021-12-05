package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/zigdon/adventofcode/common"
)

type point struct {
	X, Y int
}

func (p point) dir(p2 point) (int, int) {
	var dx, dy int
	d := p2.X - p.X
	switch {
	case d > 0:
		dx = 1
	case d < 0:
		dx = -1
	}
	d = p2.Y - p.Y
	switch {
	case d > 0:
		dy = 1
	case d < 0:
		dy = -1
	}

	return dx, dy
}

type line struct {
	P1, P2   point
	Diagonal bool
}

func mustInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("couldn't make %q into an int: %v", s, err)
	}

	return n
}

func readFile(path string) []*line {
	data := common.ReadTransformedFile(
		path,
		common.IgnoreBlankLines,
		common.Split(" -> "),
		common.Split(","),
	)

	lines := []*line{}
	for _, in := range data {
		ent := in.([][]string)
		l := &line{
			P1: point{mustInt(ent[0][0]), mustInt(ent[0][1])},
			P2: point{mustInt(ent[1][0]), mustInt(ent[1][1])},
		}

		l.Diagonal = (l.P1.X != l.P2.X) && (l.P1.Y != l.P2.Y)
		lines = append(lines, l)
	}

	return lines
}

func findDanger(ls []*line, thresh int, useDiag bool) int {
	chart := make(map[point]int)
	for _, l := range ls {
		if !useDiag && l.Diagonal {
			continue
		}
		dx, dy := l.P1.dir(l.P2)
		x := l.P1.X
		y := l.P1.Y
		for {
			chart[point{x, y}]++
			if x == l.P2.X && y == l.P2.Y {
				break
			}
			x += dx
			y += dy
		}
	}

	count := 0
	for _, d := range chart {
		if d >= thresh {
			count++
		}
	}

	return count
}

func main() {
	lines := readFile("input.txt")
	fmt.Printf("Danger: %d\n", findDanger(lines, 2, false))
	fmt.Printf("Danger: %d\n", findDanger(lines, 2, true))
}
