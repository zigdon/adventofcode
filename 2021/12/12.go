package main

import (
	"log"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type set struct {
	Data map[string]bool
}

func emptySet() *set {
	return &set{Data: make(map[string]bool)}
}

func newSet(members ...string) *set {
	s := &set{Data: make(map[string]bool)}
	for _, m := range members {
		s.Add(m)
	}
	return s
}

func (s *set) Add(m string) {
	s.Data[m] = true
}

func (s *set) Extend(ms []string) {
	for _, m := range ms {
		s.Data[m] = true
	}
}

func (s *set) Contains(m string) bool {
	_, ok := s.Data[m]
	return ok
}

func (s *set) List() []string {
	out := []string{}
	for k := range s.Data {
		out = append(out, k)
	}

	return out
}

func (s *set) Copy() *set {
	ns := newSet()
	for k := range s.Data {
		ns.Add(k)
	}

	return ns
}

type caves struct {
	Edges      map[string]*set
	SmallCaves *set
}

func (c *caves) iterate(k string) map[string]bool {
	return c.Edges[k].Data
}

func (c *caves) findPaths(start string, seen *set) []string {
	out := []string{}
	if isSmall(start) { // small cave
		seen.Add(start)
	}

	for e := range c.iterate(start) {
		if seen.Contains(e) {
			continue
		}
		if e == "end" {
			out = append(out, start+","+e)
			continue
		}

		for _, p := range c.findPaths(e, seen.Copy()) {
			out = append(out, start+","+p)
		}
	}

	return out
}

func (c *caves) findLongerPaths() []string {
	s := emptySet()
	for _, small := range c.SmallCaves.List() {
		if small == "start" || small == "end" {
			continue
		}
		s.Extend(c.findPaths2("start", emptySet(), small))
	}
	return s.List()
}

func (c *caves) findPaths2(start string, seen *set, freeSpace string) []string {
	out := []string{}
	if strings.ToLower(start) == start { // small cave
		if start == freeSpace {
			freeSpace = ""
		} else {
			seen.Add(start)
		}
	}

	for e := range c.iterate(start) {
		if seen.Contains(e) {
			continue
		}
		if e == "end" {
			out = append(out, start+","+e)
			continue
		}

		for _, p := range c.findPaths2(e, seen.Copy(), freeSpace) {
			out = append(out, start+","+p)
		}
	}

	return out
}

func isSmall(s string) bool {
	return strings.ToLower(s) == s
}

func readFile(path string) *caves {
	data := common.ReadTransformedFile(path,
		common.IgnoreBlankLines,
		common.Split("-"),
	)
	c := &caves{Edges: make(map[string]*set), SmallCaves: emptySet()}
	for _, l := range data {
		d := l.([]string)
		if e, ok := c.Edges[d[0]]; ok {
			e.Add(d[1])
		} else {
			c.Edges[d[0]] = newSet(d[1])
		}
		if e, ok := c.Edges[d[1]]; ok {
			e.Add(d[0])
		} else {
			c.Edges[d[1]] = newSet(d[0])
		}
		for _, cave := range d {
			if isSmall(cave) {
				c.SmallCaves.Add(cave)
			}
		}
	}
	return c
}

func main() {
	data := readFile("input.txt")
	log.Printf("Counted %d paths", len(data.findPaths("start", emptySet())))
	log.Printf("Counted %d longer paths", len(data.findLongerPaths()))
}
