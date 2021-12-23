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

func (c coord) sub(b coord) coord {
	return coord{c.X - b.X, c.Y - b.Y, c.Z - b.Z}
}

func (c coord) add(b coord) coord {
	return coord{c.X + b.X, c.Y + b.Y, c.Z + b.Z}
}

func (c coord) eq(b coord) bool {
	return c.X == b.X && c.Y == b.Y && c.Z == b.Z
}

func (c coord) lt(b coord) bool {
	return c.X < b.X || c.Y < b.Y || c.Z < b.Z
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

func (c coord) turn(r coord) coord {
	nc := coord{c.X, c.Y, c.Z}
	switch r.X {
	case 1:
		nc.Y, nc.Z = nc.Z, -nc.Y
	case 2:
		nc.Y, nc.Z = -nc.Y, -nc.Z
	case 3:
		nc.Y, nc.Z = -nc.Z, nc.Y
	}
	switch r.Y {
	case 1:
		nc.X, nc.Z = nc.Z, -nc.X
	case 2:
		nc.X, nc.Z = -nc.X, -nc.Z
	case 3:
		nc.X, nc.Z = -nc.Z, nc.X
	}
	switch r.Z {
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

func (p pair) String() string {
	return fmt.Sprintf("[%s, %s]", p.a, p.b)
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

func (s *scanner) shift(c coord) {
	var nb = make(map[coord]bool)
	for k := range s.Beacons {
		nb[k.add(c)] = true
	}
	s.Beacons = nb

	var nd = make(map[coord]pair)
	for k, v := range s.Deltas {
		nd[k] = pair{v.a.add(c), v.b.add(c)}
	}
	s.Deltas = nd

	s.Min = s.Min.add(c)
	s.Max = s.Max.add(c)
	s.Origin = s.Origin.add(c)
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
			if b.distance(coord{0, 0, 0}) < a.distance(coord{0, 0, 0}) {
				continue
			}
			delta := b.sub(a)
			if _, ok := s.Deltas[delta]; ok {
				log.Printf("collision at %s", delta)
			}
			s.Deltas[delta] = pair{a, b}
		}
	}
}

func (s *scanner) turn(r coord) {
	s.Orientation = s.Orientation.turn(r)

	nb := make(map[coord]bool)
	for b := range s.Beacons {
		nb[b.turn(r)] = true
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
func (s *scanner) try24(f func(*scanner, coord)) {
	var transformation coord
	for _, r := range []coord{
		{0, 0, 0},
		{1, 0, 0}, {1, 0, 0}, {1, 0, 0},
		{1, 1, 0}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0},
		{1, 1, 0}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0},
		{1, 1, 0}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0},
		{1, 2, 1}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0},
		{1, 2, 0}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0},
	} {
		s.turn(r)
		transformation = transformation.add(r)
		f(s, transformation)
	}

	s.turn(coord{2, 0, 1})
}

func (s *scanner) align(sb *scanner) (int, coord, coord, []coord) {
	var bestMatch int
	var bestOri, bestShift, turns coord
	matches := make(map[coord]bool)

	matchFunc := func(sb *scanner, rots coord) {
		// Line up pairs of beacons in the original with a pair from the overlay
		found := make(map[pair]pair)
		count := make(map[coord]bool)
		// find matching deltas, add the two ends as possible matches
		for c, over := range sb.Deltas {
			if f, ok := s.Deltas[c]; ok {
				found[f] = over
				count[f.a] = true
				count[f.b] = true
			}
		}

		if len(count) > bestMatch {
			bestShift = coord{0, 0, 0}
			bestMatch = len(count)
			bestOri = sb.Orientation
			turns = rots
			for k := range matches {
				delete(matches, k)
			}

			// each found pair shows 4 possible shifts, but we know the
			// orientation is the same, so the leftmost of each pair matches.
			for ref, over := range found {
				if ref.b.lt(ref.a) {
					ref = pair{ref.b, ref.a}
				}
				if over.b.lt(over.a) {
					over = pair{over.b, over.a}
				}
				if bestShift.isEmpty() {
					bestShift = ref.a.sub(over.a)
					log.Printf("setting bestShift: %s", bestShift)
				}
				matches[ref.a] = true
				matches[ref.b] = true
			}
		}
	}

	sb.try24(matchFunc)
	sb.turn(turns)
	if sb.Origin.isEmpty() && sb.ID != 0 {
		sb.shift(bestShift)
	}

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
