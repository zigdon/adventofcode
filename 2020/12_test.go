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

func TestWaypoint(t *testing.T) {
	s := &ship{heading: 90, wayX: 10, wayY: 1}
	gotX, gotY := s.moveToWaypoint(sample())
	if gotX != 214 || gotY != -72 {
		t.Fatalf("ship got lost: want 214, -72, got %d, %d", gotX, gotY)
	}
}
