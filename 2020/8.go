package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type cmd struct {
	Line int
	Inst string
	Arg  int
}

func compile(data []string) []*cmd {
	idx := 0

	var code []*cmd
	for i, line := range data {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		bits := strings.Split(line, " ")
		arg, err := strconv.Atoi(bits[1])
		if err != nil {
			log.Fatalf("bad arg in line %d %q: %v", i, line, err)
		}
		code = append(code, &cmd{idx, bits[0], arg})
		idx++
	}

	return code
}

func run(code []*cmd) (int, int) {
	acc := 0
	loop := -1
	ptr := 0
	seen := make(map[int]bool)
	for {
		if ptr >= len(code) {
			break
		}
		if seen[ptr] {
			return acc, ptr
		}
		seen[ptr] = true
		c := code[ptr]
		switch c.Inst {
		case "jmp":
			ptr = ptr + c.Arg
			continue
		case "acc":
			acc = acc + c.Arg
		}
		ptr++
	}

	return acc, loop
}

func hack(code []*cmd) (int, int) {
	for i, inst := range code {
		if inst.Inst == "acc" {
			continue
		}
		try, orig := "jmp", "nop"
		if inst.Inst == "jmp" {
			try, orig = orig, try
		}
		log.Printf("Trying fix %s->%s at line %d\n", orig, try, i)
		code[i].Inst = try
		acc, loop := run(code)
		if loop > -1 {
			code[i].Inst = orig
			continue
		}
		return acc, i
	}

	return -1, -1
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("can't read input: %v", err)
	}

	code := compile(strings.Split(string(data), "\n"))
	acc, loop := run(code)
	fmt.Printf("ACC=%d, LOOP=%d EOL\n", acc, loop)
	if loop > -1 {
		acc, fix := hack(code)
		fmt.Printf("ACC=%d, FIX=%d EOL\n", acc, fix)
	}
}
