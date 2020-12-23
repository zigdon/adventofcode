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

func parseInput(path string) (map[string]*validator, *ticket, []*ticket) {
	vs := make(map[string]*validator)
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
			v := newValidator(line)
			vs[v.name] = v
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

func findBadFields(t *ticket, vs map[string]*validator) []int {
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

func identifyFields(vs map[string]*validator, ts []*ticket) []string {
	allValidators := []string{}
	for name := range vs {
		allValidators = append(allValidators, name)
	}

	mapping := make(map[int]map[string]bool)
	for n := range ts[0].fs {
		ruledOut := make(map[string]bool)
		mapping[n] = make(map[string]bool)
		for _, tic := range ts {
			for name, v := range vs {
				if ruledOut[name] {
					continue
				}
				if v.ok(tic.fs[n]) {
					mapping[n][name] = true
					continue
				}
				ruledOut[name] = true
				mapping[n][name] = false
			}
		}
	}
	known := make(map[string]int)

	cont := true
	for cont {
		cont = false
		for num, opts := range mapping {
			single := ""
			for name, possible := range opts {
				if !possible {
					continue
				}
				if _, ok := known[name]; ok {
					continue
				}
				if single == "" {
					single = name
				} else {
					single = ""
					break
				}
			}
			if single != "" {
				log.Printf("column %d must be %s\n\n\n", num, single)
				known[single] = num
				cont = true
				break
			}
		}
	}

	res := make([]string, len(allValidators))
	for k, v := range known {
		res[v] = k
	}

	return res
}

func dumpKnown(mapping map[int]map[string]bool) {
	for i := 0; ; i++ {
		v, ok := mapping[i]
		if !ok {
			break
		}
		fmt.Printf("%d: ", i)
		for n, p := range v {
			if p {
				fmt.Printf("%s ", n)
			}
		}
		fmt.Println()
	}

}

func main() {
	validators, myTicket, otherTickets := parseInput(os.Args[1])
	goodTickets := []*ticket{}
	for _, t := range otherTickets {
		if len(findBadFields(t, validators)) == 0 {
			goodTickets = append(goodTickets, t)
		}
	}

	fields := identifyFields(validators, goodTickets)
	ckSum := 1
	for i, name := range fields {
		if !strings.HasPrefix(name, "departure") {
			continue
		}
		log.Printf("%s: %d", name, myTicket.fs[i])
		ckSum *= myTicket.fs[i]
	}
	fmt.Println(ckSum)
}
