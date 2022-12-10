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

func (p *Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

func (p *Point) Set(p2 *Point) {
	p.x = p2.x
	p.y = p2.y
}
func (p *Point) Move(p2 *Point) {
	// log.Printf("Move p1=%v + p2=%v", p, p2)
	p.x += p2.x
	p.y += p2.y
	// log.Printf("-> %v", p)
}
func (p *Point) Eq(p2 *Point) bool {
	return p.x == p2.x && p.y == p2.y
}
func (p *Point) Sub(p2 *Point) *Point {
	return &Point{p.x - p2.x, p.y - p2.y}
}
func (p *Point) Add(p2 *Point) *Point {
	return &Point{p.x + p2.x, p.y + p2.y}
}
func (p *Point) Mag() int {
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
	Knots    []*Point
	Max, Min *Point
	Seen     map[Point]bool
}

func (r *Rope) MarkSeen(p *Point) {
	r.Seen[*p] = true
	// log.Printf("Adding %s to seen in %v", p, &r)
	// log.Printf("seen: %v", r.Seen)
}

func (r *Rope) Tail() *Point {
	return r.Knots[len(r.Knots)-1]
}

func (r *Rope) Head() *Point {
	return r.Knots[0]
}

func NewRope(knots int) *Rope {
	r := &Rope{
		Max:  &Point{},
		Min:  &Point{},
		Seen: make(map[Point]bool),
	}
    // log.Printf("Creating new rope with %d knots -> %v", knots, &r)
	r.MarkSeen(&Point{})
	for i := 0; i < knots; i++ {
		r.Knots = append(r.Knots, &Point{})
	}

	return r
}

func (r *Rope) String() string {
	res := []string{"\n"}
	l := "     "
	for x := r.Min.x; x <= r.Max.x; x++ {
		l += fmt.Sprintf("%-2d ", x)
	}
	res = append(res, l)
	for y := r.Max.y; y >= r.Min.y; y-- {
		l = fmt.Sprintf("%3d: ", y)
		for x := r.Min.x; x <= r.Max.x; x++ {
			p := &Point{x, y}
			c := "."
			for i, k := range r.Knots {
				if k.Eq(p) {
                    if i == 0 {
                      c = "H"
                    } else if i == len(r.Knots)-1 {
                      c = "T"
                    } else {
                      c = fmt.Sprintf("%d", i)
                    }
					break
				}
			}
			if c == "." {
				if x == 0 && y == 0 {
					c = "s"
				} else if r.Seen[*p] {
					c = "#"
				}
			}
			l += fmt.Sprintf("%s  ", c)
		}
		res = append(res, l)
	}

	return strings.Join(res, "\n")
}

func (r *Rope) Nudge(i, dx, dy int) {
	// log.Printf("Nudge %d: %s + (%d,%d)", i, r.Knots[i], dx, dy)
	r.Knots[i].Move(&Point{dx, dy})
	if i == 0 {
		if r.Head().x > r.Max.x {
			r.Max.x = r.Head().x
		}
		if r.Head().y > r.Max.y {
			r.Max.y = r.Head().y
		}
		if r.Head().x < r.Min.x {
			r.Min.x = r.Head().x
		}
		if r.Head().y < r.Min.y {
			r.Min.y = r.Head().y
		}
	}
	// log.Printf("head -> %s", r.Head())
	if i+1 >= len(r.Knots) {
		// log.Printf("tail + (%d, %d) -> %s", dx, dy, r.Tail())
		r.MarkSeen(r.Tail())
		return
	}
	diff := r.Knots[i].Sub(r.Knots[i+1])
	if diff.Mag() <= 1 {
		return
	}
	var mx, my int
	if diff.x > 0 {
		mx = 1
	} else if diff.x < 0 {
		mx = -1
	}
	if diff.y > 0 {
		my = 1
	} else if diff.y < 0 {
		my = -1
	}
	r.Nudge(i+1, mx, my)
}

func (r *Rope) Move(i Inst, debug bool) {
    if debug { log.Print(i.String()) }
	for n := 0; n < i.Dist; n++ {
        if debug {
          log.Printf("\n=> %s (%d)", i, n)
          log.Printf("%s", r)
        }
		r.Nudge(0, i.Dx, i.Dy)
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
	r := NewRope(2)
	for _, i := range inst {
		r.Move(i, false)
	}
	res := 0
	for _, v := range r.Seen {
		if v {
			res++
		}
	}
	return res
}

func two(inst []Inst) int {
	r := NewRope(10)
	for _, i := range inst {
		r.Move(i, false)
	}
	res := 0
	for _, v := range r.Seen {
		if v {
			res++
		}
	}
	return res
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
