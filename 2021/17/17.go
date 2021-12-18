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

func (c *coord) move(v velocity) {
	c.X += v.X
	c.Y += v.Y
}

type area struct {
	UL, LR coord
}

func (a area) contains(c *coord) bool {
	return a.UL.X <= c.X && a.LR.X >= c.X && a.LR.Y <= c.Y && a.UL.Y >= c.Y
}

func readFile(path string) area {
	var a area
	data := common.ReadTransformedFile(
		path,
		common.IgnoreBlankLines,
		common.SplitWords,
	)

	for _, w := range data[0].([]string) {
		if strings.HasPrefix(w, "x=") {
			w = strings.TrimSuffix(w, ",")
			ns := strings.Split(w[2:], "..")
			a.UL.X = common.MustInt(ns[0])
			a.LR.X = common.MustInt(ns[1])
			if a.LR.X < a.UL.X {
				a.LR.X, a.UL.X = a.UL.X, a.LR.X
			}
		} else if strings.HasPrefix(w, "y=") {
			ns := strings.Split(w[2:], "..")
			a.UL.Y = common.MustInt(ns[0])
			a.LR.Y = common.MustInt(ns[1])
			if a.UL.Y < a.LR.Y {
				a.UL.Y, a.LR.Y = a.LR.Y, a.UL.Y
			}
		}
	}

	return a
}

type velocity struct {
	X, Y int
}

func (v *velocity) String() string {
	return fmt.Sprintf("v(%d,%d)", v.X, v.Y)
}

func (v *velocity) tick() {
	if v.X > 0 {
		v.X--
	} else if v.X < 0 {
		v.X++
	}
	v.Y -= 1
}

func draw(path []coord, a area) {
	maxY, maxX := 0, 0
	been := make(map[coord]bool)
	for _, c := range path {
		been[c] = true
		if c.X > maxX {
			maxX = c.X
		}
		if c.Y > maxY {
			maxY = c.Y
		}
	}
	if maxX < a.LR.X {
		maxX = a.LR.X
	}

	for y := maxY; y >= a.LR.Y; y-- {
		line := ""
		for x := 0; x <= maxX; x++ {
			if been[coord{x, y}] {
				line += "*"
			} else if a.contains(&coord{x, y}) {
				line += "#"
			} else {
				line += "."
			}
		}
		log.Printf("%10d: %s", y, line)
	}
}

func zenith(t []coord) int {
	max := t[0].Y
	for _, c := range t {
		if max < c.Y {
			max = c.Y
		}
	}

	return max
}

func track(v velocity, a area) ([]coord, bool) {
	pos := &coord{0, 0}
	been := []coord{*pos}
	for pos.Y >= a.LR.Y && pos.X <= a.LR.X {
		pos.move(v)
		been = append(been, *pos)
		v.tick()
		if a.contains(pos) {
			return been, true
		}
	}

	return been, false
}

func target(a area) (velocity, []coord) {
	v := velocity{}

	if a.LR.Y > 0 {
		log.Fatalf("bad assumptions, %v", a)
	}

	// find the lowest vx that ends up in the area
	log.Printf("area: %v", a)
	for x := 0; x < a.LR.X; x++ {
		dx := (x + 1) * x / 2
		if a.UL.X <= dx && a.LR.X >= dx {
			v.X = x
			break
		}
	}
	if v.X == 0 {
		log.Fatalf("can't figure out vx!")
	}

	// lets assume ballistic path, so -Y should be a good starting point
	v.Y = -a.LR.Y
	for {
		if _, hit := track(v, a); hit {
			log.Printf("vy: %d works", v.Y)
			v.Y++
		} else {
			v.Y--
			break
		}
	}

	if been, hit := track(v, a); hit {
		return v, been
	} else {
		log.Fatalf("bad news, %v didn't work", v)
	}

	return v, nil
}

func targetAll(a area) []velocity {
	vs := []velocity{}

	// get the minX and maxY from part A
	v, _ := target(a)
	log.Printf("velocity ranges: %d <= x <= %d; %d <= y <= %d",
		v.X, a.LR.X, a.LR.Y, v.Y)
	for y := a.LR.Y; y <= v.Y; y++ {
		for x := v.X; x <= a.LR.X; x++ {
			tv := velocity{x, y}
			if _, hit := track(tv, a); hit {
				vs = append(vs, tv)
			}
		}
	}

	return vs
}

func main() {
	a := readFile("input.txt")
	v, track := target(a)
	draw(track, a)
	log.Printf("v = %v", v)
	log.Printf("zenith = %d", zenith(track))
	log.Printf("total = %d", len(targetAll(a)))
}
