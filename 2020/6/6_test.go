package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sample() string {
	data, err := ioutil.ReadFile("sample.txt")
	if err != nil {
		log.Fatalf("Error reading sample: %v", err)
	}
	return string(data)
}

func TestParseInput(t *testing.T) {
	got := parseInput(sample())
	want := []*form{
		{Answers: []string{"abc"}, Anyone: 3},
		{Answers: []string{"a", "b", "c"}, Anyone: 3},
		{Answers: []string{"ab", "ac"}, Anyone: 3},
		{Answers: []string{"a", "a", "a", "a"}, Anyone: 1},
		{Answers: []string{"b"}, Anyone: 1},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Bad at parsing: %s", diff)
	}
}

func TestSample(t *testing.T) {
	got := parseInput(sample())
	total := 0
	for _, f := range got {
		total = total + f.Anyone
	}
	if total != 11 {
		t.Errorf("Failed to count to 11, got %d", total)
	}
}
