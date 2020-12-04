package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	reqFields = []string{"iyr", "byr", "eyr", "ecl", "hcl", "pid", "hgt"} // Not cid
)

type passport map[string]string

func readPassports(data string) []passport {
	var res []passport
	cur := make(passport)
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, ":") {
			if validatePassport(cur) {
				res = append(res, cur)
			}
			cur = make(passport)
			continue
		}

		bits := strings.Split(line, " ")
		for _, bit := range bits {
			kv := strings.Split(bit, ":")
			if len(kv) == 2 {
				cur[kv[0]] = kv[1]
			}
		}
	}

	return res
}

func validatePassport(ps passport) bool {
	for _, f := range reqFields {
		if _, ok := ps[f]; !ok {
			return false
		}
	}

	return true
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("Can't read input: %v", err)
	}
	pps := readPassports(string(data))
	fmt.Printf("%d valid passports", len(pps))
}
