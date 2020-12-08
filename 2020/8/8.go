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

func run(inst []string) (int, int) {
	acc := 0
	loop := -1
	idx := 0

	var code []*cmd
	for i, line := range inst {
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

	ptr := 0
	seen := make(map[int]bool)
	for {
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
		if ptr >= len(code) {
			break
		}
	}

	return acc, loop
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("can't read input: %v", err)
	}

	acc, loop := run(strings.Split(string(data), "\n"))
	fmt.Printf("ACC=%d, LOOP=%d EOL\n", acc, loop)
}
