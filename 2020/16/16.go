package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type validator struct {
	name  string
	rules []string
	funcs []func(int) bool
}

func newValidator(line string) *validator {
	bits := strings.Split(line, ": ")
	v := &validator{name: bits[0]}
	for _, r := range strings.Split(bits[1], " or ") {
		v.rules = append(v.rules, r)
		nums := strings.Split(r, "-")
		low, err := strconv.Atoi(nums[0])
		if err != nil {
			log.Fatalf("can't parse range %q in %q: %v", r, line, err)
		}
		high, err := strconv.Atoi(nums[1])
		if err != nil {
			log.Fatalf("can't parse range %q in %q: %v", r, line, err)
		}
		v.funcs = append(
			v.funcs,
			func(n int) bool { return n >= low && n <= high },
		)
	}
	return v
}

func (v *validator) ok(n int) bool {
	for _, f := range v.funcs {
		if f(n) {
			return true
		}
	}

	return false
}

type ticket struct {
	fs []int
}

func newTicket(line string) *ticket {
	t := &ticket{}
	for _, bit := range strings.Split(line, ",") {
		n, err := strconv.Atoi(bit)
		if err != nil {
			log.Fatalf("can't parse ticket %q: %v", line, err)
		}
		t.fs = append(t.fs, n)
	}
	return t
}

func parseInput(path string) ([]*validator, *ticket, []*ticket) {
	vs := []*validator{}
	t := &ticket{}
	ts := []*ticket{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("can't read %q: %v", path, err)
	}

	section := 0
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			section++
			continue
		}
		if section == 0 { // rules
			vs = append(vs, newValidator(line))
		} else if section == 1 { // my ticket
			if strings.Contains(line, ":") {
				continue
			}
			t = newTicket(line)
		} else if section == 2 { // other otherTickets
			if strings.Contains(line, ":") {
				continue
			}
			ts = append(ts, newTicket(line))
		}
	}

	return vs, t, ts
}

func findBadFields(t *ticket, vs []*validator) []int {
	bad := []int{}
	for _, f := range t.fs {
		valid := false
		for _, v := range vs {
			if v.ok(f) {
				valid = true
				break
			}
		}
		if !valid {
			bad = append(bad, f)
		}
	}
	return bad
}

func main() {
	validators, _, otherTickets := parseInput(os.Args[1])
	bad := 0
	for _, t := range otherTickets {
		for _, b := range findBadFields(t, validators) {
			bad += b
		}
	}
	fmt.Println(bad)
}
