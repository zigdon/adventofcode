package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type Stack struct {
	Cargo []byte
}

func (s *Stack) String() string {
	res := []string{}
	for _, c := range s.Cargo {
		if c > 0 {
			res = append(res, fmt.Sprintf("[%c]", c))
		}
	}

	return strings.Join(res, " ")
}
func (s *Stack) Push(b []byte) {
	s.Cargo = append(s.Cargo, b...)
}
func (s *Stack) Pop(n int) []byte {
	if n > len(s.Cargo) {
		n = len(s.Cargo)
	}
	b := s.Top(n)
	s.Cargo = s.Cargo[:len(s.Cargo)-n]
	return b
}
func (s *Stack) Top(n int) []byte {
	return s.Cargo[len(s.Cargo)-n:]
}

type Instruction struct {
	Qty, From, To int
	Fancy         bool
}

func (i Instruction) Do(ss []*Stack) {
	if i.Fancy {
		ss[i.To-1].Push(ss[i.From-1].Pop(i.Qty))
	} else {
		for n := 0; n < i.Qty; n++ {
			ss[i.To-1].Push(ss[i.From-1].Pop(1))
		}
	}
}

func one(s []*Stack, inst []Instruction) string {
	for _, i := range inst {
		i.Do(s)
	}
	tops := []byte{}
	for _, stack := range s {
		tops = append(tops, stack.Top(1)...)
	}

	return string(tops)
}

func two(s []*Stack, inst []Instruction) string {
	for _, i := range inst {
		i.Fancy = true
		i.Do(s)
	}
	tops := []byte{}
	for _, stack := range s {
		tops = append(tops, stack.Top(1)...)
	}

	return string(tops)
}

func readFile(path string) ([]*Stack, []Instruction, error) {
	res := common.ReadTransformedFile(path, common.Block)
	sStr := common.AsStrings(res[0])
	// Drop the last line, it's just indexes
	sStr = sStr[:len(sStr)-1]
	s := make([]*Stack, len(sStr[0])/4+1)
	for i := range sStr {
		r := sStr[len(sStr)-i-1]
		c := 0
		for c <= len(r)/4 {
			pos := c*4 + 1
			if r[pos] != ' ' {
				if s[c] == nil {
					s[c] = &Stack{Cargo: []byte{}}
				}
				s[c].Cargo = append(s[c].Cargo, r[pos])
			}
			c += 1
		}
	}

	i := []Instruction{}
	for _, r := range common.AsStrings(res[1]) {
		words := strings.Split(r, " ")
		i = append(i, Instruction{
			common.MustInt(words[1]),
			common.MustInt(words[3]),
			common.MustInt(words[5]),
			false,
		})
	}

	return s, i, nil
}

func main() {
	stacks, insts, err := readFile(os.Args[1])
	if err != nil {
		log.Fatalf("%v", err)
	}

	res := one(stacks, insts)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%v\n", res)

	stacks, insts, _ = readFile(os.Args[1])
	res = two(stacks, insts)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%v\n", res)
}
