package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	data := readFile("sample.txt")
	if len(data) != 5 {
		t.Errorf("expected 5 scanners, got %d", len(data))
	}
	tests := []struct {
		scanner int
		beacons []coord
	}{
		{
			scanner: 0,
			beacons: []coord{
				{404, -588, -901},
				{528, -643, 409},
				{-838, 591, 734},
				{390, -675, -793},
			},
		},
		{
			scanner: 1,
			beacons: []coord{
				{-336, 658, 858},
				{95, 138, 22},
				{-476, 619, 847},
				{-340, -569, -846},
				{567, -361, 727},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("sampleing scanner #%d", tc.scanner), func(t *testing.T) {
			s := data[tc.scanner]
			for _, b := range tc.beacons {
				if !s.Beacons[b] {
					t.Errorf("expected to find %s in scanner #%d", b, tc.scanner)
				}
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	data := readFile("sample.txt")
	for i := 0; i < 10; i++ {
		data2 := readFile("sample.txt")
		if diff := cmp.Diff(data[0], data2[0]); diff != "" {
			t.Errorf("bad at reading the same thing:\n%s", diff)
		}
	}
}

func TestRotateCoord(t *testing.T) {
	tests := []struct {
		x, y, z int
		want    coord
	}{
		{
			want: coord{2, 3, 4},
		},
		{
			z:    1,
			want: coord{3, -2, 4},
		},
		{
			z:    2,
			want: coord{-2, -3, 4},
		},
		{
			z:    3,
			want: coord{-3, 2, 4},
		},
		{
			x:    1,
			want: coord{2, 4, -3},
		},
		{
			x: 2, z: 1,
			want: coord{-3, -2, -4},
		},
		{
			x: 3, z: 2,
			want: coord{-2, 4, 3},
		},
		{
			z:    3,
			want: coord{-3, 2, 4},
		},
		{
			x: 3, y: 1, z: 2,
			want: coord{-3, 4, -2},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d/%d/%d", tc.x, tc.y, tc.z), func(t *testing.T) {
			c := coord{2, 3, 4}
			got := c.turn(coord{tc.x, tc.y, tc.z})
			if tc.want.X != got.X || tc.want.Y != got.Y || tc.want.Z != got.Z {
				t.Errorf("bad rotation: want %s, got %s", tc.want, got)
			}
		})
	}
}

func TestRotate(t *testing.T) {
	testScanner := func() *scanner {
		s := newScanner(0, []string{
			"-1,-1,1",
			"-2,-2,2",
			"-3,-3,3",
			"-2,-3,1",
			"5,6,-4",
			"8,0,7"})
		return s
	}

	tests := []struct {
		x, y, z int
		want    []coord
	}{
		{
			want: []coord{
				{-1, -1, 1},
				{-2, -2, 2},
				{-3, -3, 3},
				{-2, -3, 1},
				{5, 6, -4},
				{8, 0, 7},
			},
		},
		{
			y: 3,
			want: []coord{
				{-1, -1, -1},
				{-2, -2, -2},
				{-3, -3, -3},
				{-1, -3, -2},
				{4, 6, 5},
				{-7, 0, 8},
			},
		},
		{
			x: 1, y: 2,
			want: []coord{
				{1, 1, -1},
				{2, 2, -2},
				{3, 3, -3},
				{2, 1, -3},
				{-5, -4, 6},
				{-8, 7, -0},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d/%d/%d", tc.x, tc.y, tc.z), func(t *testing.T) {
			s := testScanner()
			s.turn(coord{tc.x, tc.y, tc.z})
			missing := false
			for _, c := range tc.want {
				if !s.Beacons[c] {
					t.Errorf("exected to find %s!", c)
					missing = true
				}
			}
			if missing {
				t.Log(s)
			}
		})
	}
}

func Test24(t *testing.T) {
	s := newScanner(0, []string{
		"1,2,3",
		"4,-5,6",
	})

	seen := make(map[coord]bool)
	start := s.Orientation
	s.try24(func(s *scanner, _ coord) {
		seen[s.Orientation] = true
	})

	if len(seen) != 24 {
		t.Errorf("Only tries %d orientations!", len(seen))
	}
	if !s.Orientation.eq(start) {
		t.Errorf("Ended in the wrong orientation, want %s, got %s", start, s.Orientation)
	}
}

func TestDeltas(t *testing.T) {
	tests := []struct {
		in   *scanner
		want map[coord]coord
	}{
		{
			in: newScanner(0, []string{
				"1,1,1",
				"2,2,2",
				"-1,2,0",
			}),
			want: map[coord]coord{
				{-3, 0, -2}:  {-1, 2, 0},
				{1, 1, 1}:    {2, 2, 2},
				{-2, 1, -1}:  {-1, 2, 0},
				{-1, -1, -1}: {1, 1, 1},
				{2, -1, 1}:   {1, 1, 1},
				{3, 0, 2}:    {2, 2, 2},
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			if diff := cmp.Diff(tc.want, tc.in.Deltas); diff != "" {
				t.Errorf("bad deltas:\n%s", diff)
			}
		})
	}
}

func TestBestMatch(t *testing.T) {
	data := readFile("sample.txt")
	tests := []struct {
		req         int
		a, b        *scanner
		setOrigin   []coord
		wantMatches []coord
		wantDir     coord
		wantOrigin  coord
	}{
		{
			req:        12,
			a:          data[0],
			b:          data[1],
			wantOrigin: coord{68, -1246, -43},
			wantMatches: []coord{
				{-661, -816, -575},
				{-618, -824, -621},
				{-537, -823, -458},
				{-485, -357, 347},
				{-447, -329, 318},
				{-345, -311, 381},
				{390, -675, -793},
				{404, -588, -901},
				{423, -701, 434},
				{459, -707, 401},
				{528, -643, 409},
				{544, -627, -890},
			},
		},
		{
			req:        12,
			a:          data[1],
			b:          data[4],
			setOrigin:  []coord{{}, {68, -1246, -43}},
			wantOrigin: coord{-20, -1133, 1061},
			wantMatches: []coord{
				{459, -707, 401},
				{-739, -1745, 668},
				{-485, -357, 347},
				{432, -2009, 850},
				{528, -643, 409},
				{423, -701, 434},
				{-345, -311, 381},
				{408, -1815, 803},
				{534, -1912, 768},
				{-687, -1600, 576},
				{-447, -329, 318},
				{-635, -1737, 486},
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			if len(tc.setOrigin) > 0 {
				for i, o := range tc.setOrigin {
					if !data[i].Origin.eq(o) {
						data[i].shift(o)
						t.Logf("setting origin of %d to %s", i, o)
					}
				}
			}
			cnt, ori, shift, matches := tc.a.align(tc.b)
			if cnt != len(tc.wantMatches) {
				t.Errorf("bad count: want %d, got %d", len(tc.wantMatches), cnt)
			}
			if !tc.wantDir.isEmpty() && !ori.eq(tc.wantDir) {
				t.Errorf("bad orientation: want %s, got %s", tc.wantDir, ori)
			}
			if !tc.wantOrigin.isEmpty() && !shift.eq(tc.wantOrigin) {
				t.Errorf("bad origin: want %s, got %s", tc.wantOrigin, shift)
			}
			want := make(map[coord]bool)
			got := make(map[coord]bool)
			for _, k := range tc.wantMatches {
				want[k] = true
			}
			for _, k := range matches {
				got[k] = true
			}

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("bad matches:\n%s", diff)
			}
		})
	}
}
