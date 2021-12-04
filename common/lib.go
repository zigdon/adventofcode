package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Transformation func(int, interface{}) (interface{}, error)

var (
	bufferLine  = fmt.Errorf("Buffer this line")
	flushBuffer = fmt.Errorf("Ignore this line, return the buffer instead")
	ignoreLine  = fmt.Errorf("Transformation not applying to this line")
	skipLine    = fmt.Errorf("Skipping this line")
)

func AsInts(in []interface{}) []int {
	out := []int{}
	for _, n := range in {
		l, err := strconv.Atoi(fmt.Sprintf("%s", n))
		if err != nil {
			log.Printf("*** Couldn't cast %T to int", n)
			continue
		}
		out = append(out, l)
	}

	return out
}

func AsStrings(in []interface{}) []string {
	out := []string{}
	for _, n := range in {
		if fmt.Sprintf("%T", n) != "string" {
			log.Printf("*** Skipping %q", n)
			continue
		}
		out = append(out, fmt.Sprintf("%s", n))
	}

	return out
}

// ReadTransformedFile reads a text file.
func ReadTransformedFile(path string, fs ...Transformation) []interface{} {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("*** Can't read input: %v", err)
		return nil
	}
	var data []interface{}
	var buffer []interface{}
LINE:
	for i, l := range strings.Split(string(text), "\n") {
		var line interface{} = l
		for _, f := range fs {
			transformed, err := f(i, line)
			if err != nil {
				switch err {
				case bufferLine:
					buffer = append(buffer, transformed)
					continue LINE
				case flushBuffer:
					line = buffer
					buffer = []interface{}{}
					continue
				case ignoreLine:
					continue
				case skipLine:
					continue LINE
				default:
					log.Print(err)
				}
				continue LINE
			}
			line = transformed
		}
		data = append(data, line)
	}
	if len(buffer) != 0 {
		data = append(data, buffer)
	}

	return data
}

// Range applies the supplied transformations only to lines between start-end.
// Use -1 to leave open ended for either.
func Range(start, end int, fs ...Transformation) func(i int, in interface{}) (interface{}, error) {
	return func(i int, in interface{}) (interface{}, error) {
		if start >= 0 && i < start {
			return nil, ignoreLine
		}
		if end >= 0 && i >= end {
			return nil, ignoreLine
		}

		line := in
		var err error
		for _, f := range fs {
			line, err = f(i, line)
			if err != nil {
				return nil, err
			}
		}

		return line, nil
	}
}

// Block collects multiple lines into a single item, breaking on empty lines.
func Block(i int, in interface{}) (interface{}, error) {
	// Accept as delimiters empty strings, or empty lists
	l, ok := in.(string)
	if ok && len(l) == 0 {
		return nil, flushBuffer
	}
	ls, ok := in.([]string)
	if ok && len(ls) == 0 {
		return nil, flushBuffer
	}

	return in, bufferLine
}

// IgnoreBlankLines skips blank lines.
func IgnoreBlankLines(i int, in interface{}) (interface{}, error) {
	l, ok := in.(string)
	if !ok {
		return nil, fmt.Errorf("IgnoreBlankLines expects string, got %T", l)
	}
	if len(strings.TrimSpace(l)) == 0 {
		return "", skipLine
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
