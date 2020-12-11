package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type space int

const (
	Floor space = iota
	Empty
	Occupied
)

func (s space) String() string {
	return [...]string{"Floor", "Empty", "Occupied"}[s]
}

func (s space) Char() string {
	return [...]string{".", "L", "#"}[s]
}

type board struct {
	Spaces     [][]space
	Width      int
	Height     int
	LastChange int
}

func (b *board) InitFromString(rows []string) *board {
	if b == nil {
		b = &board{}
	}
	b.Spaces = [][]space{}
	b.Width = 0
	b.Height = 0
	for y, row := range rows {
		if len(row) == 0 {
			continue
		}
		newRow := []space{}
		for x, seat := range row {
			switch seat {
			case '.':
				newRow = append(newRow, Floor)
			case 'L':
				newRow = append(newRow, Empty)
			case '#':
				newRow = append(newRow, Occupied)
			default:
				log.Fatalf("bad seat at %d,%d: %q", x, y, seat)
			}
		}
		b.Spaces = append(b.Spaces, newRow)
		if b.Width == 0 {
			b.Width = len(newRow)
		} else if b.Width != len(newRow) {
			log.Fatalf("bad row length: expected %d, got %d", b.Width, len(newRow))
		}
	}
	b.Height = len(b.Spaces)

	return b
}

func (b *board) CountNear(x, y int) int {
	count := 0
	for _, dy := range []int{-1, 0, 1} {
		destY := y + dy
		if destY < 0 || destY >= b.Height {
			continue
		}
		for _, dx := range []int{-1, 0, 1} {
			destX := x + dx
			if destX < 0 || destX >= b.Width || (dx == 0 && dy == 0) {
				continue
			}
			if b.Spaces[destY][destX] == Occupied {
				count++
			}
		}
	}

	return count
}

func (b *board) Look(x, y, dx, dy int) space {
	var curX, curY = x, y
	for {
		curX = curX + dx
		curY = curY + dy
		if curX < 0 || curY < 0 || curX >= b.Width || curY >= b.Height {
			return Floor
		}
		if b.Spaces[curY][curX] != Floor {
			return b.Spaces[curY][curX]
		}
	}
}

func (b *board) CountSeen(x, y int) int {
	count := 0
	for _, dy := range []int{-1, 0, 1} {
		for _, dx := range []int{-1, 0, 1} {
			if dx == 0 && dy == 0 {
				continue
			}
			if b.Look(x, y, dx, dy) == Occupied {
				count++
			}
		}
	}

	return count
}

func (b *board) Evolve(tol int, los bool) *board {
	next := [][]space{}
	b.LastChange = 0
	for y, row := range b.Spaces {
		next = append(next, []space{})
		for x := range row {
			if b.Spaces[y][x] == Floor {
				next[y] = append(next[y], Floor)
				continue
			}
			var consider int
			if los {
				consider = b.CountSeen(x, y)
			} else {
				consider = b.CountNear(x, y)
			}
			var n space
			if consider == 0 {
				n = Occupied
			} else if consider >= tol {
				n = Empty
			} else {
				n = b.Spaces[y][x]
			}
			if n != b.Spaces[y][x] {
				b.LastChange++
			}
			next[y] = append(next[y], n)
		}
	}
	b.Spaces = next
	return b
}

func (b *board) CountOccupied() int {
	res := 0
	for _, row := range b.Spaces {
		for _, c := range row {
			if c == Occupied {
				res++
			}
		}
	}

	return res
}

func (b *board) String() string {
	out := ""
	for _, row := range b.Spaces {
		for _, c := range row {
			out = out + c.Char()
		}
		out = out + "\n"
	}
	return out
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}

	var b *board
	b1 := b.InitFromString(strings.Split(string(data), "\n"))
	gen := 0
	for {
		gen++
		b1.Evolve(4, false)
		log.Printf("Gen #%d: %d", gen, b1.LastChange)
		// log.Print(b1)
		if b1.LastChange == 0 {
			break
		}
	}
	fmt.Printf("Occupied: %d\n", b1.CountOccupied())

	b2 := b.InitFromString(strings.Split(string(data), "\n"))
	gen = 0
	for {
		gen++
		b2.Evolve(5, true)
		log.Printf("Gen #%d: %d", gen, b2.LastChange)
		// log.Print(b1)
		if b2.LastChange == 0 {
			break
		}
	}
	fmt.Printf("Occupied: %d\n", b2.CountOccupied())

}
