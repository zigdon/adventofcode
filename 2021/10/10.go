package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/zigdon/adventofcode/common"
)

type stack struct {
	items []rune
}

func (s *stack) push(i rune) {
	s.items = append(s.items, i)
}

func (s *stack) pop() rune {
	l := len(s.items) - 1
	if l < 0 {
		log.Fatal("Can't pop from an empty stack.")
	}
	i := s.items[l]
	s.items = s.items[:l]

	return i
}

func (s *stack) last() rune {
	if len(s.items) == 0 {
		return 0
	}
	return s.items[len(s.items)-1]
}

func (s *stack) print() string {
	if len(s.items) == 0 {
		return "<empty>"
	}
	out := fmt.Sprintf("%c", s.items[0])
	for _, r := range s.items[1:] {
		out = fmt.Sprintf("%s, %c", out, r)
	}
	return out
}

func readFile(path string) []string {
	return common.AsStrings(common.ReadTransformedFile(
		path,
		common.IgnoreBlankLines,
	))
}

var opener = map[rune]rune{
	'[': ']',
	'(': ')',
	'{': '}',
	'<': '>',
}

var closer = map[rune]rune{
	']': '[',
	')': '(',
	'}': '{',
	'>': '<',
}

var value = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

var cValue = map[rune]int64{
	')': 1,
	']': 2,
	'}': 3,
	'>': 4,
}

func isValid(l string) (int, rune) {
	s := &stack{}
	for i, r := range l {
		// log.Printf("%s -- %c", s.print(), r)
		if m, ok := closer[r]; ok {
			if m != s.last() {
				return i, r
			}
			s.pop()
			continue
		}
		s.push(r)
	}

	return 0, 0
}

func scoreSet(ls []string) int {
	score := 0
	for _, l := range ls {
		_, r := isValid(l)
		if r == 0 {
			continue
		}
		v, ok := value[r]
		if !ok {
			log.Fatalf("no value listed for %c in %q", r, l)
		}
		score += v
	}

	return score
}

func complete(l string) (string, int64) {
	out := ""
	score := int64(0)
	s := &stack{}
	for i, r := range l {
		if m, ok := closer[r]; ok {
			if m != s.last() {
				log.Fatalf("incomplete line %q was actually invalid at %d,%c", l, i, r)
			}
			s.pop()
			continue
		}
		s.push(r)
	}

	for s.last() != 0 {
		score *= 5
		r := s.pop()
		o, ok := opener[r]
		if !ok {
			log.Fatalf("can't close %c!", r)
		}
		out += fmt.Sprintf("%c", o)
		v, ok := cValue[o]
		if !ok {
			log.Fatalf("can't value %c!", o)
		}
		score += v
	}

	return out, score
}

func main() {
	data := readFile("input.txt")
	log.Printf("score: %d", scoreSet(data))

	incomplete := []string{}
	for _, l := range data {
		if _, r := isValid(l); r == 0 {
			incomplete = append(incomplete, l)
		}
	}

	vals := []int64{}
	for _, l := range incomplete {
		_, v := complete(l)
		vals = append(vals, v)
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	log.Printf("%v\n\n%d", vals, vals[len(vals)/2])
}
