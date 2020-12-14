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
	wayX, wayY    int
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

func (s *ship) moveTowardsWaypoint(dist int) {
	s.x = s.x + dist*s.wayX
	s.y = s.y + dist*s.wayY
}

func (s *ship) rotate(deg int) {
	for deg > 0 {
		s.wayX, s.wayY = s.wayY, -s.wayX
		deg = deg - 90
	}
}

func (s *ship) moveToWaypoint(dirs []string) (int, int) {
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
			s.wayY = s.wayY + n
		case "S":
			s.wayY = s.wayY - n
		case "E":
			s.wayX = s.wayX + n
		case "W":
			s.wayX = s.wayX - n
		case "F":
			s.moveTowardsWaypoint(n)
		case "L":
			s.rotate(360 - n)
		case "R":
			s.rotate(n)
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
	fmt.Printf("moved %d,%d -> %f\n", x, y, math.Abs(float64(x))+math.Abs(float64(y)))

	s = &ship{heading: 90, wayX: 10, wayY: 1}
	x, y = s.moveToWaypoint(strings.Split(string(data), "\n"))
	fmt.Printf("moved %d,%d -> %f\n", x, y, math.Abs(float64(x))+math.Abs(float64(y)))
}
