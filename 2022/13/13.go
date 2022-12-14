package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type Packet struct {
	Literal string
	Val     int
	Sub     []*Packet
	Empty   bool
}

func (p *Packet) Value() int {
	if len(p.Sub) > 0 {
		return p.Sub[0].Value()
	}
	return p.Val
}

func (p *Packet) Cmp(p2 *Packet) int {
	log.Printf("Comparing:\n%s\n%s", p, p2)
	if len(p.Sub) == 0 || len(p2.Sub) == 0 {
		if p.Empty && !p2.Empty {
			log.Printf("Short left (empty)")
			return -1
		}
		if !p.Empty && p2.Empty {
			log.Printf("Short right (empty)")
			return 1
		}
		if p.Empty && p2.Empty {
			log.Printf("== Both empty")
			return 0
		}
		v1 := p.Value()
		v2 := p2.Value()
		log.Printf("Ints: %d <? %d", v1, v2)
		if v1 < v2 {
			return -1
		} else if v1 == v2 {
			return 0
		} else {
			return 1
		}
	}
	for i := range p.Sub {
		if len(p2.Sub) <= i {
			log.Printf("Short right: %d vs %d", len(p.Sub), len(p2.Sub))
			return 1
		}
		v := p.Sub[i].Cmp(p2.Sub[i])
		if v == 0 {
			log.Printf("[%d]: %s == %s", i, p.Sub[i], p2.Sub[i])
			continue
		}
		log.Printf("%s <=> %s = %d", p.Sub[i], p2.Sub[i], v)
		return v
	}
	if len(p.Sub) < len(p2.Sub) {
		log.Printf("Short left: %d vs %d", len(p.Sub), len(p2.Sub))
		return -1
	}
	log.Printf("fallthrough: %s vs %s", p, p2)
	return 0
}

func (p *Packet) String() string {
	res := []string{}
	if p.Empty {
		return "E"
	}
	if len(p.Sub) == 0 {
		return fmt.Sprintf("%d", p.Val)
	}
	for _, s := range p.Sub {
		if s.Empty {
			res = append(res, "E")
		} else if len(s.Sub) > 0 {
			res = append(res, s.String())
		} else {
			res = append(res, fmt.Sprintf("%d", s.Val))
		}
	}
	return "[" + strings.Join(res, ",") + "]"
}

type Pair struct {
	Original    string
	Left, Right *Packet
}

func (p Pair) Parse(s string) *Packet {
	pk := &Packet{Literal: s}
	if s == "[]" {
		return &Packet{Literal: "[]", Empty: true}
	}
	// Trim leading and trailing []
	s = s[1 : len(s)-1]
	// log.Printf("%q:", s)
	for _, i := range Segment(s) {
		// log.Printf("-> %q", i)
		if i == "[]" {
			pk.Sub = append(pk.Sub, &Packet{Literal: "[]", Empty: true})
		} else if strings.HasPrefix(i, "[") {
			pk.Sub = append(pk.Sub, p.Parse(i))
		} else {
			pk.Sub = append(pk.Sub, &Packet{Literal: i, Val: common.MustInt(i)})
		}
	}
	return pk
}

func (p Pair) Ordered() bool {
	return p.Left.Cmp(p.Right) < 0
}

func NewPair(l, r string) Pair {
	p := Pair{Original: strings.Join([]string{l, r}, " || ")}
	p.Left = p.Parse(l)
	p.Right = p.Parse(r)
	return p
}

func Segment(s string) []string {
	res := []string{}

	// Scan the string, keep track of [] openings and closings
	d := 0
	start := 0
	for i, c := range s {
		switch c {
		case '[':
			d++
			continue
		case ']':
			d--
			continue
		case ',':
			if d == 0 {
				res = append(res, s[start:i])
				start = i + 1
			}
		}
	}
	if start < len(s) {
		res = append(res, s[start:])
	}

	return res
}

func one(data []Pair) int {
	res := 0
	for i, p := range data {
		log.Printf("=== %d ===: %s", i+1, p.Original)
		if p.Ordered() {
			res += i + 1
			log.Printf("Pair %d ordered: %s", i, p)
		}
	}

	return res
}

func two(data []Pair) int {
	return 0
}

func readFile(path string) []Pair {
	res := common.ReadTransformedFile(path, common.Block)
	p := []Pair{}
	for _, l := range res {
		ps := common.AsStrings(l)
		pair := NewPair(ps[0], ps[1])
		p = append(p, pair)
	}

	return p
}

func main() {
	log.Println("Reading data...")
	data := readFile(os.Args[1])

	log.Println("Part A")
	res := one(data)
	fmt.Printf("%v\n", res)

	log.Println("Part B")
	data = readFile(os.Args[1])
	res = two(data)
	fmt.Printf("%v\n", res)
}
