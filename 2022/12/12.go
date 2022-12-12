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

func (p Point) String() string {
	return fmt.Sprintf("[%d,%d]", p.X, p.Y)
}
func (p Point) Eq(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y
}
func (p Point) Up() Point {
	return Point{p.X, p.Y - 1}
}
func (p Point) Down() Point {
	return Point{p.X, p.Y + 1}
}
func (p Point) Left() Point {
	return Point{p.X - 1, p.Y}
}
func (p Point) Right() Point {
	return Point{p.X + 1, p.Y}
}
func (p Point) Arrow(pt Point) rune {
	if pt.X == p.X {
		if pt.Y > p.Y {
			return 'v'
		}
		return '^'
	}
	if pt.X > p.X {
		return '>'
	}
	return '<'
}

var pID = 0

type Path struct {
	ID    int
	Map   *Map
	Steps []Point
	Been  map[Point]bool
}

func (p *Path) At() Point {
	return p.Steps[len(p.Steps)-1]
}

func (p *Path) Valid(pt Point) bool {
	curAlt := p.Map.Alt[p.At()]
	newAlt, ok := p.Map.Alt[pt]
	if !ok {
		return false
	}
	if newAlt > curAlt+1 {
		return false
	}

	return true
}

func NewPath(m *Map, p ...Point) *Path {
	path := &Path{ID: pID, Map: m, Steps: append([]Point{}, p...)}
	pID++
	been := map[Point]bool{}
	for _, at := range p {
		been[at] = true
	}
	path.Been = been

	return path
}

func (p *Path) String() string {
	res := map[Point]rune{}
	for pt := range p.Map.Alt {
		res[pt] = '.'
	}
	res[p.At()] = '#'

	for i, pt := range p.Steps {
		if i == 0 {
			res[pt] = 'S'
		}
		if i+1 == len(p.Steps) {
			break
		}
		res[pt] = pt.Arrow(p.Steps[i+1])
	}

	out := []string{fmt.Sprintf("Path #%d", p.ID)}
	out = append(out, strings.Split(p.Map.String(), "\n")...)
	for y := 0; y < p.Map.Size.Y; y++ {
		l := out[y+1] + "   "
		for x := 0; x < p.Map.Size.X; x++ {
			l += string(res[Point{x, y}])
		}
		out[y+1] = l
	}

	return strings.Join(out, "\n")
}

func (p *Path) Add(pt Point) *Path {
	newPath := NewPath(p.Map, append(p.Steps, pt)...)
	return newPath
}

func (p *Path) NextSteps() []*Path {
	// Check the end point of the path, of the 4 available directions:
	// - Have we been there?
	// - Is it valid
	// Return a new path with the valid steps
	at := p.At()
	res := []*Path{}
	names := []string{"up", "right", "left", "down"}
	good := []string{}
	bad := []string{}
	for i, pt := range []Point{at.Up(), at.Right(), at.Left(), at.Down()} {
		if p.Been[pt] {
			bad = append(bad, names[i])
			continue
		}
		if p.Valid(pt) {
			good = append(good, names[i])
			next := p.Add(pt)
			// log.Printf("%d: %s -> %d (%s)\n%s", p.ID, names[i], next.ID, next.At(), next)
			res = append(res, next)
		}
	}
	// log.Printf("Been: %v", bad)
	return res
}

type Map struct {
	Alt   map[Point]int
	Start Point
	End   Point
	Size  Point
}

func NewMap() *Map {
	return &Map{
		Alt:   map[Point]int{},
		Start: Point{},
		End:   Point{},
	}
}

func (m *Map) String() string {
	out := []string{}
	for y := 0; y < m.Size.Y; y++ {
		l := ""
		for x := 0; x < m.Size.X; x++ {
			l += fmt.Sprintf("%c", m.Alt[Point{x, y}]+'a')
		}
		out = append(out, l)
	}

	return strings.Join(out, "\n")
}

func (m *Map) Route(s, e Point) *Path {
	paths := []*Path{NewPath(m, s)}
	best := map[Point]int{}

	for true {
		next := []*Path{}
		for _, p := range paths {
			for _, n := range p.NextSteps() {
				if n.At().Eq(e) {
					return n
				}
				l := len(n.Steps)
				if b, ok := best[n.At()]; ok && l >= b {
					// log.Printf("Dropping inefficient route %d(%d) to %s:\n%s", n.ID, l, n.At(), n)
					continue
				}
				best[n.At()] = l
				// log.Printf("Best path to %s is %d(%d)", n.At(), n.ID, l)
				next = append(next, n)
			}
		}
		ids := []int{}
		for _, n := range next {
			ids = append(ids, n.ID)
		}
		// log.Printf("%d potential paths: %v", len(next), ids)
		if len(next) == 0 {
			return nil
		}
		paths = next
		// log.Println()
	}

	return nil
}

func one(m *Map) *Path {
	return m.Route(m.Start, m.End)
}

func two(m *Map) *Path {
	return nil
}

func readFile(path string) *Map {
	res := common.AsStrings(common.ReadTransformedFile(path, common.IgnoreBlankLines))
	m := NewMap()
	m.Size = Point{len(res[0]), len(res)}
	for y, l := range res {
		for x, c := range l {
			p := Point{x, y}
			if c == 'S' {
				c = 'a'
				m.Start = p
			} else if c == 'E' {
				c = 'z'
				m.End = p
			}
			m.Alt[p] = int(c - 'a')
		}
	}

	return m
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
