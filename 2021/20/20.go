package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type coord struct {
	X, Y int
}

func (c coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func (c coord) add(b coord) coord {
	return coord{c.X + b.X, c.Y + b.Y}
}

func (c coord) sub(b coord) coord {
	return coord{c.X - b.X, c.Y - b.Y}
}

func (c coord) eq(b coord) bool {
	return c.X == b.X && c.Y == b.Y
}

func (c coord) in(a, b coord) bool {
	return c.X >= a.X && c.X <= b.X && c.Y >= a.Y && c.Y <= b.Y
}

type image struct {
	Pixels   map[coord]bool
	Default  bool
	Min, Max coord
}

func newImage(data []string) *image {
	i := &image{Pixels: make(map[coord]bool)}
	var zx, zy bool
	for y, line := range data {
		for x, p := range line {
			if p == '#' {
				if x == 0 {
					zx = true
				}
				if y == 0 {
					zy = true
				}
				i.Pixels[coord{x, y}] = true
				if x > i.Max.X {
					i.Max.X = x
				}
				if y > i.Max.Y {
					i.Max.Y = y
				}
				if i.Min.X == 0 || i.Min.X > x {
					i.Min.X = x
				}
				if i.Min.Y == 0 || i.Min.Y > y {
					i.Min.Y = y
				}
			}
		}
	}
	if zx {
		i.Min.X = 0
	}
	if zy {
		i.Min.Y = 0
	}

	return i
}

func (i *image) String() string {
	out := []string{
		fmt.Sprintf("%s-%s, default: %v", i.Min, i.Max, i.Default),
		fmt.Sprintf("  %6d", i.Min.X),
	}
	def := "."
	if i.Default {
		def = "#"
	}
	out = append(out, "       "+strings.Repeat(def, i.Max.X-i.Min.X+1), "")
	for y := i.Min.Y; y <= i.Max.Y; y++ {
		line := fmt.Sprintf("%3d: %s ", y, def)
		for x := i.Min.Y; x <= i.Max.X; x++ {
			if i.Pixels[coord{x, y}] {
				line += "#"
			} else {
				line += "."
			}
		}
		line += " " + def
		out = append(out, line)
	}
	out = append(out, "", "       "+strings.Repeat(def, i.Max.X-i.Min.X+1))

	return strings.Join(out, "\n")
}

func (i *image) count() int {
	c := 0
	for _, v := range i.Pixels {
		if v {
			c++
		}
	}

	return c
}

func (i *image) set(x, y int) {
	i.Pixels[coord{x, y}] = true
	if i.Max.X < x {
		i.Max.X = x
	}
	if i.Max.Y < y {
		i.Max.Y = y
	}
	if i.Min.X > x {
		i.Min.X = x
	}
	if i.Min.Y > y {
		i.Min.Y = y
	}
}

func (i *image) eq(ib *image) bool {
	if i.Default != ib.Default {
		log.Printf("different defaults: %v vs %v", i.Default, ib.Default)
		return false
	}
	offset := i.Min.sub(ib.Min)
	log.Printf("offset: %s", offset)
	for k, v := range ib.Pixels {
		if i.Pixels[k.add(offset)] != v {
			log.Printf("bad pixel added %s:%v vs %s:%v", k.add(offset), i.Pixels[k.add(offset)], k, v)
			return false
		}
	}
	for k, v := range i.Pixels {
		if ib.Pixels[k.sub(offset)] != v {
			log.Printf("bad pixel removed %s:%v vs %s:%v", k, v, k.sub(offset), ib.Pixels[k.sub(offset)])
			return false
		}
	}

	return true
}

func (i *image) get(c coord) bool {
	if i.Pixels[c] {
		return true
	}
	if c.in(i.Min, i.Max) {
		return false
	}
	return i.Default

}

func (i *image) getVal(x, y int) int {
	// str := ""
	out := 0
	for _, vy := range []int{y - 1, y, y + 1} {
		for _, vx := range []int{x - 1, x, x + 1} {
			out *= 2
			if i.get(coord{vx, vy}) {
				// str += "#"
				out++
			} else {
				// str += "."
			}
		}
	}

	// log.Printf("(%d,%d): %q = %d", x, y, str, out)

	return out
}

func (i *image) apply(a *algo) *image {
	ni := newImage([]string{})
	for y := i.Min.Y - 1; y <= i.Max.Y+1; y++ {
		for x := i.Min.X - 1; x <= i.Max.X+1; x++ {
			n := i.getVal(x, y)
			if a.get(n) {
				ni.set(x, y)
			}
		}
	}
	if i.Default {
		ni.Default = a.get(511)
	} else {
		ni.Default = a.get(0)
	}

	return ni
}

func (i *image) enhance(a *algo, n int) *image {
	for ; n > 0; n-- {
		i = i.apply(a)
	}
	return i
}

type algo struct {
	Bits []bool
}

func newAlgo(data string) *algo {
	a := &algo{Bits: []bool{}}
	for _, b := range data {
		a.Bits = append(a.Bits, b == '#')
	}

	return a
}

func (a *algo) get(n int) bool {
	return a.Bits[n]
}

func readFile(path string) (*algo, *image) {
	data := common.AsStrings(common.ReadTransformedFile(
		path, common.IgnoreBlankLines))

	return newAlgo(data[0]), newImage(data[1:])
}

func main() {
	algo, img := readFile("input.txt")
	log.Printf("original image:\n%s", img)
	ni := img.enhance(algo, 2)
	log.Printf("After 2:\n%s", ni)
	log.Printf("count: %d", ni.count())
	ni = img.enhance(algo, 50)
	log.Printf("After 50:\n%s", ni)
	log.Printf("count: %d", ni.count())
}
