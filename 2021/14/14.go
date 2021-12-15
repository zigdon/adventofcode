package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type transformation struct {
	Id    string
	Delta int
}

type polymer struct {
	Template        string
	Rules           map[string]string
	Histogram       map[string]int64
	Transformations map[string][]transformation
}

func newPolymer(tmpl string, rules map[string]string) *polymer {
	p := &polymer{
		Template:        tmpl,
		Rules:           rules,
		Histogram:       make(map[string]int64),
		Transformations: make(map[string][]transformation),
	}

	// Count individual letters
	for _, c := range tmpl {
		p.Histogram[string(c)]++
	}

	// Count pairs
	for i := 0; i < len(tmpl)-1; i++ {
		p.Histogram[fmt.Sprintf("%c%c", tmpl[i], tmpl[i+1])]++
	}

	// Each rule takes AB -> ACB. This means we remove one AB, add one each of
	// AC, CB, and add one C.
	for k, v := range rules {
		p.Transformations[k] = []transformation{
			{k, -1},
			{fmt.Sprintf("%c%s", k[0], v), 1},
			{fmt.Sprintf("%s%c", v, k[1]), 1},
			{v, 1},
		}
	}

	return p
}

func (p *polymer) process() {
	// iterate on all the histograms with values > 0
	// apply transformations into a new histogram
	// replace original with new
	nh := make(map[string]int64)
	for k, v := range p.Histogram {
		if v == 0 {
			continue
		}
		if len(k) == 1 { // single letter, just copy it
			nh[k] += v
			continue
		}

		ts, ok := p.Transformations[k]
		if !ok {
			log.Printf("no transformation for %s", k)
			continue
		}
		nh[k] += v
		for _, t := range ts {
			nh[t.Id] += int64(t.Delta) * v
		}
	}

	for k, v := range nh {
		if v == 0 {
			delete(nh, k)
		}
	}

	// log.Printf("%v -> %v\n%s", p.Histogram, nh, cmp.Diff(p.Histogram, nh))

	p.Histogram = nh
}

func readFile(path string) (string, map[string]string) {
	data := common.ReadTransformedFile(
		path,
		common.Range(2, -1, common.Split(" -> ")),
		common.Block,
	)

	rules := map[string]string{}
	tmpl := data[0].([]interface{})[0]
	rs := data[1].([]interface{})
	for _, l := range rs {
		from, to := l.([]string)[0], l.([]string)[1]
		rules[from] = to
	}

	return tmpl.(string), rules
}

func process(t string, rs map[string]string) string {
	next := strings.Builder{}
	for i := 0; i < len(t)-1; i++ {
		fmt.Fprintf(&next, "%c", t[i])
		if add, ok := rs[string(t[i])+string(t[i+1])]; ok {
			fmt.Fprintf(&next, "%s", add)
		} else {
			log.Printf(" *** missing rule for %c%c", t[i], t[i+1])
		}
	}
	fmt.Fprintf(&next, "%c", t[len(t)-1])

	return next.String()
}

type bucket struct {
	Id    rune
	Count int64
}

func histogram(polymer string) []bucket {
	hist := []bucket{}
	counts := map[rune]int64{}

	for _, c := range polymer {
		counts[c]++
	}

	for k, v := range counts {
		hist = append(hist, bucket{k, v})
	}

	sort.Slice(hist, func(i, j int) bool {
		return hist[i].Count > hist[j].Count
	})

	return hist
}

func main() {
	template, rules := readFile("input.txt")
	p := newPolymer(template, rules)
	for i := 0; i < 10; i++ {
		template = process(template, rules)
		log.Printf("After step %d: %d", i, len(template))
	}
	hist := histogram(template)
	l := len(hist) - 1
	log.Printf("%d - %d = %d", hist[0].Count, hist[l].Count, hist[0].Count-hist[l].Count)

	// part 2
	for i := 0; i < 40; i++ {
		log.Printf("after step %d", i)
		p.process()
	}

	var minC, maxC int64
	var minR, maxR string
	for k, v := range p.Histogram {
		if len(k) != 1 {
			continue
		}
		if minC == 0 || v < minC {
			minC = v
			minR = k
		}
		if v > maxC {
			maxC = v
			maxR = k
		}
	}
	log.Printf("%d(%s) - %d(%s) = %d", maxC, maxR, minC, minR, maxC-minC)
}
