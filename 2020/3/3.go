package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type trees struct {
	data          map[int]map[int]bool
	width, height int
}

func (t *trees) init(data string) {
	t.data = make(map[int]map[int]bool)
	t.width = 0
	t.height = 0
	t.load(data)
}

func (t *trees) get(x, y int) bool {
	if y >= t.height {
		return false
	}
	return t.data[y][x%t.width]
}

func (t *trees) load(data string) {
	lines := strings.Split(string(data), "\n")
	for l, ldata := range lines {
		line := strings.TrimSpace(ldata)
		if len(line) == 0 {
			continue
		}
		if t.data[l] == nil {
			t.data[l] = make(map[int]bool)
			if l >= t.height {
				t.height = l + 1
			}
		}
		for c, r := range line {
			t.data[l][c] = r == '#'
			if c >= t.width {
				t.width = c + 1
			}
		}
	}
}

func (t *trees) readInput(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	t.init(string(data))
	return nil
}

func (t *trees) countTrees(right, down int) int {
	var x, y, res int
	for y < t.height {
		if t.get(x, y) {
			res = res + 1
		}
		x = x + right
		y = y + down
	}
	return res
}

func main() {
	input := os.Args[1]
	terrain := &trees{}
	err := terrain.readInput(input)
	if err != nil {
		log.Fatalf("Can't read input: %v", err)
	}

	count := terrain.countTrees(3, 1)
	fmt.Println(count)

	res := 1
	for _, dir := range []struct{ right, down int }{{1, 1}, {3, 1}, {5, 1}, {7, 1}, {1, 2}} {
		res = res * terrain.countTrees(dir.right, dir.down)
	}
	fmt.Println(res)
}
