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

type Field struct {
	Objs     map[Point]Object
	Min, Max *Point
}

func NewField() *Field {
	return &Field{
		Objs: make(map[Point]Object),
	}
}
func (f *Field) String() string {
	res := []string{
		"",
		strings.Repeat(" ", -f.Min.X+5) + "0",
	}
	for y := f.Min.Y; y <= f.Max.Y; y++ {
		l := fmt.Sprintf("%3d: ", y)
		for x := f.Min.X; x <= f.Max.X; x++ {
			l += f.Objs[Point{x, y}].String()
		}
		res = append(res, l)
	}
	return strings.Join(res, "\n")
}
func (f *Field) Fill() {
	for y := f.Min.Y; y <= f.Max.Y; y++ {
		for x := f.Min.X; x <= f.Max.X; x++ {
			p := Point{x, y}
			if _, ok := f.Objs[p]; !ok {
				f.Objs[p] = Unknown
			}
		}
	}
}
func (f *Field) ParseMap(o Point, m string) {
	for i, l := range strings.Split(m, "\n") {
		y := o.Y + i
		for j, r := range l {
			x := o.X + j
			p := Point{x, y}
			var obj Object
			switch r {
			case '.':
				obj = Unknown
			case '#':
				obj = Air
			case 'S':
				obj = Sensor
			case 'B':
				obj = Beacon
			}
			f.Objs[p] = obj
		}
	}
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
func (f *Field) FillAir(p Point, r int) {
	if r > 1 {
		f.FillAir(p, r-1)
	}
	f.UpdateSize(Point{p.X + r, p.Y + r})
	f.UpdateSize(Point{p.X - r, p.Y - r})
	mkAir := func(p Point) {
		if f.Objs[p] != Unknown {
			return
		}
		f.Objs[p] = Air
	}
	for i := 0; i <= r; i++ {
		mkAir(Point{p.X + i, p.Y + r - i})
		mkAir(Point{p.X - i, p.Y + r - i})
		mkAir(Point{p.X + i, p.Y - r + i})
		mkAir(Point{p.X - i, p.Y - r + i})
	}
}
func (f *Field) AddSensor(s, b Point) {
	// We can mark anything with this distance or less as Air
	dist := s.Dist(b)
	f.UpdateSize(s)
	f.UpdateSize(b)
	f.FillAir(s, dist)

	f.Objs[s] = Sensor
	f.Objs[b] = Beacon
}

func one(data *Field) int {
	return 0
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
	res := one(data)
	fmt.Printf("%v\n", res)

	log.Println("Part B")
	data = readFile(os.Args[1])
	res = two(data)
	fmt.Printf("%v\n", res)
}
