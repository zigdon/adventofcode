package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type coord struct {
	x, y int
}

type board struct {
	finished bool
	lines    [][]int
	cache    map[int]coord
}

func (b *board) create(data interface{}) *board {
	b.cache = make(map[int]coord)
	d := data.([]interface{})
	for y, l := range d {
		line := []int{}
		for x, c := range l.([]string) {
			n := mustInt(c)
			line = append(line, n)
			b.cache[n] = coord{x, y}
		}
		b.lines = append(b.lines, line)
	}

	return b
}

func (b *board) winLine(y int) bool {
	for _, n := range b.lines[y] {
		if n != 0 {
			return false
		}
	}

	b.finished = true
	return true
}

func (b *board) winCol(x int) bool {
	for _, l := range b.lines {
		if l[x] != 0 {
			return false
		}
	}

	b.finished = true
	return true
}

func (b *board) play(n int) bool {
	if b.finished {
		return false
	}
	c, ok := b.cache[n]
	if !ok {
		return false
	}

	b.lines[c.y][c.x] = 0
	return b.winLine(c.y) || b.winCol(c.x)
}

func (b *board) sum() int {
	s := 0
	for _, l := range b.lines {
		for _, n := range l {
			s += n
		}
	}

	return s
}

func mustInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Can't convert %q to number: %v", s, err)
	}

	return n
}

func readInput(path string) ([]int, []*board) {
	data := common.ReadTransformedFile(
		path,
		common.IgnoreBlankLines,
		common.SplitWords,
	)

	calls := []int{}
	boards := []*board{}
	for _, s := range strings.Split(data[0].([]string)[0], ",") {
		calls = append(calls, mustInt(s))
	}

	l := 1
	for l < len(data) {
		b := &board{}
		boards = append(boards, b.create(data[l:l+5]))
		l += 5
	}

	return calls, boards
}

func main() {
	calls, boards := readInput("input.txt")
	running := len(boards)
	lastWin := -1
	lastPlay := 0
	first := false
	for _, call := range calls {
		lastPlay = call
		for i, board := range boards {
			if board.finished {
				continue
			}
			win := board.play(call)
			if win {
				lastWin = i
				if !first {
					fmt.Printf("First win: board %d, playing %d, sum = %d, result = %d\n",
						i, call, board.sum(), board.sum()*call)
					first = true
				}
				running--
				continue
			}
		}
		if running == 0 {
			break
		}
	}
	fmt.Printf("Last win: board %d, playing %d, sum = %d, result = %d\n",
		lastWin, lastPlay, boards[lastWin].sum(), boards[lastWin].sum()*lastPlay)
}
