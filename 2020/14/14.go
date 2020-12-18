package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type compy struct {
	mem     map[int]int
	mask    string
	maskOn  int
	maskOff int
}

func newCompy() *compy {
	c := &compy{}
	c.mem = make(map[int]int)
	for i := 0; i < 36; i++ {
		c.maskOff = c.maskOff<<1 + 1
	}
	return c
}

func (c *compy) setMask(newMask string) {
	log.Printf("setting mask to %q", newMask)
	c.mask = newMask
	c.maskOn = 0
	c.maskOff = 0
	for i := 0; i < len(newMask); i++ {
		c.maskOn <<= 1
		c.maskOff <<= 1
		if newMask[i] == '1' {
			c.maskOn++
			c.maskOff++
		} else if newMask[i] == 'X' {
			c.maskOff++
		}
	}
	log.Printf("=>  on: %b", c.maskOn)
	log.Printf("=> off: %b", c.maskOff)
}

func (c *compy) sumMem() int {
	sum := 0
	for _, val := range c.mem {
		sum += val
	}
	return sum
}

func (c *compy) set(addr, val int) {
	masked := val&c.maskOff | c.maskOn
	log.Printf("set(%d): masked(%b) = %b\n", addr, val, masked)
	c.mem[addr] = masked
}

func (c *compy) runLines(lines []string) {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		bits := strings.SplitN(line, " = ", 2)
		cmd, arg := bits[0], bits[1]
		if cmd == "mask" {
			c.setMask(arg)
			continue
		}
		addr, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(cmd, "mem["), "]"))
		if err != nil {
			log.Fatalf("bad address in %q: %v", line, err)
		}
		val, err := strconv.Atoi(arg)
		if err != nil {
			log.Fatalf("bad value in %q: %v", line, err)
		}
		c.set(addr, val)
	}
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("can't read input: %v", err)
	}
	c := newCompy()
	c.runLines(strings.Split(string(data), "\n"))
	fmt.Printf("memory sum: %d\n", c.sumMem())
}
