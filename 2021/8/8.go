package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type input struct {
	In  []string
	Out []string
}

var digit = map[string]string{
	"ABCEFG":  "0",
	"CF":      "1",
	"ACDEG":   "2",
	"ACDFG":   "3",
	"BCDF":    "4",
	"ABDFG":   "5",
	"ABDEFG":  "6",
	"ACF":     "7",
	"ABCDEFG": "8",
	"ABCDFG":  "9",
}

// mapping how many segments are on -> what segments they might be
var lens = map[int]string{
	2: "CF",
	3: "ACF",
	4: "BCDF",
}

type wireMap struct {
	Possible map[rune]map[rune]bool
	Opts     map[rune]int
	Solved   map[rune]rune
	Pairs    map[rune]string
}

func newWireMap() *wireMap {
	wm := &wireMap{
		Solved:   make(map[rune]rune),
		Possible: make(map[rune]map[rune]bool),
		Opts:     make(map[rune]int),
		Pairs:    make(map[rune]string),
	}
	for _, r := range "ABCDEFG" {
		wm.Possible[r] = make(map[rune]bool)
		for _, p := range "ABCDEFG" {
			wm.Possible[r][p] = true
		}
	}
	return wm
}

func (wm *wireMap) init(pos []string) {
	for i, p := range pos {
		wm.limit(fmt.Sprintf("%c", 'A'+i), p)
	}
}

func (wm *wireMap) check() bool {
	for {
		unchanged := true
		for given, opts := range wm.Possible {
			count := 0
			if _, ok := wm.Solved[given]; ok {
				continue
			}
			var lastOpt rune
			for opt, possible := range opts {
				if possible {
					count++
					lastOpt = opt
				}
			}
			wm.Opts[given] = count
			if count == 1 {
				wm.mark(given, lastOpt)
				unchanged = false
			} else if count == 0 {
				log.Fatalf("no options left for %q! %#v", given, wm)
			}
		}

		if unchanged {
			break
		}
	}

	return len(wm.Solved) == len(wm.Possible)
}

func (wm *wireMap) prune(src string, dst rune) {
	for given, opts := range wm.Possible {
		if strings.ContainsRune(src, given) {
			continue
		}
		for opt := range opts {
			if opt == dst {
				wm.Possible[given][opt] = false
			}
		}
	}
}

func (wm *wireMap) limit(segs string, opts string) {
	for _, s := range segs {
		for given := range wm.Possible[s] {
			if !strings.ContainsRune(opts, given) {
				wm.Possible[s][given] = false
			}
		}
	}
}

func (wm *wireMap) string() string {
	out := "   abcdefg\n"
	for r := 'A'; r <= 'G'; r++ {
		out += fmt.Sprintf("%c: ", r)
		for s := 'A'; s <= 'G'; s++ {
			if wm.Solved[r] == s {
				out += "*"
			} else if wm.Possible[r][s] {
				out += "+"
			} else {
				out += "-"
			}
		}
		out += fmt.Sprintf("  (%d)\n", wm.Opts[r])
	}

	return out
}

func (wm *wireMap) cur(g rune) string {
	out := ""
	for opt, pos := range wm.Possible[g] {
		if pos {
			out = fmt.Sprintf("%s%c", out, opt)
		}
	}

	return out
}

func (wm *wireMap) findPair(name string) {
	var pair, dst string
	for given, opts := range wm.Opts {
		if _, ok := wm.Pairs[given]; ok {
			continue
		}
		if opts == 2 {
			pair = fmt.Sprintf("%s%c", pair, given)
			dst = wm.cur(given)
		}
	}
	if len(pair) != 2 {
		log.Fatalf("Can't find %s pair:\n%s", name, wm.string())
	}
	wm.Pairs[rune(pair[0])] = name
	wm.Pairs[rune(pair[1])] = name
	for _, d := range dst {
		wm.prune(pair, d)
	}
}

func (wm *wireMap) getPair(name string) string {
	out := ""
	for r, n := range wm.Pairs {
		if n != name {
			continue
		}
		out = fmt.Sprintf("%s%c", out, r)
	}

	return out
}

func (wm *wireMap) mapsTo(dst rune) string {
	for s, d := range wm.Solved {
		if d == dst {
			return fmt.Sprintf("%c", s)
		}
	}

	log.Fatalf("nothing mapts to %c", dst)
	return ""
}

func (wm *wireMap) solve(data []string) {
	for _, d := range data {
		if p, ok := lens[len(d)]; ok {
			wm.limit(d, p)
		}
	}
	wm.check()

	// isolate CF (1, 7)
	wm.findPair("CF")
	wm.check()

	// segment A should be identified now

	// isolate BD
	wm.findPair("BD")
	wm.check()

	// isolate EG
	wm.findPair("EG")
	wm.check()

	// find 9 (ACF + BDG) -> identify G (and E)
	pattern := wm.getPair("CF") + wm.getPair("BD") + wm.mapsTo('A') + "."
	match := wm.match(data, pattern)
	for _, r := range match {
		if !strings.ContainsRune(pattern, r) {
			wm.mark(r, 'G')
			break
		}
	}

	// find 3 -> identify D (and B)
	pattern = wm.getPair("CF") + wm.mapsTo('A') + wm.mapsTo('G') + "."
	match = wm.match(data, pattern)
	for _, r := range match {
		if !strings.ContainsRune(pattern, r) {
			wm.mark(r, 'D')
			break
		}
	}

	// find 2 -> identify C (and F)
	pattern = wm.mapsTo('A') + wm.mapsTo('E') + wm.mapsTo('D') + wm.mapsTo('G') + "."
	match = wm.match(data, pattern)
	for _, r := range match {
		if !strings.ContainsRune(pattern, r) {
			wm.mark(r, 'C')
			break
		}
	}

	wm.check()
}

func (wm *wireMap) mark(r rune, dst rune) {
	wm.Solved[r] = dst
	wm.prune(fmt.Sprintf("%c", r), dst)
	wm.check()
}

func (wm *wireMap) match(data []string, pattern string) string {
	targetLen := len(pattern)
	for _, d := range data {
		if len(d) != targetLen {
			continue
		}

		valid := true
		for _, c := range pattern {
			if c == '.' {
				continue
			}

			if !strings.ContainsRune(d, c) {
				valid = false
				break
			}
		}

		if valid {
			return d
		}
	}

	log.Fatalf("can't find a match for %q: %v", pattern, data)
	return ""
}

func (wm *wireMap) decode(data []string) string {
	out := ""
	for _, d := range data {
		key := ""
		decode := []rune{}
		for _, r := range d {
			decode = append(decode, wm.Solved[r])
		}
		sort.Slice(decode, func(i, j int) bool {
			return decode[i] < decode[j]
		})
		for _, r := range decode {
			key = fmt.Sprintf("%s%c", key, r)
		}
		if n, ok := digit[key]; ok {
			out += n
		} else {
			log.Printf("can't decode %q", key)
		}
	}

	return out
}

func readFile(path string) []input {
	data := common.ReadTransformedFile(
		path,
		common.IgnoreBlankLines,
		common.ToUpper,
		common.SplitWords)

	res := []input{}
	for _, l := range data {
		ls := l.([]string)
		res = append(res, input{ls[0:10], ls[11:15]})
	}

	return res
}

func countEasy(data [][]string) int {
	count := 0
	easy := map[int]string{
		2: "1",
		3: "7",
		4: "4",
		7: "8",
	}
	for _, l := range data {
		for _, w := range l {
			if _, ok := easy[len(w)]; ok {
				count++
			}
		}
	}

	return count
}

func main() {
	data := readFile("input.txt")
	outs := [][]string{}
	for _, d := range data {
		outs = append(outs, d.Out)
	}
	log.Printf("easy digits: %d", countEasy(outs))

	sum := 0
	for _, d := range data {
		wm := newWireMap()
		words := []string{}
		words = append(words, d.In...)
		words = append(words, d.Out...)
		wm.solve(words)
		dec := wm.decode(d.Out)
		n, err := strconv.Atoi(dec)
		if err != nil {
			log.Printf("can't make %q into int: %v", dec, err)
		}
		sum += n
	}
	log.Printf("sum: %d", sum)

}
