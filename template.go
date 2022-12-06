package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zigdon/adventofcode/common"
)

func one(data []int) int {
	return 0
}

func two(data []int) int {
	return 0
}

func readFile(path string) ([]int, error) {
	res := common.ReadTransformedFile(path, common.IgnoreBlankLines)

	return common.AsInts(res), nil
}

func main() {
	data, err := readFile(os.Args[1])
	if err != nil {
		log.Fatalf("%v", err)
	}

	res := one(data)
	fmt.Printf("%v\n", res)

	res = two(data)
	fmt.Printf("%v\n", res)
}
