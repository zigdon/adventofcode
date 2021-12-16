package main

import (
	"log"
	"time"

	"github.com/zigdon/adventofcode/common"
)

type coord struct {
	X, Y int
}

type board struct {
	Danger   [][]int
	Exit     coord
	Best     int
	Attempt  *trail
	Solution *trail
}

type step struct {
	Pos    coord
	Danger int
	Seen   map[coord]bool
}

type dir int

const (
	RIGHT dir = iota
	DOWN
	LEFT
	UP
)

func (d dir) String() string {
	return []string{"->", "vv", "<-", "^^"}[d]
}

func move(cur *step, d dir) *step {
	next := &step{
		Pos:  cur.Pos,
		Seen: make(map[coord]bool),
	}
	switch d {
	case RIGHT:
		next.Pos.X++
	case LEFT:
		next.Pos.X--
	case UP:
		next.Pos.Y--
	case DOWN:
		next.Pos.Y++
	}

	for k, v := range cur.Seen {
		next.Seen[k] = v
	}
	next.Seen[next.Pos] = true

	return next
}

type trail struct {
	Steps []*step
}

func (b board) valid(s *step) bool {
	if s.Pos.X < 0 {
		return false
	}
	if s.Pos.Y < 0 {
		return false
	}
	if s.Pos.X > b.Exit.X {
		return false
	}
	if s.Pos.Y > b.Exit.Y {
		return false
	}

	return true
}

func (b board) guess() *step {
	// Look at the last step in the trail, pick the next step we haven't tried before
	// Mark it as tried in the current step, create a new step, add it to the trail, and return it.
	// If nothing left to try, return nil
	cur := b.Attempt.Steps[len(b.Attempt.Steps)-1]
	for _, d := range []dir{RIGHT, DOWN, LEFT, UP} {
		if n := move(cur, d); !cur.Seen[n.Pos] && b.valid(n) {
			n.Danger = cur.Danger + b.Danger[n.Pos.Y][n.Pos.X]
			return n
		}
	}
	return nil
}

func (b board) backtrack() bool {
	attempts := len(b.Attempt.Steps)
	if attempts > 1 {
		b.Attempt.Steps = b.Attempt.Steps[:attempts-1]
		return true
	}
	return false
}

func (b board) solve(x, y int) (*trail, int) {
	// - at each position, try RDLU (avoiding locations you've been to), add the danger
	// - if the danger is higher than best-case, pop
	lastPrint := time.Now()
	count := int64(0)
	for {
		count++
		cur := b.Attempt.Steps[len(b.Attempt.Steps)-1]
		next := b.guess()
		if next == nil {
			if time.Now().Sub(lastPrint) > 10*time.Second {
				draw(b.Attempt, b.Exit)
				lastPrint = time.Now()
				log.Printf("guesses: %d", count)
			}
			if !b.backtrack() {
				break
			}
			continue
		}
		cur.Seen[next.Pos] = true
		if b.Best != 0 && next.Danger > b.Best {
			continue
		}
		b.Attempt.Steps = append(b.Attempt.Steps, next)
		if next.Pos == b.Exit {
			b.Best = next.Danger
			sol := &trail{}
			for _, s := range b.Attempt.Steps {
				sol.Steps = append(sol.Steps, &step{Pos: s.Pos, Danger: s.Danger})
			}
			b.Solution = sol
			log.Printf("Got to the exit (%d steps, %d danger)!", len(sol.Steps), b.Best)
			draw(sol, b.Exit)
			if !b.backtrack() {
				break
			}
		}
	}

	draw(b.Solution, b.Exit)
	return b.Solution, b.Best
}

func draw(t *trail, exit coord) {
	m := make(map[coord]bool)
	for _, s := range t.Steps {
		m[s.Pos] = true
	}
	for y := 0; y <= exit.Y; y++ {
		line := ""
		for x := 0; x <= exit.X; x++ {
			if m[coord{x, y}] {
				line += "#"
			} else {
				line += "."
			}
		}
		log.Print(line)
	}
}

func readFile(path string) board {
	data := common.AsIntGrid(common.ReadTransformedFile(
		path,
		common.IgnoreBlankLines,
		common.Split(""),
	))

	start := &step{
		Pos:  coord{0, 0},
		Seen: map[coord]bool{{0, 0}: true},
	}
	return board{
		Danger: data,
		Exit:   coord{len(data[0]) - 1, len(data) - 1},
		Attempt: &trail{
			Steps: []*step{start},
		},
	}
}

func main() {
	b := readFile("input.txt")
	p, sc := b.solve(0, 0)
	draw(p, b.Exit)
	log.Printf("Danger = %d", sc)
}
