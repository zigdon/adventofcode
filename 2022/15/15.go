package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/zigdon/adventofcode/common"
)

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("[%d,%d]", p.X, p.Y)
}
func (p Point) Dist(p2 Point) int {
	return int(math.Abs(float64(p2.X-p.X)) + math.Abs(float64(p2.Y-p.Y)))
}

type Object int

const (
	Unknown Object = iota
	Air
	Sensor
	Beacon
)

func (o Object) String() string {
	return []string{"?", ".", "S", "B"}[o]
}

type Range struct {
	From, To int
}

func NewRange(f, t int) *Range {
	if t < f {
		f, t = t, f
	}
	return &Range{f, t}
}
func (r *Range) String() string {
	return fmt.Sprintf("R[%d:%d]", r.From, r.To)
}
func (r *Range) Contains(x int) bool {
	return x >= r.From && x <= r.To
}
func (r *Range) Merge(r2 *Range) bool {
	if r.To+1 < r2.From || r.From-1 > r2.To {
		// log.Printf("Not merging %s and %s", r, r2)
		return false
	}
	// log.Printf("Merging %s and %s", r, r2)
	r.From = int(math.Min(float64(r.From), float64(r2.From)))
	r.To = int(math.Max(float64(r.To), float64(r2.To)))
	// log.Printf("-> %s", r)
	return true
}
func (r *Range) Copy() *Range {
	return &Range{r.From, r.To}
}

type Field struct {
	Objs     map[Point]Object
	Min, Max *Point
	Gaps     map[int][]*Range
}

func NewField() *Field {
	return &Field{
		Objs: make(map[Point]Object),
		Gaps: make(map[int][]*Range),
	}
}
func (f *Field) String() string {
	ruler := ""
	for n := 0; n < f.Max.X; n += 10 {
		ruler += fmt.Sprintf("%2d        ", n)
	}
	res := []string{
		"",
		strings.Repeat(" ", -f.Min.X+5) + ruler,
	}
	for y := f.Min.Y; y <= f.Max.Y; y++ {
		l := fmt.Sprintf("%3d: ", y)
		rs := f.Gaps[y]
		for x := f.Min.X; x <= f.Max.X; x++ {
			o, ok := f.Objs[Point{x, y}]
			if ok {
				l += o.String()
				continue
			}
			found := false
			for _, r := range rs {
				if r.Contains(x) {
					l += "#"
					found = true
					break
				}
			}
			if !found {
				l += "."
			}
		}
		res = append(res, l)
	}
	return strings.Join(res, "\n")
}
func (f *Field) UpdateSize(p Point) {
	if f.Min == nil {
		f.Min = &Point{p.X, p.Y}
	}
	if f.Max == nil {
		f.Max = &Point{p.X, p.Y}
	}
	if f.Min.X > p.X {
		f.Min.X = p.X
	}
	if f.Min.Y > p.Y {
		f.Min.Y = p.Y
	}
	if f.Max.X < p.X {
		f.Max.X = p.X
	}
	if f.Max.Y < p.Y {
		f.Max.Y = p.Y
	}
}
func (f *Field) FillSparseAir(p Point, r int) {
	f.UpdateSize(Point{p.X + r, p.Y + r})
	f.UpdateSize(Point{p.X - r, p.Y - r})
	for i := 0; i <= r; i++ {
		rng := NewRange(p.X+i-r, p.X+r-i)
		for _, y := range []int{p.Y + i, p.Y - i} {
			if f.Gaps[y] == nil {
				f.Gaps[y] = []*Range{}
			}
			nrs := []*Range{rng.Copy()}
			nrs = append(nrs, f.Gaps[y]...)
			f.Gaps[y] = MergeRanges(nrs)
		}
	}
}
func (f *Field) AddSensor(s, b Point) {
	// We can mark anything with this distance or less as Air
	dist := s.Dist(b)
	log.Printf("New sensor+beacon: S:%s B:%s dist=%d", s, b, dist)
	f.UpdateSize(s)
	f.UpdateSize(b)
	f.FillSparseAir(s, dist)

	f.Objs[s] = Sensor
	f.Objs[b] = Beacon
}

func MergeRanges(gaps []*Range) []*Range {
	keep := map[*Range]bool{}
	for _, g := range gaps {
		found := false
		rs := []*Range{}
		for r, v := range keep {
			if !v {
				continue
			}
			rs = append(rs, r)
		}
		for _, r := range rs {
			keep[r] = false
			if r.Merge(g) {
				// log.Printf("Merged %s -> %s", g, r)
				keep[r] = true
				keep[g] = false
				found = true
				break
			}
			// log.Printf("Keeping %s", r)
			keep[r] = true
		}
		if !found {
			// log.Printf("Adding %s", g)
			keep[g] = true
		}
	}
	res := []*Range{}
	for r, v := range keep {
		if v {
			// log.Printf("Returning %s", r)
			res = append(res, r.Copy())
		}
	}
	sort.Slice(res, func(a, b int) bool {
		return res[a].From < res[b].From || res[a].From == res[b].From && res[a].To < res[b].To
	})
	return res
}

func one(f *Field, y int) int {
	res := 0
	rs := MergeRanges(f.Gaps[y])
	for _, r := range rs {
		res += r.To - r.From
	}
	return res
}

func two(f *Field, bound int) int64 {
	// There has to be exactly one spot where a beacon can be
	// This means there are either two gaps on that line
	// #####.###
	// or a single gap
	// .########

	var found *Point
	log.Printf("Scanning 0/%d gaps...", len(f.Gaps))
	var last time.Time
	stats := map[string]int{}
	for y := 0; found == nil && y <= bound; y++ {
		if time.Now().Sub(last) > time.Second {
			log.Printf("Scanning %d/%d gaps: %v", y, len(f.Gaps), stats)
			last = time.Now()
		}
		if len(f.Gaps[y]) > 2 {
			stats["many"]++
			continue
		}
		gs := f.Gaps[y]
		if len(gs) == 1 {
			stats["single"]++
			if gs[0].From == 1 && gs[0].To >= bound {
				found = &Point{0, y}
			} else if gs[0].From <= 0 && gs[0].To == bound-1 {
				found = &Point{bound, y}
			}
			continue
		} else {
			log.Printf("double %d: %v", y, gs)
			stats["two"]++
			if gs[0].To+2 == gs[1].From {
				found = &Point{gs[0].To + 1, y}
			}
		}
	}
	if found != nil {
		log.Printf("Found: %s", found)
		return int64(found.X*4000000) + int64(found.Y)
	}
	log.Printf("none found: %v", stats)

	return int64(-1)
}

func readFile(path string) *Field {
	res := common.AsStrings(common.ReadTransformedFile(path, common.IgnoreBlankLines))

	f := NewField()
	for _, l := range res {
		ints := common.ExtractInts(l)
		f.AddSensor(Point{ints[0], ints[1]}, Point{ints[2], ints[3]})
	}

	return f
}

func main() {
	log.Println("Reading data...")
	data := readFile(os.Args[1])

	log.Println("Part A")
	res := one(data, 2000000)
	fmt.Printf("%v\n", res)

	log.Println("Part B")
	// data = readFile(os.Args[1])
	res2 := two(data, 4000000)
	fmt.Printf("%v\n", res2)
}
