package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Transformation func(int, interface{}) (interface{}, error)

func AsInts(in []interface{}) []int {
  out := []int{}
  for _, n := range in {
	l, ok := n.(int)
	if !ok {
		log.Fatalf("Couldn't cast %T to int", n)
	}
	out = append(out, l)
  }

  return out
}

func AsStrings(in []interface{}) []string {
  out := []string{}
  for _, n := range in {
	out = append(out, fmt.Sprintf("%s", n))
  }

  return out
}

// ReadTransformedFile reads a text file.
func ReadTransformedFile(path string, fs ...Transformation) []interface{} {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Can't read input: %v", err)
	}
	var data []interface{}
LINE:
	for i, l := range strings.Split(string(text), "\n") {
		var line interface{} = l
		for _, f := range fs {
			line, err = f(i, line)
			if err != nil {
				log.Print(err)
				continue LINE
			}
		}
		data = append(data, line)
	}

	return data
}

// IgnoreBlankLines trims spaces, skips blank lines.
func IgnoreBlankLines(i int, in interface{}) (interface{}, error) {
	l, ok := in.(string)
	if !ok {
		return nil, fmt.Errorf("IgnoreBlankLines expects string, got %T", l)
	}
	if len(strings.TrimSpace(l)) == 0 {
		return "", fmt.Errorf("Skipping blank line %d", i)
	}
	return l, nil
}

// SplitWords splits lines on spaces.
func SplitWords(i int, in interface{}) (interface{}, error) {
	l, ok := in.(string)
	if !ok {
		return nil, fmt.Errorf("SplitWords expects string, got %T", l)
	}
	return strings.Fields(l), nil
}

// Ints transforms lines to ints
func Ints(i int, in interface{}) (interface{}, error) {
	l, ok := in.(string)
	if !ok {
		return nil, fmt.Errorf("Ints expects string, got %T", l)
	}
	n, err := strconv.Atoi(strings.TrimSpace(l))
	if err != nil {
		return 0, fmt.Errorf("Can't convert %q at line %d to number: %v", l, i, err)
	}

	return n, nil
}
