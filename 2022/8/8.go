package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type grid struct {
	Trees [][]int
	W, H  int
}

func (g grid) String() string {
	res := []string{}
	for y := range g.Trees {
		line := fmt.Sprintf("%3d: ", y)
		for x := range g.Trees[y] {
			line += fmt.Sprintf(" %d", g.Trees[y][x])
		}
		res = append(res, line)
	}

	return strings.Join(res, "\n")
}
func (g grid) Do(f func(x, y int)) {
	for y := range g.Trees {
		for x := range g.Trees[y] {
			f(x, y)
		}
	}
}
func (g grid) ScenicScore(x, y int) int {
	if x == 0 || y == 0 || x == g.W-1 || y == g.H-1 {
		return 0
	}
	res := 1
	h := g.Trees[y][x]
	seen := 0
	for i := 0; i < g.W; i++ {
		if i == x {
			res *= seen
			// log.Printf("ToLeft: %d -> %d", seen, res)
			seen = 0
			continue
		}
		seen++
		if g.Trees[y][i] >= h {
			// log.Printf("blocked at %d,%d", i, y)
			if i > x {
				break
			}
			seen = 1
		}
	}
	res *= seen
	// log.Printf("ToRight: %d -> %d", seen, res)
	seen = 0
	for i := 0; i < g.H; i++ {
		if i == y {
			res *= seen
			// log.Printf("ToUp: %d -> %d", seen, res)
			seen = 0
			continue
		}
		seen++
		if g.Trees[i][x] >= h {
			// log.Printf("blocked at %d,%d", x, i)
			if i > y {
				break
			}
			seen = 1
		}
	}
	res *= seen
	// log.Printf("ToDown: %d -> %d", seen, res)

	return res
}
func (g grid) IsVisible(x, y int) bool {
	if x == 0 || y == 0 || x == g.W-1 || y == g.H-1 {
		return true
	}
	visLeft := true
	visRight := true
	visUp := true
	visDown := true
	h := g.Trees[y][x]
	for i := 0; i < g.W; i++ {
		t := g.Trees[y][i]
		if i == x {
			continue
		}
		if i < x {
			if !visLeft {
				continue
			}
			if t >= h {
				visLeft = false
			}
		} else {
			if !visRight {
				continue
			}
			if t >= h {
				visRight = false
			}
		}
	}
	for i := 0; i < g.H; i++ {
		t := g.Trees[i][x]
		if i == y {
			continue
		}
		if i < y {
			if !visUp {
				continue
			}
			if t >= h {
				visUp = false
			}
		} else {
			if !visDown {
				continue
			}
			if t >= h {
				visDown = false
			}
		}
	}
	if visLeft || visRight || visUp || visDown {
		// log.Printf("(%d,%d) visibility:", x, y)
		// log.Printf("  %v", visUp)
		// log.Printf("%v %v", visLeft, visRight)
		// log.Printf("  %v", visDown)
		return true
	}

	return false
}

func makeGrid(data [][]int) grid {
	return grid{Trees: data, H: len(data), W: len(data[0])}
}

func one(data grid) int {
	count := 0
	data.Do(func(x, y int) {
		if data.IsVisible(x, y) {
			count++
		}
	})
	return count
}

func two(data grid) int {
    best := 0
    data.Do(func(x, y int) {
      sc := data.ScenicScore(x, y)
      if sc > best {
        best = sc
        log.Printf("New best view at (%d,%d): %d", x, y, best)
      }
    })
	return best
}

func readFile(path string) (grid, error) {
	res := common.ReadTransformedFile(path, common.IgnoreBlankLines, common.Split(""))

	return makeGrid(common.AsIntGrid(res)), nil
}

func main() {
	data, err := readFile(os.Args[1])
	if err != nil {
		log.Fatalf("%v", err)
	}

	res := one(data)
	fmt.Printf("%v\n", res)

	res = two(data)
	fmt.Printf("%v\n", res)
}
