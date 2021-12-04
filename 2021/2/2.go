package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/zigdon/adventofcode/common"
)

type direction int

const (
	dir_forward direction = iota + 1
	dir_up
	dir_down
)

type command struct {
	Dir  direction
	Dist int
}

func parseLine(i int, in interface{}) (string, interface{}, error) {
	name := "parseLine"
	l, ok := in.([]string)
	if !ok {
		return name, nil, fmt.Errorf("parseLine expects []string, got %T", l)
	}

	if len(l) != 2 {
		return name, nil, fmt.Errorf("parseLine expects 2 words, got %v in line %d", l, i)
	}

	n, err := strconv.Atoi(l[1])
	if err != nil {
		return name, nil, fmt.Errorf("invalid number %q at line %d", l[1], i)
	}

	var d direction
	switch l[0] {
	case "forward":
		d = dir_forward
	case "up":
		d = dir_up
	case "down":
		d = dir_down
	default:
		return name, nil, fmt.Errorf("invalid direction %q at line %d", l[0], i)
	}

	return name, command{d, n}, nil
}

func readFile(path string) ([]command, error) {
	cmds := []command{}
	for _, l := range common.ReadTransformedFile(path,
		common.IgnoreBlankLines,
		common.SplitWords,
		parseLine,
	) {
		c, ok := l.(command)
		if !ok {
			return nil, fmt.Errorf("bad line %+v (%T)", l, l)
		}
		cmds = append(cmds, c)
	}

	return cmds, nil
}

func pilotSub(inst []command) (int, int, error) {
	var curX, curY int
	for _, c := range inst {
		switch c.Dir {
		case dir_forward:
			curX += c.Dist
		case dir_up:
			curY -= c.Dist
		case dir_down:
			curY += c.Dist
		}
	}

	return curX, curY, nil
}

func pilotSubWithAim(inst []command) (int, int, error) {
	var curX, curY, aim int
	for _, c := range inst {
		switch c.Dir {
		case dir_forward:
			curX += c.Dist
			curY += c.Dist * aim
		case dir_up:
			aim -= c.Dist
		case dir_down:
			aim += c.Dist
		}
	}

	return curX, curY, nil
}

func main() {
	cmds, err := readFile("input.txt")
	if err != nil {
		log.Fatalf("Can't read input: %v\n", err)
	}
	x, y, _ := pilotSub(cmds)

	fmt.Println(x * y)

	x, y, _ = pilotSubWithAim(cmds)

	fmt.Println(x * y)
}
