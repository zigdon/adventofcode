package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBits(t *testing.T) {
	tests := []struct {
		in         string
		chunksReqs []int
		wantChunks []int
		leftover   int
	}{
		{
			in:         "D2FE28",
			chunksReqs: []int{3, 3, 5, 5, 5},
			wantChunks: []int{6, 4, 23, 30, 5},
			leftover:   3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			b := newBits(tc.in)
			for i, c := range tc.chunksReqs {
				got := b.Consume(c)
				if got != tc.wantChunks[i] {
					t.Errorf("bad consume(%d): want %d, got %d", c, tc.wantChunks[i], got)
				}
			}
			if len(b.bits) != tc.leftover {
				t.Errorf("unexpected number of bits left over: want %d, got %d", tc.leftover, len(b.bits))
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		in               string
		wantValue        int64
		wantPacketValues []int64
	}{
		{
			in:        "D2FE28",
			wantValue: 2021,
		},
		{
			in:               "38006F45291200",
			wantPacketValues: []int64{10, 20},
		},
		{
			in:               "EE00D40C823060",
			wantPacketValues: []int64{1, 2, 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			p := newPacketFromString(tc.in, 0)
			got := p.Value
			if got != tc.wantValue {
				t.Errorf("bad literal value: want %d, got %d", tc.wantValue, got)
			}

			if len(p.Packets) < len(tc.wantPacketValues)-1 {
				t.Fatalf("missing sub-packets, want %d, got %d", len(tc.wantPacketValues), len(p.Packets))
			}
			for i, pv := range tc.wantPacketValues {
				if p.Packets[i].Value != pv {
					t.Errorf("bad sub-packet #%d: want %d, got %d", i, pv, p.Packets[i].Value)
				}
			}
		})
	}
}

func TestPacket(t *testing.T) {
	tests := []struct {
		in          string
		wantVersion int
		wantTypeID  typeID
	}{
		{
			in:          "D2FE28",
			wantVersion: 6,
			wantTypeID:  LITERAL,
		},
		{
			in:          "38006F45291200",
			wantVersion: 1,
			wantTypeID:  OP_LT,
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			p := newPacketFromString(tc.in, 0)
			if p.Version != tc.wantVersion {
				t.Errorf("bad version: want %d, got %d", tc.wantVersion, p.Version)
			}
			if p.TypeID != tc.wantTypeID {
				t.Errorf("bad version: want %s, got %s", tc.wantTypeID, p.TypeID)
			}
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		in  string
		sum int
	}{
		{
			in:  "8A004A801A8002F478",
			sum: 16,
		},
		{
			in:  "620080001611562C8802118E34",
			sum: 12,
		},
		{
			in:  "C0015000016115A2E0802F182340",
			sum: 23,
		},
		{
			in:  "A0016C880162017C3686B18A3D4780",
			sum: 31,
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			p := newPacketFromString(tc.in, 0)
			got := add(p, getVersion)
			if got != tc.sum {
				t.Errorf("bad sum, wanted %d, got %d", tc.sum, got)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		in     string
		values []int64
		out    int64
	}{
		{"C200B40A82", []int64{1, 2}, 3},
		{"04005AC33890", []int64{6, 9}, 54},
		{"880086C3E88112", []int64{7, 8, 9}, 7},
		{"CE00C43D881120", []int64{7, 8, 9}, 9},
		{"D8005AC2A8F0", []int64{5, 15}, 1},
		{"F600BC2D8F", []int64{5, 15}, 0},
		{"9C005AC2F8F0", []int64{5, 15}, 0},
		{"9C0141080250320F1802104A08", []int64{4, 4}, 1},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			p := newPacketFromString(tc.in, 0)
			got := p.Execute()
			pvs := []int64{}
			for _, sp := range p.Packets {
				pvs = append(pvs, sp.Value)
			}
			if diff := cmp.Diff(tc.values, pvs); diff != "" {
				t.Errorf("bad sub values:\n%s", diff)
			}
			if got != tc.out {
				t.Errorf("bad exec: want %d, got %d", tc.out, got)
			}
		})
	}
}
