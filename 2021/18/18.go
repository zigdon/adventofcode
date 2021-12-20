package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

var debugLevel = 0

func debug(lvl int, tmpl string, args ...interface{}) {
	if lvl > debugLevel {
		return
	}
	tmpl = fmt.Sprintf("[%d] %s", lvl, tmpl)
	log.Printf(tmpl, args...)
}

type direction bool

const (
	LEFT  direction = false
	RIGHT direction = true
)

func (d direction) String() string {
	if d {
		return "RIGHT"
	} else {
		return "LEFT"
	}
}

type snail struct {
	LeftVal, RightVal     int
	LeftSnail, RightSnail *snail
	Parent                *snail
}

func (s *snail) String() string {
	var l, r string
	if s.LeftSnail == nil {
		l = fmt.Sprint(s.LeftVal)
	} else {
		l = s.LeftSnail.String()
	}
	if s.RightSnail == nil {
		r = fmt.Sprint(s.RightVal)
	} else {
		r = s.RightSnail.String()
	}

	return fmt.Sprintf("[%s,%s]", l, r)
}

func (s *snail) getValue(d direction) int {
	if d == LEFT {
		return s.LeftVal
	} else {
		return s.RightVal
	}
}

func (s *snail) addValue(val int, d direction) {
	if d == LEFT {
		s.LeftVal += val
	} else {
		s.RightVal += val
	}
}

func (s *snail) getSnail(d direction) *snail {
	if d == LEFT {
		return s.LeftSnail
	} else {
		return s.RightSnail
	}
}

func (s *snail) addAdjecent(d direction) {
	// go up until we didn't come from direction d, then go down one level in
	// direction d, then !d until you find a value in !d
	val := s.getValue(d)
	debug(3, "adding %d from %s of %s(%p)", val, d, s, s)
	cur := s
	prev := s
	for {
		cur = cur.Parent
		debug(3, "up to %s(%p)", cur, cur)
		if cur == nil {
			return
		}
		if cur.getSnail(d) != prev {
			break
		}

		prev = cur
	}

	if cur.getSnail(d) == nil {
		debug(3, "-> adding %d to %s of %s(%p)", val, d, cur, cur)
		cur.addValue(val, d)
		return
	}
	cur = cur.getSnail(d)

	for {
		if cur.getSnail(!d) == nil {
			debug(3, "-> adding %d to %s of %s (%p)", val, !d, cur, cur)
			cur.addValue(val, !d)
			return
		}
		cur = cur.getSnail(!d)
		debug(3, "down to %s", cur)
	}
}

func (s *snail) explode() bool {
	for _, l1 := range []*snail{s.LeftSnail, s.RightSnail} {
		if l1 == nil {
			continue
		}
		for _, l2 := range []*snail{l1.LeftSnail, l1.RightSnail} {
			if l2 == nil {
				continue
			}
			for _, l3 := range []*snail{l2.LeftSnail, l2.RightSnail} {
				if l3 == nil || (l3.LeftSnail == nil && l3.RightSnail == nil) {
					continue
				}
				if l3.LeftSnail != nil {
					debug(2, "Exploding %s in %s", l3.LeftSnail, s)
					l3.LeftSnail.addAdjecent(LEFT)
					l3.LeftSnail.addAdjecent(RIGHT)
					l3.LeftSnail = nil
					l3.LeftVal = 0
				} else {
					debug(2, "Exploding %s in %s", l3.RightSnail, s)
					l3.RightSnail.addAdjecent(LEFT)
					l3.RightSnail.addAdjecent(RIGHT)
					l3.RightSnail = nil
					l3.RightVal = 0
				}
				return true
			}
		}
	}

	return false
}

func (s *snail) checkSplit() bool {
	debug(3, "checking %s for splits", s)
	for _, d := range []direction{LEFT, RIGHT} {
		if s.getValue(d) > 9 {
			debug(2, "splitting %s of %s", d, s)
			l := s.getValue(d) / 2
			r := l
			if l*2 != s.getValue(d) {
				r++
			}
			ns, _ := makeSnail(fmt.Sprintf("[%d,%d]", l, r), s)
			if d == LEFT {
				s.LeftVal = 0
				s.LeftSnail = ns
			} else {
				s.RightVal = 0
				s.RightSnail = ns
			}
			return true
		}

		if s.getSnail(d) != nil && s.getSnail(d).checkSplit() {
			return true
		}
	}
	return false
}

func (s *snail) reduce() {
	var before string
	for {
		if before != "" {
			debug(2, "reduced:\n%s", common.StringDiff(before, s.String()))
		}
		before = s.String()
		if s.explode() {
			continue
		}
		if s.checkSplit() {
			continue
		}

		return
	}
}

func (s *snail) magnitude() int {
	out := 0
	if s.LeftSnail != nil {
		out += 3 * s.LeftSnail.magnitude()
	} else {
		out += 3 * s.LeftVal
	}
	if s.RightSnail != nil {
		out += 2 * s.RightSnail.magnitude()
	} else {
		out += 2 * s.RightVal
	}

	return out
}

func (s *snail) add(sb *snail) *snail {
	debug(1, "\n\n***** Adding %s + %s", s, sb)
	ns, _ := makeSnail(fmt.Sprintf("[%s,%s]", s, sb), nil)
	ns.reduce()
	return ns
}

func consumePrefix(s string, p string) string {
	if !strings.HasPrefix(s, p) {
		log.Fatalf("missing %s in %q!", p, s)
	}
	return strings.TrimPrefix(s, p)
}

func consumeNumber(s string) (int, string) {
	n := 0
	for s[0] >= '0' && s[0] <= '9' {
		n *= 10
		n += int(s[0] - '0')
		s = s[1:]
	}

	return n, s
}

func makeSnail(in string, parent *snail) (*snail, string) {
	sn := &snail{Parent: parent}
	in = consumePrefix(in, "[")
	if strings.HasPrefix(in, "[") {
		sn.LeftSnail, in = makeSnail(in, sn)
	} else {
		sn.LeftVal, in = consumeNumber(in)
	}

	in = consumePrefix(in, ",")

	if strings.HasPrefix(in, "[") {
		sn.RightSnail, in = makeSnail(in, sn)
	} else {
		sn.RightVal, in = consumeNumber(in)
	}

	in = consumePrefix(in, "]")

	return sn, in
}

func readFile(path string) []*snail {
	data := common.ReadTransformedFile(
		path,
		common.IgnoreBlankLines,
	)

	out := []*snail{}
	for _, s := range data {
		sn, leftover := makeSnail(s.(string), nil)
		if leftover != "" {
			log.Fatalf("unbalanced snail %q!", s)
		}
		out = append(out, sn)
	}

	return out
}

func main() {
	ss := readFile("input.txt")
	s := ss[0]
	// fmt.Printf("  %s\n", s)
	for _, sb := range ss[1:] {
		s = s.add(sb)
		// fmt.Printf("+ %s\n= %s\n", sb, s)
	}
	fmt.Printf("magnitude: %d\n", s.magnitude())

	max := 0
	var ma, mb *snail
	for i, sa := range ss {
		for j, sb := range ss {
			if i == j {
				continue
			}
			m := sa.add(sb).magnitude()
			if m > max {
				log.Printf("%s + %s -> %d", sa, sb, m)
				max = m
				ma, mb = sa, sb
			}
		}
	}

	log.Printf("max: %d (%s + %s)", max, ma, mb)
}
