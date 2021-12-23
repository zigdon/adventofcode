package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type coords []coord

func (bs coords) less(i, j int) bool {
	return bs[i].lt(bs[j])
}

func (bs coords) String() string {
	out := []string{}
	for _, b := range bs {
		out = append(out, b.String())
	}

	return strings.Join(out, "\n")
}

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

func (c coord) mod(n int) coord {
	return coord{c.X % n, c.Y % n, c.Z % n}
}

func (c coord) eq(b coord) bool {
	return c.X == b.X && c.Y == b.Y && c.Z == b.Z
}

func (c coord) lt(b coord) bool {
	if c.X < b.X {
		return true
	} else if c.X > b.X {
		return false
	}
	if c.Y < b.Y {
		return true
	} else if c.Y > b.Y {
		return false
	}
	if c.Z < b.Z {
		return true
	} else if c.Z > b.Z {
		return false
	}

	return false
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

func (c coord) manhattan(d coord) int {
	return int(
		math.Abs(float64(c.X-d.X)) +
			math.Abs(float64(c.Y-d.Y)) +
			math.Abs(float64(c.Z-d.Z)))
}

func (c coord) turn(r coord) coord {
	r = r.mod(4)
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
	A, B coord
}

func (p pair) String() string {
	return fmt.Sprintf("[%s, %s]", p.A, p.B)
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
		out = append(out, fmt.Sprintf("    %s (%s, %s)", c, t.A, t.B))
	}

	return strings.Join(out, "\n")
}

func (s *scanner) listBeacons() string {
	out := coords{}
	for k := range s.Beacons {
		out = append(out, k)
	}

	sort.Slice(out, out.less)

	return out.String()
}

func (s *scanner) shift(c coord) {
	var nb = make(map[coord]bool)
	for k := range s.Beacons {
		nb[k.add(c)] = true
	}
	s.Beacons = nb

	var nd = make(map[coord]pair)
	for k, v := range s.Deltas {
		nd[k] = pair{v.A.add(c), v.B.add(c)}
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
			if p, ok := s.Deltas[delta]; ok {
				log.Printf("collision at %s: %s and %s", delta, p, pair{a, b})
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
// orientation as it started. If f returns true, stop trying.
func (s *scanner) try24(f func(*scanner) bool) {
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
		if f(s) {
			return
		}
	}

	s.turn(coord{2, 0, 1})
}

func (s *scanner) align(sb *scanner, req int) []coord {
	var bestMatch int
	var bestShift, bestOrientation coord
	matches := make(map[coord]bool)

	// Line up pairs of beacons in the original with a pair from the overlay
	matchFunc := func(sb *scanner) bool {
		found := make(map[pair]pair)
		count := make(map[coord]bool)

		// find matching deltas, add the two ends as possible matches
		for delta, over := range sb.Deltas {
			if ref, ok := s.Deltas[delta]; ok {
				found[ref] = over
				count[ref.A] = true
				count[ref.B] = true
			}
		}

		if len(count) > bestMatch {
			bestOrientation = sb.Orientation
			bestShift = coord{0, 0, 0}
			bestMatch = len(count)
			for k := range matches {
				delete(matches, k)
			}

			// each found pair shows 4 possible shifts, but we know the
			// orientation is the same, so the leftmost of each pair matches.
			for ref, over := range found {
				if ref.B.lt(ref.A) {
					ref = pair{ref.B, ref.A}
				}
				if over.B.lt(over.A) {
					over = pair{over.B, over.A}
				}
				if bestShift.isEmpty() {
					bestShift = ref.A.sub(over.A)
				}
				matches[ref.A] = true
				matches[ref.B] = true
			}
		}
		return false // try all 24
	}

	sb.try24(matchFunc)

	res := []coord{}
	if bestMatch >= req {
		// log.Printf("adjusting #%d:\n%s", sb.ID, sb.listBeacons())
		// try24 returns the scanner to its origin, so turn it until it's at bestOrientation
		sb.turnTo(bestOrientation)
		// log.Printf("turned #%d %v:\n%s", sb.ID, turns, sb.listBeacons())
		if sb.Origin.isEmpty() && sb.ID != 0 {
			sb.shift(bestShift)
			// log.Printf("shifted #%d %v:\n%s", sb.ID, bestShift, sb.listBeacons())
		}

		for k := range matches {
			res = append(res, k)
		}
	}
	return res
}

func (s *scanner) turnTo(r coord) {
	s.try24(func(sb *scanner) bool {
		if !sb.Orientation.eq(r) {
			return false
		}
		return true
	})
}

func (s *scanner) merge(sb *scanner) {
	nb := coords{}
	for c := range sb.Beacons {
		if !s.Beacons[c] {
			nb = append(nb, c)
			s.Beacons[c] = true
		}
	}
	if len(nb) > 0 {
		sort.Slice(nb, nb.less)
		// log.Printf("Adding from #%d to #%d:\n%s", sb.ID, s.ID, nb)
	}

	s.makeDeltas()
}

func turnTo(o coord) coord {
	var res coord
	switch {
	case o.X == -2:
		res.Z += 2
	case o.Y == 2:
		res.Z += 1
	case o.Y == -2:
		res.Z += 3
	case o.Z == 2:
		res.Y += 1
	case o.Z == -2:
		res.Y += 3
	}
	switch {
	case o.X == 1:
		res.Z += 2
	case o.Y == 2:
		res.Z += 1
	case o.Y == -2:
		res.Z += 3
	case o.Z == 2:
		res.Y += 1
	case o.Z == -2:
		res.Y += 3
	}

	return res
}

func solve(ss []*scanner) (*scanner, bool) {
	var solved = map[int]bool{ss[0].ID: true}
	cont := true

	for cont {
		cont = false
		for _, a := range ss {
			if !solved[a.ID] {
				continue
			}
			for _, b := range ss {
				if solved[b.ID] {
					continue
				}
				got := a.align(b, 12)
				// log.Printf("Trying to align %d to %d: got %d", b.ID, a.ID, len(got))
				if len(got) < 12 {
					continue
				}
				log.Printf("solved #%d, shifted %s!", b.ID, b.Origin)
				solved[b.ID] = true
				cont = true
				break
			}
		}
	}

	good := true
	for _, s := range ss {
		log.Printf("%d: %v (%s)", s.ID, solved[s.ID], s.Origin)
		if !solved[s.ID] {
			good = false
		}
		if s.ID != ss[0].ID {
			ss[0].merge(s)
		}
	}

	return ss[0], good
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
	data := readFile("input.txt")
	res, ok := solve(data)
	if !ok {
		log.Printf("failed to solve!")
	}
	log.Printf("got %d beacons", len(res.Beacons))

	max := 0
	for _, a := range data {
		for _, b := range data {
			d := a.Origin.manhattan(b.Origin)
			if d > max {
				log.Printf("#%d and #%d are %d apart (%s - %s)", a.ID, b.ID, d, a.Origin, b.Origin)
				max = d
			}
		}
	}
}
