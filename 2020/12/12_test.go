package main

import "testing"

func sample() []string {
	return []string{
		"F10",
		"N3",
		"F7",
		"R90",
		"F11",
	}
}

func TestTrack(t *testing.T) {
	s := &ship{heading: 90}
	gotX, gotY := s.track(sample())
	if gotX != 17 || gotY != -8 {
		t.Fatalf("ship got lost: want 17, -8, got %d, %d", gotX, gotY)
	}
}
