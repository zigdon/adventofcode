package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

func hasDup(b string) bool {
	for i, c := range b {
		if strings.ContainsRune(b[:i]+b[i+1:], c) {
			return true
		}
	}

	return false
}

func findMarker(data string, bufSize int) int {
	buf := ""
	for i, c := range data {
		if len(buf) < bufSize {
			buf = buf + string(c)
			continue
		} else {
			buf = buf[1:] + string(c)
		}
		if !hasDup(buf) {
			return i + 1
		}
	}

	return 0
}

func one(data string) int {
	return findMarker(data, 4)
}

func two(data string) int {
	return findMarker(data, 14)
}

func readFile(path string) ([]string, error) {
	res := common.ReadTransformedFile(path, common.IgnoreBlankLines)

	return common.AsStrings(res), nil
}

func main() {
	data, err := readFile(os.Args[1])
	if err != nil {
		log.Fatalf("%v", err)
	}

	res := one(data[0])
	fmt.Printf("%v\n", res)

	res = two(data[0])
	fmt.Printf("%v\n", res)
}
