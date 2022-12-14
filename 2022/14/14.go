package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type Point struct {
	X, Y int
}

func NewPoint(s string) Point {
	ns := strings.Split(s, ",")
	return Point{common.MustInt(ns[0]), common.MustInt(ns[1])}
}
func (p Point) String() string {
	return fmt.Sprintf("[%d,%d]", p.X, p.Y)
}

type Object int

const (
	Air Object = iota
	Rock
	Sand
	Source
)

func (o Object) String() string {
	switch o {
	case Air:
		return "."
	case Rock:
		return "#"
	case Sand:
		return "o"
	case Source:
		return "+"
	default:
		return "."
	}
}

type Cave struct {
	Scan     map[Point]Object
	Min, Max *Point
	Floor    bool
}

func NewCave() *Cave {
	return &Cave{
		Scan: make(map[Point]Object),
	}
}
func (c *Cave) String() string {
	res := []string{}
	for y := c.Min.Y; y <= c.Max.Y; y++ {
		l := fmt.Sprintf("%3d: ", y)
		for x := c.Min.X; x <= c.Max.X; x++ {
			l += c.Scan[Point{x, y}].String()
		}
		res = append(res, l)
	}

	return strings.Join(res, "\n")
}
func (c *Cave) UpdateSize(p Point) {
	if c.Min == nil {
		c.Min = &Point{p.X, p.Y}
	} else {
		if p.X < c.Min.X {
			c.Min.X = p.X
		}
		if p.Y < c.Min.Y {
			c.Min.Y = p.Y
		}
	}
	if c.Max == nil {
		c.Max = &Point{p.X, p.Y}
	} else {
		if p.X > c.Max.X {
			c.Max.X = p.X
		}
		if p.Y > c.Max.Y {
			c.Max.Y = p.Y
		}
	}
}
func (c *Cave) AddLine(a, b Point) {
	if b.X < a.X || b.Y < a.Y {
		a, b = b, a
	}
	c.UpdateSize(a)
	c.UpdateSize(b)
	if a.X == b.X {
		for y := a.Y; y <= b.Y; y++ {
			c.Scan[Point{a.X, y}] = Rock
		}
	} else {
		for x := a.X; x <= b.X; x++ {
			c.Scan[Point{x, a.Y}] = Rock
		}
	}
}
func (c *Cave) AddStructure(l string) {
	ps := strings.Split(l, " -> ")
	var cur Point
	for n, coord := range ps {
		p := NewPoint(coord)
		if n == 0 {
			cur = p
			continue
		}
		c.AddLine(cur, p)
		cur = p
	}
}
func (c *Cave) Drop(p Point, floor int) *Point {
	pos := &Point{p.X, p.Y}
	if c.Scan[*pos] != Source {
		log.Printf("%s already occupied: %s", p, c.Scan[*pos])
		return nil
	}
	// Starting from the drop point, go down until you hit something, then
	// down-left, then down-right, then stop.
	if floor == 0 {
		floor = c.Max.Y
	}
	for pos.Y <= floor {
		pos.Y += 1
		if c.Scan[*pos] == Air {
			continue
		}
		pos.X -= 1
		if c.Scan[*pos] == Air {
			continue
		}
		pos.X += 2
		if c.Scan[*pos] == Air {
			continue
		}
		pos.Y -= 1
		pos.X -= 1
		break
	}
	if !c.Floor && pos.Y > c.Max.Y {
		return nil
	}
	c.Scan[*pos] = Sand
	c.UpdateSize(*pos)

	return pos
}

func one(c *Cave) int {
	start := Point{500, 0}
	n := 0
	for true {
		got := c.Drop(start, 0)
		if got == nil {
			return n
		}
		n++
	}
	return 0
}

func two(c *Cave) int {
	c.Floor = true
	start := Point{500, 0}
	n := 0
	floor := c.Max.Y
	for true {
		got := c.Drop(start, floor)
		if got == nil {
			return n
		}
		n++
	}
	return 0
}

func readFile(path string) *Cave {
	res := common.AsStrings(common.ReadTransformedFile(path, common.IgnoreBlankLines))
	c := NewCave()
	for _, l := range res {
		c.AddStructure(l)
	}
	source := Point{500, 0}
	c.Scan[source] = Source
	c.UpdateSize(source)

	return c
}

func main() {
	log.Println("Reading data...")
	data := readFile(os.Args[1])

	log.Println("Part A")
	res := one(data)
	fmt.Printf("%v\n", res)

	log.Println("Part B")
	data = readFile(os.Args[1])
	res = two(data)
	fmt.Printf("%v\n", res)
}
