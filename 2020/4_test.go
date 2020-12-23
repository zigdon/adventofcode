package main

import (
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func samplePassports() []passport {
	return []passport{
		{"ecl": "gry", "pid": "860033327", "eyr": "2020", "hcl": "#fffffd",
			"byr": "1937", "iyr": "2017", "cid": "147", "hgt": "183cm"},
		{"iyr": "2013", "ecl": "amb", "cid": "350", "eyr": "2023", "pid": "028048884",
			"hcl": "#cfa07d", "byr": "1929"},
		{"hcl": "#ae17e1", "iyr": "2013",
			"eyr": "2024",
			"ecl": "brn", "pid": "760753108", "byr": "1931",
			"hgt": "179cm"},
		{"hcl": "#cfa07d", "eyr": "2025", "pid": "166559648",
			"iyr": "2011", "ecl": "brn", "hgt": "59in"},
	}

}

func TestValidatePassport(t *testing.T) {
	sample := samplePassports()
	tests := []struct {
		input map[string]string
		want  bool
	}{
		{
			input: sample[0],
			want:  true,
		},
		{
			input: sample[1],
			want:  false,
		},
		{
			input: sample[2],
			want:  true,
		},
		{
			input: sample[3],
			want:  false,
		},
	}

	for _, tc := range tests {
		got := validatePassport(tc.input)
		if got != tc.want {
			t.Errorf("Bad passport validation for %s", tc.input["pid"])
		}
	}
}

func TestReadPassports(t *testing.T) {
	sample := samplePassports()
	data, err := ioutil.ReadFile("sample.txt")
	if err != nil {
		t.Fatalf("Can't read sample data: %v", err)
	}
	got := readPassports(string(data))
	if diff := cmp.Diff([]passport{sample[0], sample[2]}, got); diff != "" {
		t.Errorf("Bad passports read: -want +got\n%s", diff)
	}
}

func TestValidateField(t *testing.T) {
	tests := []struct {
		field string
		val   string
		want  bool
	}{
		{field: "byr", val: "2002", want: true},
		{field: "byr", val: "2003", want: false},
		{field: "hgt", val: "60in", want: true},
		{field: "hgt", val: "190cm", want: true},
		{field: "hgt", val: "190in", want: false},
		{field: "hgt", val: "190", want: false},
		{field: "hcl", val: "#123abc", want: true},
		{field: "hcl", val: "#123abz", want: false},
		{field: "hcl", val: "123abc", want: false},
		{field: "ecl", val: "brn", want: true},
		{field: "ecl", val: "wat", want: false},
		{field: "pid", val: "000000001", want: true},
		{field: "pid", val: "0123456789", want: false},
	}

	for i, tc := range tests {
		if validateField(tc.field, tc.val) != tc.want {
			t.Errorf("Bad field validation in %d", i)
		}
	}
}
