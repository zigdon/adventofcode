package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	got, err := readFile("sample.txt")
	if err != nil {
		t.Fatalf("Can't read sample.txt: %v", err)
	}

	want := []command{
		{dir_forward, 5},
		{dir_down, 5},
		{dir_forward, 8},
		{dir_up, 3},
		{dir_down, 8},
		{dir_forward, 2},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Bad reading of file: %s", diff)
	}
}

func TestPilotSub(t *testing.T) {
	cmds, err := readFile("sample.txt")
	if err != nil {
		t.Fatalf("Can't load sample.txt: %v", err)
	}

	gotX, gotY, err := pilotSub(cmds)
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	if gotX != 15 {
		t.Errorf("Bad x: want 15, got %d", gotX)
	}
	if gotY != 10 {
		t.Errorf("Bad y: want 10, got %d", gotY)
	}

}

func TestPilotSubWithAim(t *testing.T) {
	cmds, err := readFile("sample.txt")
	if err != nil {
		t.Fatalf("Can't load sample.txt: %v", err)
	}

	gotX, gotY, err := pilotSubWithAim(cmds)
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	if gotX != 15 {
		t.Errorf("Bad x: want 15, got %d", gotX)
	}
	if gotY != 60 {
		t.Errorf("Bad y: want 60, got %d", gotY)
	}

}
