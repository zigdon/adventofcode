package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type Monkey struct {
	ID    int
	Items []int
	Op    func(int) int
	Test  int
	True  int
	False int
}

var ops = map[string]func(a, b int) int{
	"+": func(a, b int) int { return a + b },
	"-": func(a, b int) int { return a - b },
	"*": func(a, b int) int { return a * b },
	"/": func(a, b int) int { return a / b },
}

func ExtractInts(in string) []int {
	res := []int{}
	for _, w := range strings.Split(in, " ") {
		n, err := strconv.Atoi(strings.TrimRight(w, ":,"))
		if err != nil {
			continue
		}
		res = append(res, n)
	}

	return res
}

func NewMonkey(lines []string) *Monkey {
	m := &Monkey{Items: []int{}}
	m.ID = ExtractInts(strings.TrimSuffix(lines[0], ":"))[0]
	m.Items = ExtractInts(lines[1])
	m.Test = ExtractInts(lines[3])[0]
	m.True = ExtractInts(lines[4])[0]
	m.False = ExtractInts(lines[5])[0]

	formula := strings.Split(lines[2], " = ")[1]
	parts := strings.Split(formula, " ")
	if parts[0] != "old" {
		log.Printf("Invalid formula: %q", formula)
		os.Exit(1)
	}
	if parts[2] == "old" {
		m.Op = func(old int) int { return ops[parts[1]](old, old) }
	} else {
		m.Op = func(old int) int { return ops[parts[1]](old, common.MustInt(parts[2])) }
	}

	return m
}

type Troop struct {
	Monkeys []Monkey
}

func one(data []*Monkey) int {
	return 0
}

func two(data []*Monkey) int {
	return 0
}

func readFile(path string) ([]*Monkey, error) {
	res := common.ReadTransformedFile(path, common.Block)
	m := []*Monkey{}
	for _, ent := range res {
		m = append(m, NewMonkey(common.AsStrings(ent)))
	}

	return m, nil
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
