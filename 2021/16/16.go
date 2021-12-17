package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type typeID int

const (
	OP_SUM typeID = iota
	OP_PRODUCT
	OP_MINIMUM
	OP_MAXIMUM
	LITERAL
	OP_GT
	OP_LT
	OP_EQ
)

func (t typeID) String() string {
	return []string{
		"OP_SUM",
		"OP_PRODUCT",
		"OP_MINIMUM",
		"OP_MAXIMUM",
		"LITERAL",
		"OP_GT",
		"OP_LT",
		"OP_EQ",
	}[t]
}

type lengthTypeID bool

const (
	TOTAL_LENGTH     lengthTypeID = false
	SUB_PACKET_COUNT lengthTypeID = true
)

func (l lengthTypeID) String() string {
	if l {
		return "TOTAL_LENGTH"
	} else {
		return "SUB_PACKET_COUNT"
	}
}

type bits struct {
	Raw  string
	bits []bool
}

func newBits(in string) *bits {
	h2b := map[rune][]bool{
		'0': {false, false, false, false},
		'1': {false, false, false, true},
		'2': {false, false, true, false},
		'3': {false, false, true, true},
		'4': {false, true, false, false},
		'5': {false, true, false, true},
		'6': {false, true, true, false},
		'7': {false, true, true, true},
		'8': {true, false, false, false},
		'9': {true, false, false, true},
		'A': {true, false, true, false},
		'B': {true, false, true, true},
		'C': {true, true, false, false},
		'D': {true, true, false, true},
		'E': {true, true, true, false},
		'F': {true, true, true, true},
	}
	b := &bits{
		Raw:  in,
		bits: []bool{},
	}

	for _, c := range in {
		bs, ok := h2b[c]
		if !ok {
			log.Fatalf("Can't parse %c in %q!", c, in)
		}
		b.bits = append(b.bits, bs...)
	}

	return b
}

func (b *bits) Consume(n int) int {
	if len(b.bits) < n {
		log.Fatalf("There aren't %d bits left to consume!", n)
	}
	// log.Printf("Consuming %d bits", n)
	out := 0
	bits := ""
	for n > 0 {
		n--
		p := int(math.Pow(2, float64(n)))
		v := b.bits[0]
		b.bits = b.bits[1:]
		if v {
			out += p
			bits += "1"
		} else {
			bits += "0"
		}
	}
	// log.Printf("consumed: %s (%d)", bits, out)

	return out
}

func (b *bits) String() string {
	out := ""
	for _, bit := range b.bits {
		if bit {
			out += "1"
		} else {
			out += "0"
		}
	}
	return out
}

func (b *bits) ParseLiteral(depth int) int64 {
	log.Printf("(%d) parsing %q literal packet:", depth, b.String())
	out := int64(0)
	cont := true
	for {
		out *= 16
		if b.Consume(1) == 0 {
			cont = false
		}
		n := b.Consume(4)
		out += int64(n)
		// log.Printf("byte: %d -> total: %d", n, out)
		if !cont {
			break
		}
	}

	return out
}

type packet struct {
	PacketID     int
	Data         *bits
	Version      int
	TypeID       typeID
	LengthTypeID lengthTypeID
	Packets      []*packet
	Value        int64
}

func newPacketFromString(in string, depth int) *packet {
	log.Printf("Parsing packet: %q", in)
	return newPacket(newBits(in), depth)
}

var lastID = 0

func newPacket(in *bits, depth int) *packet {
	log.Printf("(%d) Parsing packet: %s (%d bits)", depth, in, len(in.bits))
	p := &packet{PacketID: lastID, Data: in}
	lastID++
	p.Version = p.Data.Consume(3)
	p.TypeID = typeID(p.Data.Consume(3))
	log.Printf("(%d) Version: %d, TypeID: %s", depth, p.Version, p.TypeID)

	if p.TypeID == LITERAL {
		p.Value = int64(p.Data.ParseLiteral(depth))
	} else {
		p.ParseOperator(depth)
	}

	return p
}

func (p *packet) String() string {
	out := []string{}

	out = append(out,
		fmt.Sprintf("Packet(%d) {", p.PacketID),
		fmt.Sprintf("  Version: %d", p.Version),
		fmt.Sprintf("  TypeID: %s", p.TypeID),
		fmt.Sprintf("  LengthTypeID: %s", p.LengthTypeID),
		fmt.Sprintf("  Value: %d", p.Value))
	if len(p.Packets) > 0 {
		out = append(out, "  Sub packets: {")
		for _, ps := range p.Packets {
			for _, l := range strings.Split(ps.String(), "\n") {
				out = append(out, "    "+l)
			}
		}
		out = append(out, "  }")
	}
	out = append(out, "}")

	return strings.Join(out, "\n")
}

func (p *packet) ParseOperator(depth int) {
	p.LengthTypeID = lengthTypeID(p.Data.Consume(1) == 1)
	log.Printf("(%d) Operator packet: %s", depth, p.LengthTypeID)
	if p.LengthTypeID == TOTAL_LENGTH {
		length := p.Data.Consume(15)
		sub := &bits{bits: p.Data.bits[:length]}
		p.Data.bits = p.Data.bits[length:]
		log.Printf("(%d) parsing packets for %d bits, %d left over", depth, length, len(p.Data.bits))
		for len(sub.bits) > 6 {
			np := newPacket(sub, depth+1)
			p.Packets = append(p.Packets, np)
			sub = np.Data
		}
	} else {
		count := p.Data.Consume(11)
		log.Printf("(%d) parsing %d packets", depth, count)
		sub := &bits{bits: p.Data.bits[:]}
		for len(p.Packets) < count {
			np := newPacket(sub, depth+1)
			p.Packets = append(p.Packets, np)
			sub = np.Data
		}
		log.Printf("(%d) %d bits left over", depth, len(sub.bits))
		p.Data = sub
	}
}

func (p *packet) Execute() int64 {
	out := int64(0)
	var f func(int, int64)
	switch p.TypeID {
	case OP_SUM:
		f = func(_ int, v int64) {
			out += v
		}
	case OP_PRODUCT:
		f = func(i int, v int64) {
			if i == 0 {
				out = v
			} else {
				out *= v
			}
		}
	case OP_MINIMUM:
		f = func(i int, v int64) {
			if i == 0 || v < out {
				out = v
			}
		}
	case OP_MAXIMUM:
		f = func(i int, v int64) {
			if v > out {
				out = v
			}
		}
	case OP_GT:
		f = func(i int, v int64) {
			if i == 0 {
				out = v
			} else {
				if out > v {
					out = 1
				} else {
					out = 0
				}
			}
		}
	case OP_LT:
		f = func(i int, v int64) {
			if i == 0 {
				out = v
			} else {
				if out < v {
					out = 1
				} else {
					out = 0
				}
			}
		}
	case OP_EQ:
		f = func(i int, v int64) {
			if i == 0 {
				out = v
			} else {
				if out == v {
					out = 1
				} else {
					out = 0
				}
			}
		}

	}

	log.Printf("[%03d] Exec %s on %d packets...", p.PacketID, p.TypeID, len(p.Packets))
	vs := []int64{}
	pids := []int{}
	for i, sp := range p.Packets {
		pids = append(pids, sp.PacketID)
		if len(sp.Packets) > 0 {
			// log.Printf("[%03d] Evaluating sub packet #%d", p.PacketID, i)
			sp.Value = sp.Execute()
			// log.Printf("[%03d] sub packet #%d -> %d", p.PacketID, i, sp.Value)
		}
		vs = append(vs, sp.Value)
		// log.Printf("[%03d] %d %s %d...", p.PacketID, out, p.TypeID, sp.Value)
		f(i, sp.Value)
		// log.Printf("[%03d] ... %d", p.PacketID, out)
	}
	log.Printf("[%03d] Exec %s on %d packets (%v) -> %v -> %d",
		p.PacketID, p.TypeID, len(p.Packets), pids, vs, out)

	return out
}

func getVersion(p *packet) int {
	return p.Version
}

func add(p *packet, f func(*packet) int) int {
	sum := f(p)
	for _, sp := range p.Packets {
		sum += add(sp, f)
	}
	return sum
}

func main() {
	data := common.AsStrings(common.ReadTransformedFile("input.txt",
		common.IgnoreBlankLines))[0]
	p := newPacketFromString(data, 0)
	log.Printf("version sum: %d", add(p, getVersion))
	log.Printf("data:\n%s", p)
	log.Printf("exec: %d", p.Execute())
}
