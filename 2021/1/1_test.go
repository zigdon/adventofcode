package main

import (
	"testing"
)

func TestCountIncreaseSample(t *testing.T) {
	data := readFile("sample.txt")
	count := countIncrease(data)
	if count != 7 {
		t.Errorf("Bad increase, want 7, got %d", count)
	}
}

func TestCountIncreaseWindowSample(t *testing.T) {
	data := readFile("sample.txt")
	count := countWindowIncrease(data, 3)
	if count != 5 {
		t.Errorf("Bad increase, want 5, got %d", count)
	}
}
