package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zigdon/adventofcode/common"
)

type Point struct {
  X, Y int
}
func(p Point)String() string {
  return fmt.Sprintf("[%d,%d]", p.X, p.Y)
}

type Map struct {
  Alt map[Point]int
  Start Point
  End Point
  Size Point
}

func NewMap() *Map {
  return &Map{
    Alt: map[Point]int{},
    Start: Point{},
    End: Point{},
  }
}

func one(m *Map) int {
	return 0
}

func two(m *Map) int {
	return 0
}

func readFile(path string) *Map {
	res := common.AsStrings(common.ReadTransformedFile(path, common.IgnoreBlankLines))
    m := NewMap()
    m.Size = Point{len(res[0]), len(res)}
    for y, l := range res {
      for x, c := range l {
        p := Point{x, y}
        if c == 'S' {
          c = 'a'
          m.Start = p
        } else if c == 'E' {
          c = 'z'
          m.End = p
        }
        m.Alt[p] = int(c - 'a')
      }
    }

	return m
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
