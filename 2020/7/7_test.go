package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sample() string {
	return strings.Join([]string{
		"light red bags contain 1 bright white bag, 2 muted yellow bags.",
		"dark orange bags contain 3 bright white bags, 4 muted yellow bags.",
		"bright white bags contain 1 shiny gold bag.",
		"muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.",
		"shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.",
		"dark olive bags contain 3 faded blue bags, 4 dotted black bags.",
		"vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.",
		"faded blue bags contain no other bags.",
		"dotted black bags contain no other bags.",
	}, "\n")
}

func getSampleSyn() (*Syntax, error) {
	parser := getParser()

	syn := &Syntax{}
	err := parser.ParseString("", sample(), syn)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %v", err)
	}

	return syn, nil
}

func TestParseRules(t *testing.T) {
	got := parseRules(getRules(sample()))
	tests := []struct {
		color string
		cc    []string
		cbi   []string
	}{
		{
			"brightwhite",
			[]string{"darkolive", "dottedblack", "fadedblue", "shinygold", "vibrantplum"},
			[]string{"darkorange", "lightred"},
		},
		{
			"darkolive",
			[]string{"dottedblack", "fadedblue"},
			[]string{"brightwhite", "darkorange", "lightred", "shinygold", "mutedyellow"},
		},
		{
			"darkorange",
			[]string{"dottedblack", "fadedblue", "darkolive", "brightwhite", "shinygold", "vibrantplum", "mutedyellow"},
			[]string{},
		},
	}

	for _, tc := range tests {
		want := &Container{
			CanContain: make(map[string]bool),
			CanBeIn:    make(map[string]bool),
		}
		for _, c := range tc.cc {
			want.CanContain[c] = true
		}
		for _, c := range tc.cbi {
			want.CanBeIn[c] = true
		}
		if diff := cmp.Diff(want, got[tc.color]); diff != "" {
			t.Errorf("bad containers map for %q: %s", tc.color, diff)
		}
	}
}

func TestParser(t *testing.T) {
	got, err := getSampleSyn()
	if err != nil {
		t.Fatal(err)
	}

	True := true
	want := &Syntax{
		Rules: []*Rule{
			&Rule{
				Container: Color{"lightred"},
				Bags: []*CountedBag{
					&CountedBag{Count: 1, Color: Color{"brightwhite"}},
					&CountedBag{Count: 2, Color: Color{"mutedyellow"}},
				}},
			&Rule{
				Container: Color{"darkorange"},
				Bags: []*CountedBag{
					&CountedBag{Count: 3, Color: Color{"brightwhite"}},
					&CountedBag{Count: 4, Color: Color{"mutedyellow"}},
				}},
			&Rule{
				Container: Color{"brightwhite"},
				Bags: []*CountedBag{
					&CountedBag{Count: 1, Color: Color{"shinygold"}},
				}},
			&Rule{
				Container: Color{"mutedyellow"},
				Bags: []*CountedBag{
					&CountedBag{Count: 2, Color: Color{"shinygold"}},
					&CountedBag{Count: 9, Color: Color{"fadedblue"}},
				}},
			&Rule{
				Container: Color{"shinygold"},
				Bags: []*CountedBag{
					&CountedBag{Count: 1, Color: Color{"darkolive"}},
					&CountedBag{Count: 2, Color: Color{"vibrantplum"}},
				}},
			&Rule{
				Container: Color{"darkolive"},
				Bags: []*CountedBag{
					&CountedBag{Count: 3, Color: Color{"fadedblue"}},
					&CountedBag{Count: 4, Color: Color{"dottedblack"}},
				}},
			&Rule{
				Container: Color{"vibrantplum"},
				Bags: []*CountedBag{
					&CountedBag{Count: 5, Color: Color{"fadedblue"}},
					&CountedBag{Count: 6, Color: Color{"dottedblack"}},
				}},
			&Rule{
				Container: Color{"fadedblue"},
				Empty:     &True,
			},
			&Rule{
				Container: Color{"dottedblack"},
				Empty:     &True,
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bad parsing: %s", diff)
	}
}

func TestMakeMap(t *testing.T) {
	syn, err := getSampleSyn()
	if err != nil {
		t.Fatal(err)
	}

	got := makeMap(syn)
	for color, rule := range got {
		if color != rule.Container.Color {
			t.Errorf("Wrong rule for %q: %v", color, rule)
		}
	}
}

func TestFindContents(t *testing.T) {
	syn, err := getSampleSyn()
	if err != nil {
		t.Fatal(err)
	}

	rules := makeMap(syn)

	tests := []struct {
		color string
		want  []string
	}{
		{"dottedblack", []string{}},
		{"vibrantplum", []string{"dottedblack", "fadedblue"}},
		{"shinygold", []string{"darkolive", "dottedblack", "fadedblue", "vibrantplum"}},
	}

	for _, tc := range tests {
		got := findContents(tc.color, rules, nil)
		sort.Strings(got)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("bad contents for %q: %s", tc.color, diff)
		}
	}
}

func TestCountContainers(t *testing.T) {
	containers := parseRules(getRules(sample()))
	got := countContainers("shinygold", containers)
	if got != 4 {
		t.Errorf("misplaced bag: expected 4, found %d", got)
	}
}

func TestContContents(t *testing.T) {
	rules := getRules(sample())

	tests := []struct {
		color string
		want  int
	}{
		{"dottedblack", 0},
		{"vibrantplum", 11},
		{"shinygold", 32},
	}

	for _, tc := range tests {
		got := countContents(tc.color, rules)
		if got != tc.want {
			t.Errorf("someone messed with our %q bag, expected %d, got %d", tc.color, tc.want, got)
		}
	}
}
