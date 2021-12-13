package main

import (
	"log"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type coord struct {
	X, Y int
}

type paper struct {
	Marks map[coord]bool
	Size  coord
}

func (p *paper) fold(is []inst) {
	gone := []coord{}
	for _, i := range is {
		var f func(coord) coord
		log.Printf("Before folding %v:\n%s", i, p.draw())
		if i.Dir == "y" {
			p.Size.Y = i.Pos
			f = func(c coord) coord {
				if c.Y <= i.Pos {
					return c
				}
				return coord{c.X, 2*i.Pos - c.Y}
			}
		} else {
			p.Size.X = i.Pos
			f = func(c coord) coord {
				if c.X <= i.Pos {
					return c
				}
				return coord{2*i.Pos - c.X, c.Y}
			}
		}

		for c, there := range p.Marks {
			if !there {
				continue
			}
			nc := f(c)
			if nc != c {
				p.Marks[c] = false
				p.Marks[nc] = true
				gone = append(gone, c)
			}
		}
		log.Printf("After folding %v:\n%s", i, p.draw())
	}

	for _, c := range gone {
		delete(p.Marks, c)
	}
}

func (p *paper) count() int {
	res := 0
	for _, v := range p.Marks {
		if v {
			res++
		}
	}

	return res
}

func (p *paper) draw() string {
	log.Printf("drawing %v", p.Size)
	out := ""
	mY, mX := p.Size.Y, p.Size.X
	if mY > 100 {
		mY = 100
	}
	if mX > 100 {
		mX = 100
	}
	for y := 0; y <= mY; y++ {
		for x := 0; x <= mX; x++ {
			if p.Marks[coord{x, y}] {
				out += "#"
			} else {
				out += "."
			}
		}
		out += "\n"
	}

	return out
}

type inst struct {
	Dir string
	Pos int
}

func readFile(path string) (*paper, []inst) {
	p := &paper{Marks: make(map[coord]bool)}
	i := []inst{}

	data := common.ReadTransformedFile(
		path,
		common.Split(","),
		common.Block,
		common.IgnoreBlankLines,
	)

	dots := data[0].([]interface{})
	folds := data[1].([]interface{})
	for _, c := range dots {
		l := common.AsStrings(c)
		co := coord{common.MustInt(l[0]), common.MustInt(l[1])}
		p.Marks[co] = true
		if co.X > p.Size.X {
			p.Size.X = co.X
		}
		if co.Y > p.Size.Y {
			p.Size.Y = co.Y
		}
	}

	for _, f := range folds {
		in := inst{}
		s := f.([]string)[0]
		if strings.Contains(s, "x=") {
			in.Dir = "x"
		} else {
			in.Dir = "y"
		}

		in.Pos = common.MustInt(s[strings.Index(s, "=")+1:])

		i = append(i, in)
	}

	return p, i
}

func main() {
	p, ins := readFile("input.txt")
	p.fold([]inst{ins[0]})
	log.Printf("%d marks after folding.", p.count())
	p.fold(ins[1:])
}
