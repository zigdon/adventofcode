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
	bits    int
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
	c.bits = len(newMask)
	for i := 0; i < c.bits; i++ {
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

func (c *compy) setMaskedVal(addr, val int) {
	masked := val&c.maskOff | c.maskOn
	log.Printf("setMaskedVal(%d): masked(%b) = %b\n", addr, val, masked)
	c.mem[addr] = masked
}

func (c *compy) setMaskedAddr(addr, val int) []int {
	var res []int
	bits := []int{0}
	format := fmt.Sprintf("%%0%db", c.bits)
	masked := []byte(fmt.Sprintf(format, addr))
	log.Printf("addr:   "+format, addr)
	log.Printf("mask:   %s", c.mask)
	for i := c.bits - 1; i >= 0; i-- {
		if c.mask[i] == '1' {
			masked[i] = '1'
		} else if c.mask[i] == 'X' {
			masked[i] = '0'
			bit := 1 << (c.bits - i - 1)
			var newBits []int
			for _, b := range bits {
				newBits = append(newBits, b+bit)
			}
			bits = append(bits, newBits...)
		}
	}
	base, err := strconv.ParseInt(string(masked), 2, 64)
	log.Printf("masked: %s (%d)", masked, base)
	log.Printf("bits: %v", bits)
	if err != nil {
		log.Fatalf("bad mask conversion %q: %v", masked, err)
	}

	for _, a := range bits {
		c.mem[a+int(base)] = val
		res = append(res, a+int(base))
	}

	return res
}

func (c *compy) runLines(lines []string, mode bool) {
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

		if mode {
			c.setMaskedAddr(addr, val)
		} else {
			c.setMaskedVal(addr, val)
		}
	}
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("can't read input: %v", err)
	}
	c := newCompy()
	c.runLines(strings.Split(string(data), "\n"), false)
	fmt.Printf("memory sum: %d\n", c.sumMem())

	c = newCompy()
	c.runLines(strings.Split(string(data), "\n"), true)
	fmt.Printf("memory sum: %d\n", c.sumMem())
}
