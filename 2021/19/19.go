package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type coord struct {
	X, Y, Z int
}

func (c coord) String() string {
	return fmt.Sprintf("(%d,%d,%d)", c.X, c.Y, c.Z)
}

func (c coord) eq(b coord) bool {
	return c.X == b.X && c.Y == b.Y && c.Z == b.Z
}

func (c coord) sub(b coord) coord {
	return coord{c.X - b.X, c.Y - b.Y, c.Z - b.Z}
}

func (c coord) add(b coord) coord {
	return coord{c.X + b.X, c.Y + b.Y, c.Z + b.Z}
}

func (c coord) max(b coord) coord {
	if b.X > c.X {
		c.X = b.X
	}
	if b.Y > c.Y {
		c.Y = b.Y
	}
	if b.Z > c.Z {
		c.Z = b.Z
	}

	return c
}

func (c coord) min(b coord) coord {
	if b.X < c.X {
		c.X = b.X
	}
	if b.Y < c.Y {
		c.Y = b.Y
	}
	if b.Z < c.Z {
		c.Z = b.Z
	}

	return c
}

func (c coord) isEmpty() bool {
	return c.X == 0 && c.Y == 0 && c.Z == 0
}

func (c coord) distance(d coord) float64 {
	return math.Sqrt(
		math.Pow(float64(c.X-d.X), 2) +
			math.Pow(float64(c.Y-d.Y), 2) +
			math.Pow(float64(c.Z-d.Z), 2))
}

func (c coord) turn(x, y, z int) coord {
	nc := coord{c.X, c.Y, c.Z}
	switch x {
	case 1:
		nc.Y, nc.Z = nc.Z, -nc.Y
	case 2:
		nc.Y, nc.Z = -nc.Y, -nc.Z
	case 3:
		nc.Y, nc.Z = -nc.Z, nc.Y
	}
	switch y {
	case 1:
		nc.X, nc.Z = nc.Z, -nc.X
	case 2:
		nc.X, nc.Z = -nc.X, -nc.Z
	case 3:
		nc.X, nc.Z = -nc.Z, nc.X
	}
	switch z {
	case 1:
		nc.X, nc.Y = nc.Y, -nc.X
	case 2:
		nc.X, nc.Y = -nc.X, -nc.Y
	case 3:
		nc.X, nc.Y = -nc.Y, nc.X
	}

	return nc
}

type pair struct {
	a, b coord
}

type scanner struct {
	ID          int
	Origin      coord
	Orientation coord // 2 forward, 1 up
	Beacons     map[coord]bool
	Deltas      map[coord]pair
	Min, Max    coord
}

func newScanner(n int, bs []string) *scanner {
	s := &scanner{
		ID:          n,
		Beacons:     make(map[coord]bool),
		Deltas:      make(map[coord]pair),
		Orientation: coord{2, 0, 1},
	}
	for _, l := range bs {
		s.add(l)
	}

	s.makeDeltas()

	return s
}

func (s *scanner) String() string {
	out := []string{
		fmt.Sprintf("--- Scanner %d ---", s.ID),
		fmt.Sprintf("  Origin: %s", s.Origin),
	}
	if len(s.Beacons) > 0 {
		out = append(out, "  Beacons:")
	}
	for c := range s.Beacons {
		out = append(out, "    "+c.String())
	}
	if len(s.Deltas) > 0 {
		out = append(out, "  Deltas:")
	}
	for c, t := range s.Deltas {
		out = append(out, fmt.Sprintf("    %s (%s, %s)", c, t.a, t.b))
	}

	return strings.Join(out, "\n")
}

func (s *scanner) makeDeltas() {
	// calculate the deltas between all the beacons

	s.Deltas = make(map[coord]pair)

	bs := []coord{}
	for b := range s.Beacons {
		bs = append(bs, b)
	}
	for i := range bs {
		for j := range bs {
			if i == j {
				continue
			}

			a := bs[i]
			b := bs[j]
			delta := b.sub(a)
			if _, ok := s.Deltas[delta]; ok {
				log.Printf("collision at %s", delta)
			}
			s.Deltas[delta] = pair{a, b}
		}
	}
}

func (s *scanner) turn(x, y, z int) {
	s.Orientation = s.Orientation.turn(x, y, z)

	nb := make(map[coord]bool)
	for b := range s.Beacons {
		nb[b.turn(x, y, z)] = true
	}
	s.Beacons = nb

	s.makeDeltas()
}

func (s *scanner) add(b string) {
	ns := strings.Split(b, ",")
	c := coord{
		common.MustInt(ns[0]),
		common.MustInt(ns[1]),
		common.MustInt(ns[2]),
	}
	s.Beacons[c] = true
	s.Min = s.Min.min(c)
	s.Max = s.Max.max(c)
}

// run f in all 24 possible orientations, ensuring s ends up at the same
// orientation as it started.
func (s *scanner) try24(f func(*scanner)) {
	for _, r := range []coord{
		{0, 0, 0},
		{1, 0, 0}, {1, 0, 0}, {1, 0, 0},
		{1, 1, 0}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0},
		{1, 1, 0}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0},
		{1, 1, 0}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0},
		{1, 2, 1}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0},
		{1, 2, 0}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0},
	} {
		s.turn(r.X, r.Y, r.Z)
		f(s)
	}

	s.turn(2, 0, 1)
}

func (s *scanner) align(sb *scanner) (int, coord, coord, []coord) {
	var bestMatch int
	var bestOri, bestShift coord
	matches := make(map[coord]coord)

	matchFunc := func(sb *scanner) {
		found := make(map[coord]pair)
		for c, orig := range sb.Deltas {
			if f, ok := s.Deltas[c]; ok {
				found[f.a] = pair{s.Deltas[c].a, orig.a}
				found[f.b] = pair{s.Deltas[c].b, orig.b}
			}
		}

		if len(found) > bestMatch {
			bestMatch = len(found)
			bestOri = sb.Orientation
			for k := range matches {
				delete(matches, k)
			}

			// found shows delta
			for _, v := range found {
				if bestShift.isEmpty() {
					bestShift = v.a.sub(v.b).add(s.Origin)
					if sb.Origin.isEmpty() && sb.ID != 0 {
						sb.Origin = bestShift
					} else if !sb.Origin.eq(bestShift) {
						log.Fatalf("moving origin of %d from %s to %s", sb.ID, sb.Origin, bestShift)
					}
					log.Printf("shift of %s relative to %s -> %s", v.a.sub(v.b), s.Origin, sb.Origin)
				} else {
					if bestShift != v.a.sub(v.b).add(s.Origin) {
						log.Fatalf("inconsistent shifting: both %s and %s", bestShift, v.a.sub(v.b))
					}
				}
				matches[v.b.add(sb.Origin)] = v.a.add(s.Origin)
			}
		}
	}

	sb.try24(matchFunc)

	res := []coord{}
	for k := range matches {
		res = append(res, k)
	}

	return bestMatch, bestOri, bestShift, res
}

func readFile(path string) []*scanner {
	data := common.ReadTransformedFile(
		path,
		common.Block,
		common.IgnoreBlankLines)

	ss := []*scanner{}
	for _, d := range data {
		ls := d.([]interface{})
		n := common.MustInt(strings.Split(ls[0].(string), " ")[2])
		s := newScanner(n, nil)

		for _, b := range ls[1:] {
			s.add(b.(string))
		}

		s.makeDeltas()
		ss = append(ss, s)
	}

	return ss
}

func main() {
	fmt.Println("vim-go")
}
