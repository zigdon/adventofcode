package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type dumbo struct {
	// step: y, x -> energy
	Layout  map[int]*[][]int
	Flashes map[int]int
}

type coord struct {
	x, y int
}

func newDumbo() *dumbo {
	d := &dumbo{
		Layout:  make(map[int]*[][]int),
		Flashes: map[int]int{0: 0},
	}

	return d
}

func (d *dumbo) string(step int) string {
	out := ""
	layout, ok := d.Layout[step]
	if !ok {
		return "<empty>"
	}

	for _, l := range *layout {
		line := ""
		for _, n := range l {
			line += fmt.Sprintf("%d ", n)
		}
		out += line + "\n"
	}

	return out
}

func (d *dumbo) init(step int, data ...string) {
	layout := [][]int{}
	for _, l := range data {
		line := []int{}
		for _, o := range l {
			line = append(line, int(o-'0'))
		}
		layout = append(layout, line)
	}
	d.Layout[step] = &layout
}

func (d *dumbo) step(s int) int {
	if s == 0 {
		return 0
	}
	if f, ok := d.Flashes[s]; ok {
		return f
	}
	// log.Printf("calculating flashes at step %d", s)

	f := 0
	pastFlash, ok := d.Flashes[s-1]
	if ok {
		f += pastFlash
	} else {
		// log.Printf("calculating missing flashes for step %d", s-1)
		f += d.step(s - 1)
	}
	layout, ok := d.Layout[s-1]
	if !ok {
		log.Fatalf("missing layout at step %d", s-1)
	}

	// make a copy, add 1 to each
	nl := [][]int{}
	for _, l := range *layout {
		line := []int{}
		for _, n := range l {
			line = append(line, n+1)
		}
		nl = append(nl, line)
	}
	d.Layout[s] = &nl

	newFlashes := d.flash(s)
	f += newFlashes
	d.Flashes[s] = f

	// log.Printf("at step %d -> %d (+%d) flashes", s, f, newFlashes)
	return f
}

func (d *dumbo) flash(step int) int {
	f := 0
	layout := d.Layout[step]

	// loop:
	//   scan for >9 -> Add (x,y) to queue, once only
	//   if queue has no true entries, done
	flashed := make(map[coord]bool)
	cont := true
	for cont {
		for y, l := range *layout {
			for x, n := range l {
				c := coord{x, y}
				_, ok := flashed[c]
				if n > 9 && !ok {
					flashed[c] = true
				}
			}
		}

		//   scan queue, add 1 to surrounding numbers
		//   mark entry false
		toFlash := []coord{}
		for k, v := range flashed {
			if v {
				toFlash = append(toFlash, k)
				flashed[k] = false
			}
		}

		for _, c := range toFlash {
			f++
			for _, dy := range []int{-1, 0, 1} {
				for _, dx := range []int{-1, 0, 1} {
					if dx == 0 && dy == 0 {
						continue
					}
					nx, ny := c.x+dx, c.y+dy
					if nx < 0 || ny < 0 || nx >= len((*layout)[0]) || ny >= len(*layout) {
						continue
					}
					(*layout)[ny][nx]++
				}
			}
		}

		cont = len(toFlash) > 0
	}

	// any squid that was flashed, has their energy set to 0
	for k := range flashed {
		(*layout)[k.y][k.x] = 0
	}

	d.Layout[step] = layout
	// fmt.Printf("after flashing step %d:\n%s", step, d.string(step))

	return f
}

func (d *dumbo) isSynced(step int) bool {
	d.step(step)
	for _, l := range *d.Layout[step] {
		for _, n := range l {
			if n > 0 {
				return false
			}
		}
	}

	return true
}

func mustAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("can't convert %q to int: %v", n, err)
	}

	return n
}

func readFile(path string) *dumbo {
	data := common.ReadTransformedFile(
		path,
		common.Block,
		common.IgnoreBlankLines,
	)

	d := newDumbo()
	for _, in := range data {
		step := 0
		// first line is either a step indicator or just the data
		s := common.AsStrings(in.([]interface{}))
		if len(s) == 0 {
			continue
		}
		if strings.Contains(s[0], ":") {
			if i := strings.IndexAny(s[0], "0123456789"); i > -1 {
				num := s[0][i:]
				step = mustAtoi(num[:len(num)-1])
			}
			s = s[1:]
		}
		d.init(step, s...)
	}

	return d
}

func main() {
	data := readFile("input.txt")
	log.Printf("after 100 steps, %d flashes", data.step(100))

	s := 1
	for !data.isSynced(s) {
		s++
		log.Printf("checking sync at %d", s)
	}
}
