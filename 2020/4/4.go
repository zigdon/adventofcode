package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	reqFields = []string{"iyr", "byr", "eyr", "ecl", "hcl", "pid", "hgt"} // Not cid
	validEyes = map[string]bool{
		"amb": true, "blu": true, "brn": true, "gry": true,
		"grn": true, "hzl": true, "oth": true,
	}
)

type passport map[string]string

func validateField(key, s string) bool {
	n := atoi(s)
	switch key {
	case "byr":
		return n >= 1920 && n <= 2002
	case "iyr":
		return n >= 2010 && n <= 2020
	case "eyr":
		return n >= 2020 && n <= 2030
	case "hgt":
		n = atoi(s[:len(s)-2])
		if strings.HasSuffix(s, "in") {
			return n >= 59 && n <= 76
		}
		if strings.HasSuffix(s, "cm") {
			return n >= 150 && n <= 193
		}
		return false
	case "hcl":
		if !strings.HasPrefix(s, "#") || len(s) != 7 {
			return false
		}
		for i, r := range s {
			if i == 0 {
				continue
			}
			if (r >= 'a' && r <= 'f') || (r >= '0' && r <= '9') {
				continue
			}
			return false
		}
		return true
	case "ecl":
		_, ok := validEyes[s]
		return ok
	case "pid":
		return len(s) == 9 && n > 0
	case "cid":
		return true
	}
	return false
}

func atoi(a string) int {
	ret, err := strconv.Atoi(a)
	if err != nil {
		return 0
	}
	return ret
}

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
		if !validateField(f, ps[f]) {
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
