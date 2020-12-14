package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type ship struct {
	x, y, heading int
}

func (s *ship) move(dist int) {
	for s.heading < 0 {
		s.heading = s.heading + 360
	}
	s.heading = s.heading % 360
	switch s.heading {
	case 0:
		s.y = s.y + dist
	case 90:
		s.x = s.x + dist
	case 180:
		s.y = s.y - dist
	case 270:
		s.x = s.x - dist
	}
}

func (s *ship) track(dirs []string) (int, int) {
	for _, step := range dirs {
		step := strings.TrimSpace(step)
		if len(step) == 0 {
			continue
		}
		op, arg := step[0:1], step[1:]
		n, err := strconv.Atoi(arg)
		if err != nil {
			log.Fatalf("bad arg %q: %v", step, err)
		}
		switch op {
		case "N":
			s.y = s.y + n
		case "S":
			s.y = s.y - n
		case "E":
			s.x = s.x + n
		case "W":
			s.x = s.x - n
		case "F":
			s.move(n)
		case "L":
			s.heading = s.heading - n
		case "R":
			s.heading = s.heading + n
		}
	}

	return s.x, s.y
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}
	s := &ship{heading: 90}
	x, y := s.track(strings.Split(string(data), "\n"))
	fmt.Printf("moved %d,%d -> %f", x, y, math.Abs(float64(x))+math.Abs(float64(y)))
}
