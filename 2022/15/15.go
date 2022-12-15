package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"

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

func (r *Range) String() string {
	return fmt.Sprintf("R[%d:%d]", r.From, r.To)
}
func (r *Range) Contains(x int) bool {
	return x >= r.From && x <= r.To
}
func (r *Range) Merge(r2 *Range) bool {
	if r.To+1 < r2.From || r.From-1 > r2.To {
		return false
	}
	r.From = int(math.Min(float64(r.From), float64(r2.From)))
	r.To = int(math.Max(float64(r.To), float64(r2.To)))
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
	res := []string{
		"",
		strings.Repeat(" ", -f.Min.X+5) + "0",
	}
	for y := f.Min.Y; y <= f.Max.Y; y++ {
		l := fmt.Sprintf("%3d: ", y)
		rs := MergeRanges(f.Gaps[y])
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
		rng := &Range{p.X + i - r, p.X + r - i}
		for _, y := range []int{p.Y + i, p.Y - i} {
			if f.Gaps[y] == nil {
				f.Gaps[y] = []*Range{rng.Copy()}
			} else {
				f.Gaps[y] = MergeRanges(append(f.Gaps[y], rng))
			}
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
	rs := []*Range{}
	for _, g := range gaps {
		found := false
		for _, r := range rs {
			if r.Merge(g) {
				found = true
				break
			}
		}
		if !found {
			rs = append(rs, g.Copy())
		}
	}
	return rs
}

func one(f *Field, y int) int {
	res := 0
	rs := MergeRanges(f.Gaps[y])
	for _, r := range rs {
		res += r.To - r.From
	}
	return res
}

func two(data *Field) int {
	return 0
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
	data = readFile(os.Args[1])
	res = two(data)
	fmt.Printf("%v\n", res)
}
