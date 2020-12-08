package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/alecthomas/participle"
)

type Syntax struct {
	Rules []*Rule `@@*`
}

type Color struct {
	Color string `@Ident @Ident`
}

type Rule struct {
	Container Color         `@@ "bags" "contain"`
	Empty     *bool         `  ( @"no" "other" "bags" "."`
	Bags      []*CountedBag `| @@ ("," @@ | ".")+ )`
}

type CountedBag struct {
	Count int   `@Int`
	Color Color `@@ ("bag" | "bags")`
}

type Container struct {
	CanContain map[string]bool
	CanBeIn    map[string]bool
}

func getParser() *participle.Parser {
	parser, err := participle.Build(&Syntax{})
	if err != nil {
		log.Fatalf("Error creating parser: %v", err)
	}

	return parser
}

// Find all the colors that can be in a given bag (the supplied color), recursively
func findContents(color string, rules map[string]*Rule, seen map[string]bool) []string {
	if seen == nil {
		seen = make(map[string]bool)
	}

	res := []string{}
	seen[color] = true
	for _, bag := range rules[color].Bags {
		if seen[bag.Color.Color] {
			continue
		}
		res = append(res, bag.Color.Color)
		res = append(res, findContents(bag.Color.Color, rules, seen)...)
	}

	return res
}

func makeMap(syn *Syntax) map[string]*Rule {
	res := make(map[string]*Rule)
	for _, rule := range syn.Rules {
		res[rule.Container.Color] = rule
	}

	return res
}

func getRules(data string) map[string]*Rule {
	parser := getParser()
	syn := &Syntax{}
	err := parser.ParseString("", data, syn)
	if err != nil {
		log.Fatalf("error parsing data: %v", err)
	}

	return makeMap(syn)
}

func parseRules(rules map[string]*Rule) map[string]*Container {
	res := make(map[string]*Container)
	for color := range rules {
		c := &Container{
			CanContain: make(map[string]bool),
			CanBeIn:    make(map[string]bool),
		}
		for _, bag := range findContents(color, rules, nil) {
			c.CanContain[bag] = true
		}
		res[color] = c
	}

	for k, v := range res {
		for color := range v.CanContain {
			res[color].CanBeIn[k] = true
		}
	}

	return res
}

func countContainers(color string, rules map[string]*Container) int {
	res := 0
	for _, cont := range rules {
		if cont.CanContain[color] {
			res++
		}
	}

	return res
}

func main() {
	input := os.Args[1]
	color := os.Args[2]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("can't read input: %v", err)
	}
	rules := getRules(string(data))
	containers := parseRules(rules)
	fmt.Printf("# of bags that can contain %q: %d\n", color, countContainers(color, containers))
}
