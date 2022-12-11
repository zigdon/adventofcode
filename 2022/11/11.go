package main

import (
	"fmt"
	"log"
	"os"
    "sort"
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

func (m *Monkey) Has(item int) bool {
	for _, i := range m.Items {
		if i == item {
			return true
		}
	}

	return false
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
	Monkeys     []*Monkey
	Inspections []int
    Decay int
}

func NewTroop(m []*Monkey, d int) *Troop {
	return &Troop{Monkeys: m, Decay: d, Inspections: make([]int, len(m))}
}

func (t *Troop) Turn(i int) int {
	m := t.Monkeys[i]
	cnt := 0

	for _, item := range m.Items {
		cnt++
		// log.Printf("M#%d examining %d", i, item)
		next := m.Op(item)
		// log.Printf(" -> %d / decay -> %d", next, next / 3)
		next /= t.Decay
		if next%m.Test == 0 {
			// log.Printf("%d divisible by %d -> %d", next, m.Test, m.True)
			t.Monkeys[m.True].Items = append(t.Monkeys[m.True].Items, next)
		} else {
			// log.Printf("%d not divisible by %d -> %d", next, m.Test, m.False)
			t.Monkeys[m.False].Items = append(t.Monkeys[m.False].Items, next)
		}
	}
	m.Items = []int{}

	return cnt
}

func (t *Troop) Round() {
	for n := range t.Monkeys {
		t.Inspections[n] += t.Turn(n)
	}
}

func one(data []*Monkey) int {
    tr := NewTroop(data,3)
    for n := 0; n < 20; n++ {
      tr.Round()
    }
    ints := []int{}
    for n := range tr.Monkeys {
      ints = append(ints, tr.Inspections[n])
    }
    log.Printf("After 20 rounds: %v", ints)
    sort.Sort(sort.Reverse(sort.IntSlice(ints)))
	return ints[0] * ints[1]
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
