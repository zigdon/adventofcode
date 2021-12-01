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
		{Answers: []string{"abc"}, Anyone: 3, Everyone: 3},
		{Answers: []string{"a", "b", "c"}, Anyone: 3, Everyone: 0},
		{Answers: []string{"ab", "ac"}, Anyone: 3, Everyone: 1},
		{Answers: []string{"a", "a", "a", "a"}, Anyone: 1, Everyone: 1},
		{Answers: []string{"b"}, Anyone: 1, Everyone: 1},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Bad at parsing: %s", diff)
	}
}

func TestSample(t *testing.T) {
	got := parseInput(sample())
	totalAny := 0
	totalEvery := 0
	for _, f := range got {
		totalAny = totalAny + f.Anyone
		totalEvery = totalEvery + f.Everyone
	}
	if totalAny != 11 {
		t.Errorf("Failed to count to 11, got %d", totalAny)
	}
	if totalEvery != 6 {
		t.Errorf("Failed to count to 6, got %d", totalEvery)
	}
}
