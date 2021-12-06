package main

import (
	"log"

	"github.com/zigdon/adventofcode/common"
)

type school struct {
	Ages []int64
}

func (s *school) Size() int64 {
	count := int64(0)
	for _, f := range s.Ages {
		count += f
	}

	return count
}

func (s *school) Breed(days int) int64 {
	log.Printf("spawning for %d days", days)
	count := int64(0)
	for days > 0 {
		count = int64(0)
		spawning := s.Ages[0]
		for age, n := range s.Ages[1:] {
			count += n
			s.Ages[age] = n
		}
		s.Ages[6] += spawning
		s.Ages[8] = spawning
		count += 2 * spawning
		days--
		log.Printf("%d days left, %d fish", days, count)
	}

	return count
}

func readFile(path string) *school {
	data := common.AsInts(common.ReadTransformedFile(
		path,
		common.Split(","),
	))

	s := &school{Ages: make([]int64, 9)}

	for _, f := range data {
		s.Ages[f]++
	}
	return s
}

func main() {
	fish := readFile("input.txt")
	fish.Breed(80)
	fish.Breed(256 - 80)
}
