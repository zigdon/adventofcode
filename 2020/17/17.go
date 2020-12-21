package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type coord struct {
	X, Y, Z, W int
}

type space struct {
	Data           map[coord]bool
	Lx, Ly, Lz, Lw int
	Hx, Hy, Hz, Hw int
	Initialized    bool
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
		s.Hw = p.W
		s.Lw = p.W
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
		if p.W > s.Hw {
			s.Hw = p.W
		}
		if p.W < s.Lw {
			s.Lw = p.W
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

func (s *space) loadSlice(slice string, z, w int) {
	for y, line := range strings.Split(slice, "\n") {
		for x, c := range line {
			if c == '#' {
				s.set(coord{x, y, z, w}, true)
			}
		}
	}
}

func (s *space) count(p coord) int {
	res := 0
	for _, x := range []int{-1, 0, 1} {
		for _, y := range []int{-1, 0, 1} {
			for _, z := range []int{-1, 0, 1} {
				for _, w := range []int{-1, 0, 1} {
					if x == 0 && y == 0 && z == 0 && w == 0 {
						continue
					}
					if s.get(coord{p.X + x, p.Y + y, p.Z + z, p.W + w}) {
						res++
					}
				}
			}
		}
	}
	return res
}

func (s *space) dump() string {
	var res string
	for w := s.Lw; w <= s.Hw; w++ {
		for z := s.Lz; z <= s.Hz; z++ {
			res += fmt.Sprintf("Layer: %d, %d\n", z, w)
			for y := s.Ly; y <= s.Hy; y++ {
				for x := s.Lx; x <= s.Hx; x++ {
					if s.get(coord{x, y, z, w}) {
						res += "#"
					} else {
						res += "."
					}
				}
				res += "\n"
			}
			res += "\n\n"
		}
	}

	return res
}

func (s *space) evolve() {
	next := newSpace()
	for x := s.Lx - 1; x <= s.Hx+1; x++ {
		for y := s.Ly - 1; y <= s.Hy+1; y++ {
			for z := s.Lz - 1; z <= s.Hz+1; z++ {
				p := coord{x, y, z, 0}
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
	s.Lx, s.Ly, s.Lz, s.Lw = next.Lx, next.Ly, next.Lz, next.Lw
	s.Hx, s.Hy, s.Hz, s.Lw = next.Hx, next.Hy, next.Hz, next.Hw
}

func (s *space) evolve4d() {
	next := newSpace()
	for x := s.Lx - 1; x <= s.Hx+1; x++ {
		for y := s.Ly - 1; y <= s.Hy+1; y++ {
			for z := s.Lz - 1; z <= s.Hz+1; z++ {
				for w := s.Lw - 1; w <= s.Hw+1; w++ {
					p := coord{x, y, z, w}
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
	}
	s.Data = next.Data
	s.Lx, s.Ly, s.Lz, s.Lw = next.Lx, next.Ly, next.Lz, next.Lw
	s.Hx, s.Hy, s.Hz, s.Hw = next.Hx, next.Hy, next.Hz, next.Hw
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("can't read %q: %v", input, err)
	}
	s := newSpace()
	s.loadSlice(string(data), 0, 0)
	for i := 1; i <= 6; i++ {
		s.evolve()
	}
	fmt.Println(len(s.getTrue()))

	s = newSpace()
	s.loadSlice(string(data), 0, 0)
	for i := 1; i <= 6; i++ {
		s.evolve4d()
	}
	fmt.Println(len(s.getTrue()))
}
