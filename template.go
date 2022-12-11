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

func readFile(path string) []int {
	res := common.ReadTransformedFile(path, common.IgnoreBlankLines)

	return common.AsInts(res)
}

func main() {
    log.Println("Reading data...")
	data := readFile(os.Args[1])

    log.Println("Part A")
	res := one(data)
	fmt.Printf("%v\n", res)

    log.Println("Part B")
	data = readFile(os.Args[1])
	res = two(data)
	fmt.Printf("%v\n", res)
}
