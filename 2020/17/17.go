package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type coord struct {
	X, Y, Z int
}

type space struct {
	Data        map[coord]bool
	Lx, Ly, Lz  int
	Hx, Hy, Hz  int
	Initialized bool
}

func newSpace() *space {
	return &space{Data: make(map[coord]bool)}
}

func (s *space) set(p coord, v bool) {
	if !s.Initialized {
		s.Hx = p.X
		s.Lx = p.X
		s.Hy = p.Y
		s.Ly = p.Y
		s.Hz = p.Z
		s.Lz = p.Z
		s.Initialized = true
	} else if v {
		if p.X > s.Hx {
			s.Hx = p.X
		}
		if p.X < s.Lx {
			s.Lx = p.X
		}
		if p.Y > s.Hy {
			s.Hy = p.Y
		}
		if p.Y < s.Ly {
			s.Ly = p.Y
		}
		if p.Z > s.Hz {
			s.Hz = p.Z
		}
		if p.Z < s.Lz {
			s.Lz = p.Z
		}
	}

	s.Data[p] = v
}

func (s *space) get(p coord) bool {
	return s.Data[p]
}

func (s *space) getTrue() []coord {
	var res []coord
	for k, v := range s.Data {
		if v {
			res = append(res, k)
		}
	}

	return res
}

func (s *space) loadSlice(slice string, z int) {
	for y, line := range strings.Split(slice, "\n") {
		for x, c := range line {
			if c == '#' {
				s.set(coord{x, y, z}, true)
			}
		}
	}
}

func (s *space) count(p coord) int {
	res := 0
	for _, x := range []int{-1, 0, 1} {
		for _, y := range []int{-1, 0, 1} {
			for _, z := range []int{-1, 0, 1} {
				if x == 0 && y == 0 && z == 0 {
					continue
				}
				if s.get(coord{p.X + x, p.Y + y, p.Z + z}) {
					res++
				}
			}
		}
	}
	return res
}

func (s *space) dump() string {
	var res string
	for z := s.Lz; z <= s.Hz; z++ {
		res += fmt.Sprintf("Layer: %d\n", z)
		for y := s.Ly; y <= s.Hy; y++ {
			for x := s.Lx; x <= s.Hx; x++ {
				if s.get(coord{x, y, z}) {
					res += "#"
				} else {
					res += "."
				}
			}
			res += "\n"
		}
	}

	return res
}

func (s *space) evolve() {
	next := newSpace()
	for x := s.Lx - 1; x <= s.Hx+1; x++ {
		for y := s.Ly - 1; y <= s.Hy+1; y++ {
			for z := s.Lz - 1; z <= s.Hz+1; z++ {
				p := coord{x, y, z}
				cnt := s.count(p)
				if s.get(p) {
					if cnt == 2 || cnt == 3 {
						next.set(p, true)
					}
				} else {
					if cnt == 3 {
						next.set(p, true)
					}
				}
			}
		}
	}
	s.Data = next.Data
	s.Lx, s.Ly, s.Lz = next.Lx, next.Ly, next.Lz
	s.Hx, s.Hy, s.Hz = next.Hx, next.Hy, next.Hz
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("can't read %q: %v", input, err)
	}
	s := newSpace()
	s.loadSlice(string(data), 0)
	for i := 1; i <= 6; i++ {
		s.evolve()
	}
	fmt.Println(len(s.getTrue()))
}
