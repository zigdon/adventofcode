package main

import (
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

func TestParseRules(t *testing.T) {
	tests := []struct {
		text string
		want map[string]rule
	}{
		{
			text: strings.Join([]string{
				"dotted black bags contain no other bags.",
			}, "\n"),
			want: map[string]rule{
				"dotted black": {
					CanContain:   map[string]int{},
					CanBeFoundIn: map[string]bool{},
				},
			},
		},
		{
			text: strings.Join([]string{
				"vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.",
				"faded blue bags contain no other bags.",
				"dotted black bags contain no other bags.",
			}, "\n"),
			want: map[string]rule{
				"vibrant plum": {
					CanContain: map[string]int{
						"faded blue":   5,
						"dotted black": 6,
					},
					CanBeFoundIn: map[string]bool{},
				},
				"faded blue": {
					CanContain:   map[string]int{},
					CanBeFoundIn: map[string]bool{},
				},
				"dotted black": {
					CanContain:   map[string]int{},
					CanBeFoundIn: map[string]bool{},
				},
			},
		},
		{
			text: strings.Join([]string{
				"shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.",
				"dark olive bags contain 3 faded blue bags, 4 dotted black bags.",
				"vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.",
				"faded blue bags contain no other bags.",
				"dotted black bags contain no other bags.",
			}, "\n"),
			want: map[string]rule{
				"shiny gold": {
					CanContain: map[string]int{
						"dark olive":   1,
						"vibrant plum": 2,
					},
					CanBeFoundIn: map[string]bool{},
				},
				"dark olive": {
					CanContain: map[string]int{
						"faded blue":   3,
						"dotted black": 4,
					},
					CanBeFoundIn: map[string]bool{},
				},
				"vibrant plum": {
					CanContain: map[string]int{
						"faded blue":   5,
						"dotted black": 6,
					},
					CanBeFoundIn: map[string]bool{},
				},
				"faded blue": {
					CanContain:   map[string]int{},
					CanBeFoundIn: map[string]bool{},
				},
				"dotted black": {
					CanContain:   map[string]int{},
					CanBeFoundIn: map[string]bool{},
				},
			},
		},
	}

	for i, tc := range tests {
		got := parseRules(tc.text)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("Bad rule #%d parsing: %s", i, diff)
		}
	}
}
