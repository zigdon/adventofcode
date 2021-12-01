package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Pw struct {
	Min      int
	Max      int
	Req      rune
	Password string
}

func parseData(data string) []Pw {
	var res []Pw
	for _, line := range strings.Split(data, "\n") {
		if len(line) == 0 {
			continue
		}
		p := Pw{}
		bits := strings.Split(line, " ")
		ns := strings.Split(bits[0], "-")
		p.Min, _ = strconv.Atoi(ns[0])
		p.Max, _ = strconv.Atoi(ns[1])
		p.Req = rune(bits[1][0])
		p.Password = bits[2]
		res = append(res, p)
	}

	return res
}

func findInvalid(entries []Pw) []Pw {
	var res []Pw

	for _, ent := range entries {
		counts := make(map[rune]int)
		for _, r := range ent.Password {
			counts[r] = counts[r] + 1
		}
		if counts[ent.Req] > ent.Max || counts[ent.Req] < ent.Min {
			res = append(res, ent)
		}
	}
	return res
}

func findInvalid2(entries []Pw) []Pw {
	var res []Pw
	for _, ent := range entries {
		f := 0
		if rune(ent.Password[ent.Min-1]) == ent.Req {
			f = f + 1
		}
		if rune(ent.Password[ent.Max-1]) == ent.Req {
			f = f + 1
		}
		if f != 1 {
			res = append(res, ent)
		}
	}

	return res
}

func printEnt(e Pw) string {
	return fmt.Sprintf("%d-%d %s %s", e.Min, e.Max, string(e.Req), e.Password)
}

func main() {
	path := os.Args[1]
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	all := parseData(string(data))
	bad := findInvalid(all)
	fmt.Println(len(all) - len(bad))
	bad2 := findInvalid2(all)
	fmt.Println(len(all) - len(bad2))
}
