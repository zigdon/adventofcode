package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

// Attempting to implement https://en.wikipedia.org/wiki/A*_search_algorithm

type coord struct {
	X, Y int
}

func (c coord) neighbours() []coord {
	return []coord{
		{c.X, c.Y + 1},
		{c.X, c.Y - 1},
		{c.X + 1, c.Y},
		{c.X - 1, c.Y},
	}
}

type board struct {
	Danger [][]int
	Exit   coord
}

func (b *board) duplicate(n int) {
	nb := [][]int{}
	for y := range b.Danger {
		line := []int{}
		for i := 0; i < n; i++ {
			for _, n := range b.Danger[y] {
				d := n + i
				if d > 9 {
					d -= 9
				}
				line = append(line, d)
			}
		}
		nb = append(nb, line)
	}

	lnb := len(nb)
	for i := 1; i < n; i++ {
		for y := 0; y < lnb; y++ {
			line := []int{}
			for _, n := range nb[y] {
				d := n + i
				if d > 9 {
					d -= 9
				}
				line = append(line, d)
			}
			nb = append(nb, line)
		}
	}

	b.Danger = nb
	b.Exit.X = len(b.Danger[0]) - 1
	b.Exit.Y = len(b.Danger) - 1
}

type dir int

const (
	RIGHT dir = iota
	DOWN
	LEFT
	UP
)

func (d dir) String() string {
	return []string{"->", "vv", "<-", "^^"}[d]
}

type queue struct {
	CoordToPriority map[coord]int
	PriorityToCoord map[int][]coord
	Count           int
	Lowest          int
}

func newQueue() *queue {
	return &queue{
		CoordToPriority: make(map[coord]int),
		PriorityToCoord: make(map[int][]coord),
	}
}

func (q *queue) add(c coord, p int) {
	_, ok := q.CoordToPriority[c]
	if ok {
		return
	}
	q.CoordToPriority[c] = p
	q.PriorityToCoord[p] = append(q.PriorityToCoord[p], c)
	q.Count++
	if q.Lowest == 0 || p < q.Lowest {
		q.Lowest = p
	}
}

func (q *queue) isEmpty() bool {
	return q.Count == 0
}

func (q *queue) getLowest() (coord, error) {
	if q.isEmpty() {
		return coord{0, 0}, fmt.Errorf("empty queue")
	}

	c := q.PriorityToCoord[q.Lowest][0]
	q.Count--
	if len(q.PriorityToCoord[q.Lowest]) > 1 {
		q.PriorityToCoord[q.Lowest] = q.PriorityToCoord[q.Lowest][1:]
	} else {
		delete(q.PriorityToCoord, q.Lowest)
		q.Lowest = 0
		for p, cs := range q.PriorityToCoord {
			if len(cs) == 0 {
				continue
			}
			if q.Lowest == 0 || p < q.Lowest {
				q.Lowest = p
			}
		}
	}

	return c, nil
}

func (q *queue) has(c coord) bool {
	_, ok := q.CoordToPriority[c]
	return ok
}

type cameFrom struct {
	Src  map[coord]coord
	Cost map[coord]int
}

func (cf cameFrom) getSrc(c coord) coord {
	return cf.Src[c]
}

func (cf cameFrom) getScore(c coord) int {
	return cf.Cost[c]
}

func (cf cameFrom) set(src, dst coord, score int) {
	cf.Cost[dst] = score
	cf.Src[dst] = src
}

func (cf cameFrom) path(c coord) []coord {
	log.Printf("getting path ending at %v", c)
	out := []coord{c}
	for c.X != 0 && c.Y != 0 {
		c = cf.getSrc(c)
		log.Printf(" -> %v", c)
		out = append(out, c)
	}

	return out
}

func (b *board) aStar() ([]coord, int) {
	h := func(c coord) int {
		return int(math.Abs(float64(b.Exit.Y-c.Y)) + math.Abs(float64(b.Exit.X-c.X)))
	}

	q := newQueue()
	start := coord{0, 0}
	q.add(start, h(start))

	cf := cameFrom{
		Src:  map[coord]coord{},
		Cost: map[coord]int{{0, 0}: 0},
	}

	for !q.isEmpty() {
		current, err := q.getLowest()
		if err != nil {
			log.Fatal(err)
		}
		if current.X == b.Exit.X && current.Y == b.Exit.Y {
			return cf.path(current), cf.getScore(b.Exit)
		}

		for _, c := range current.neighbours() {
			if c.Y < 0 || c.X < 0 || c.Y > b.Exit.Y || c.X > b.Exit.X {
				continue
			}
			tScore := cf.getScore(current) + b.Danger[c.Y][c.X]
			if c.X == 0 && c.Y == 0 {
				continue
			}
			if cf.getScore(c) == 0 || tScore < cf.getScore(c) {
				cf.set(current, c, tScore)
			}
			if !q.has(c) {
				q.add(c, tScore+h(c))
			}
		}
	}

	log.Print("ran out of ideas:")
	for to, from := range cf.Src {
		log.Printf(" %v(%d) -> %v(%d)", from, cf.Cost[from], to, cf.Cost[to])
	}

	return nil, 0
}

func (b *board) draw(p []coord) {
	out := [][]string{}
	for y := 0; y <= b.Exit.Y; y++ {
		line := []string{}
		for x := 0; x <= b.Exit.X; x++ {
			line = append(line, ".")
		}
		out = append(out, line)
	}

	for _, s := range p {
		out[s.Y][s.X] = "#"
	}

	for _, l := range out {
		log.Print(strings.Join(l, " "))
	}
}

func readFile(path string) board {
	data := common.AsIntGrid(common.ReadTransformedFile(
		path,
		common.IgnoreBlankLines,
		common.Split(""),
	))

	return board{
		Danger: data,
		Exit:   coord{len(data[0]) - 1, len(data) - 1},
	}
}

func main() {
	b := readFile("input.txt")
	p, s := b.aStar()
	log.Printf("%v", p)
	log.Printf("score: %d", s)

	b = readFile("input.txt")
	b.duplicate(5)
	p, s = b.aStar()
	log.Printf("%v", p)
	log.Printf("other score: %d", s)
}
