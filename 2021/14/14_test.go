package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	template, rules := readFile("sample.txt")

	if template != "NNCB" {
		t.Errorf("bad template: want NNCB, got %q", template)
	}

	want := map[string]string{
		"CH": "B",
		"HH": "N",
		"CB": "H",
		"NH": "C",
		"HB": "C",
		"HC": "B",
		"HN": "C",
		"NN": "C",
		"BH": "H",
		"NC": "B",
		"NB": "B",
		"BN": "B",
		"BB": "N",
		"BC": "B",
		"CC": "N",
		"CN": "C",
	}

	if diff := cmp.Diff(want, rules); diff != "" {
		t.Errorf("bad rules:\n%s", diff)
	}
}

func TestProcess(t *testing.T) {
	want := []string{
		"NCNBCHB",
		"NBCCNBBBCBHCB",
		"NBBBCNCCNBBNBNBBCHBHHBCHB",
		"NBBNBNBBCCNBCNCCNBBNBBNBBBNBBNBBCBHCBHHNHCBBCBHCB",
	}

	template, rules := readFile("sample.txt")
	for i, w := range want {
		t.Logf("starting %d: %s", i, template)
		template = process(template, rules)
		t.Logf(" ->       : %s", template)
		if template != w {
			t.Errorf("bad polymer after step %d:\nwant: %q\n got: %q", i, w, template)
		}
	}
}

func TestHistogram(t *testing.T) {
	template, rules := readFile("sample.txt")
	for i := 0; i < 10; i++ {
		template = process(template, rules)
	}
	if len(template) != 3073 {
		t.Errorf("bad length after 10 steps, want 3073, got %d", len(template))
	}

	hist := histogram(template)
	if diff := cmp.Diff([]bucket{{'B', 1749}, {'H', 161}}, []bucket{hist[0], hist[len(hist)-1]}); diff != "" {
		t.Errorf("bad histogram:\n%s", diff)
	}
}
