package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

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
	for i := 0; i < 10; i++ {
		template = process(template, rules)
		log.Printf("After step %d: %d", i, len(template))
	}
	hist := histogram(template)
	l := len(hist) - 1
	log.Printf("%d - %d = %d", hist[0].Count, hist[l].Count, hist[0].Count-hist[l].Count)
}
