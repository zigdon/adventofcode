package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type CPUState struct {
	X int
}

func (s CPUState) String() string {
	return fmt.Sprintf("state[x=%d]", s.X)
}

type CPU struct {
	X       int
	Cycle   int
	History []CPUState
	Screen  [][]rune
}

func NewCPU() *CPU {
	c := &CPU{
		X:       1,
		History: []CPUState{},
	}

	c.Screen = [][]rune{
		[]rune(strings.Repeat(".", 40)),
		[]rune(strings.Repeat(".", 40)),
		[]rune(strings.Repeat(".", 40)),
		[]rune(strings.Repeat(".", 40)),
		[]rune(strings.Repeat(".", 40)),
		[]rune(strings.Repeat(".", 40)),
	}

	c.Show()
	return c
}

func (cpu *CPU) Tick() {
	p := cpu.Cycle%40
	l := cpu.Cycle / 40
	if -1 <= p-cpu.X && p-cpu.X <= 1 {
		// log.Printf("At cycle #%d, setting (%d,%d)", cpu.Cycle, p, l)
		cpu.Screen[l][p] = '#'
	}
}

func (cpu *CPU) Run(op Op) {
    // log.Printf(" x: %d  => %s", cpu.X, op)
	for i := 0; i < op.Cycles; i++ {
		cpu.Tick()
		cpu.History = append(cpu.History, CPUState{cpu.X})
        cpu.Cycle++
	}
	if op.Do != nil {
		op.Do(cpu, op.Arg)
	}
}

func (cpu *CPU) Read() int {
	res := 0
	for n := 20; n <= len(cpu.History); n += 40 {
		// log.Printf("%d: %s", n, cpu.History[n-1])
		res += n * cpu.History[n-1].X
	}
	return res
}

func (cpu *CPU) Show() {
	for _, l := range cpu.Screen {
		fmt.Println(string(l))
	}
}

type Op struct {
	Name   string
	Cycles int
	Do     func(*CPU, int)
	Arg    int
}

func (o Op) String() string {
	return fmt.Sprintf("%s[%d]", o.Name, o.Arg)
}

var defs = map[string]Op{
	"noop": {Name: "noop", Cycles: 1},
	"addx": {Name: "addx", Cycles: 2, Do: func(cpu *CPU, arg int) {
		cpu.X += arg
	}},
}

func one(data []Op) int {
	cpu := NewCPU()
	for _, o := range data {
		cpu.Run(o)
	}

	return cpu.Read()
}

func two(data []Op) []string {
	cpu := NewCPU()
	for _, o := range data {
		cpu.Run(o)
		cpu.Show()
		fmt.Println()
	}

    res := []string{}
    for _, l := range cpu.Screen{
      res = append(res, string(l))
    }
	return res
}

func readFile(path string) ([]Op, error) {
	lines := common.AsStrings(common.ReadTransformedFile(path, common.IgnoreBlankLines))
	ops := []Op{}
	for n, l := range lines {
		words := strings.Split(l, " ")
		if op, ok := defs[words[0]]; !ok {
			return nil, fmt.Errorf("Bad op in line %d: %q", n, op)
		} else {
			if len(words) > 1 {
				op.Arg = common.MustInt(words[1])
			}
			ops = append(ops, op)
		}
	}

	return ops, nil
}

func main() {
	data, err := readFile(os.Args[1])
	if err != nil {
		log.Fatalf("%v", err)
	}

	res := one(data)
	fmt.Printf("%v\n", res)

	two(data)
}
