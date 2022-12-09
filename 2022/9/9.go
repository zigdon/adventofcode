package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type Point struct {
	x, y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

func (p Point) Eq(p2 Point) bool {
	return p.x == p2.x && p.y == p2.y
}
func (p Point) Sub(p2 Point) Point {
	return Point{p.x - p2.x, p.y - p2.y}
}
func (p Point) Add(p2 Point) Point {
	return Point{p.x + p2.x, p.y + p2.y}
}
func (p Point) Mag() int {
	x := p.x
	if x < 0 {
		x = -x
	}
	y := p.y
	if y < 0 {
		y = -y
	}
	if x > y {
		return x
	}
	return y
}

type Rope struct {
	Head, Tail, Max, Min Point
	Seen                 map[Point]bool
}

func NewRope() *Rope {
	return &Rope{
		Head: Point{},
		Tail: Point{},
		Max:  Point{},
		Min:  Point{},
		Seen: map[Point]bool{
			Point{}: true,
		},
	}
}

func (r *Rope) String() string {
	res := []string{}
	l := "     "
	for x := r.Min.x; x <= r.Max.x; x++ {
		l += fmt.Sprintf("%-2d ", x)
	}
	res = append(res, l)
	for y := r.Max.y; y >= r.Min.y; y-- {
		l = fmt.Sprintf("%3d: ", y)
		for x := r.Min.x; x <= r.Max.x; x++ {
			p := Point{x, y}
			var c string
			switch {
			case r.Head.Eq(p):
				c = "H"
			case r.Tail.Eq(p):
				c = "T"
			case x == 0 && y == 0:
				c = "s"
			case r.Seen[p]:
				c = "#"
			default:
				c = "."
			}
			l += fmt.Sprintf("%s  ", c)
		}
		res = append(res, l)
	}

	return strings.Join(res, "\n")
}

func (r *Rope) Nudge(dx, dy int) {
	r.Head.x += dx
	r.Head.y += dy
	if r.Head.x > r.Max.x {
		r.Max.x = r.Head.x
	}
	if r.Head.y > r.Max.y {
		r.Max.y = r.Head.y
	}
	if r.Head.x < r.Min.x {
		r.Min.x = r.Head.x
	}
	if r.Head.y < r.Min.y {
		r.Min.y = r.Head.y
	}
	log.Printf("head -> %s", r.Head)
	diff := r.Head.Sub(r.Tail)
	if diff.Mag() <= 1 {
		return
	}
	var mx, my int
	if diff.x > 1 {
		mx = 1
	} else if diff.x < -1 {
		mx = -1
	}
	if diff.y > 1 {
		my = 1
	} else if diff.y < -1 {
		my = -1
	}
	r.Tail.Add(Point{mx, my})
	log.Printf("tail (%d, %d) -> %s", mx, my, r.Tail)
	r.Seen[r.Tail] = true
}

func (r *Rope) Move(i Inst) {
	log.Print(i.String())
	for n := 0; n < i.Dist; n++ {
		r.Nudge(i.Dx, i.Dy)
	}
}

type Inst struct {
	Dx, Dy, Dist int
}

func (i Inst) String() string {
	return fmt.Sprintf("I[%d,%d x %d]", i.Dx, i.Dy, i.Dist)
}

func NewInst(d string, q int) Inst {
	i := Inst{Dist: q}
	if d == "U" {
		i.Dy = 1
	} else if d == "D" {
		i.Dy = -1
	} else if d == "R" {
		i.Dx = 1
	} else if d == "L" {
		i.Dx = -1
	}
	return i
}

func one(inst []Inst) int {
	r := NewRope()
	for _, i := range inst {
		r.Move(i)
		log.Printf("\n%s\n", r)
	}
	res := 0
	for _, v := range r.Seen {
		if v {
			res++
		}
	}
	return res
}

func two(i []Inst) int {
	return 0
}

func readFile(path string) ([]Inst, error) {
	i := common.AsStrings(common.ReadTransformedFile(path, common.IgnoreBlankLines))

	inst := []Inst{}
	for _, l := range i {
		bits := strings.Split(l, " ")
		inst = append(inst, NewInst(bits[0], common.MustInt(bits[1])))
	}

	return inst, nil
}

func main() {
	i, err := readFile(os.Args[1])
	if err != nil {
		log.Fatalf("%v", err)
	}

	res := one(i)
	fmt.Printf("%v\n", res)

	res = two(i)
	fmt.Printf("%v\n", res)
}
