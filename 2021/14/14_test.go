package main

import (
	"sort"
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
		template = process(template, rules)
		if template != w {
			t.Errorf("bad polymer after step %d:\nwant: %q\n got: %q", i, w, template)
		}
	}
}

func TestPolymerProcess(t *testing.T) {
	want := []string{
		"NCNBCHB",
		"NBCCNBBBCBHCB",
		"NBBBCNCCNBBNBNBBCHBHHBCHB",
		"NBBNBNBBCCNBCNCCNBBNBBNBBBNBBNBBCBHCBHHNHCBBCBHCB",
	}

	template, rules := readFile("sample.txt")
	p := newPolymer(template, rules)
	for i, w := range want {
		p.process()
		if diff := cmp.Diff(newPolymer(w, nil).Histogram, p.Histogram); diff != "" {
			t.Errorf("bad polymer after step %d: wanted %q\n%s", i, w, diff)
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

func TestPolymer(t *testing.T) {
	tmpl, rules := readFile("sample.txt")
	got := newPolymer(tmpl, rules)
	want := &polymer{
		Template: tmpl,
		Rules:    rules,
		Histogram: map[string]int{
			"N": 2, "C": 1, "B": 1,
			"NN": 1, "NC": 1, "CB": 1,
		},
	}

	ignoreTransformations := cmp.FilterPath(
		func(p cmp.Path) bool {
			return p.String() == "Transformations"
		}, cmp.Ignore())

	if diff := cmp.Diff(want, got, ignoreTransformations); diff != "" {
		t.Errorf("bad polymer:\n%s", diff)
	}

	// Spot check some transformations

	wantTransformations := map[string][]transformation{
		"CH": {{"CH", -1}, {"B", 1}, {"CB", 1}, {"BH", 1}},
		"BC": {{"BC", -1}, {"B", 1}, {"BB", 1}, {"BC", 1}},
	}

	for k, v := range wantTransformations {
		gt := got.Transformations[k]
		sort.Slice(gt, func(i, j int) bool {
			return gt[i].Id < gt[j].Id
		})
		sort.Slice(v, func(i, j int) bool {
			return v[i].Id < v[j].Id
		})
		if diff := cmp.Diff(v, gt); diff != "" {
			t.Errorf("bad (%s) transformations:\n%s", k, diff)
		}
	}
}
