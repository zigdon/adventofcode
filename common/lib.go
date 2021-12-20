package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const logDebug = false

type Transformation func(int, interface{}) (string, interface{}, error)

var (
	bufferLine  = fmt.Errorf("Buffer this line")
	flushBuffer = fmt.Errorf("Ignore this line, return the buffer instead")
	ignoreLine  = fmt.Errorf("Transformation not applying to this line")
	skipLine    = fmt.Errorf("Skipping this line")
)

func debug(template string, args ...interface{}) {
	if logDebug {
		log.Printf(template, args...)
	}
}

func AsInts(in interface{}) []int {
	out := []int{}
	var list []string
	switch in.(type) {
	case int:
		return []int{in.(int)}
	case []int:
		return in.([]int)
	case string:
		list = append(list, in.(string))
	case []string:
		list = in.([]string)
	case []interface{}:
		for _, i := range in.([]interface{}) {
			out = append(out, AsInts(i)...)
		}
		return out
	default:
		log.Printf("*** Couldn't make ints out of %#v(%T)", in, in)
	}
	for _, n := range list {
		l, err := strconv.Atoi(fmt.Sprintf("%s", n))
		if err != nil {
			log.Printf("*** Couldn't cast %#v(%T) to int", n, n)
		} else {
			out = append(out, l)
		}
	}

	return out
}

func AsIntGrid(in []interface{}) [][]int {
	out := [][]int{}
	for _, l := range in {
		out = append(out, AsInts(l))
	}

	return out
}

func AsStrings(in interface{}) []string {
	out := []string{}
	var list []string
	switch in.(type) {
	case string:
		list = append(list, in.(string))
	case []string:
		list = in.([]string)
	case []interface{}:
		for _, i := range in.([]interface{}) {
			out = append(out, AsStrings(i)...)
		}
		return out
	default:
		log.Printf("*** Couldn't make strings out of %#v(%T)", in, in)
	}
	for _, s := range list {
		out = append(out, s)
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
		debug("line: %q", l)
		var line interface{} = l
		for _, f := range fs {
			name, transformed, err := f(i, line)
			debug(" %s -> %#v", name, transformed)
			if err != nil {
				debug(" ** %v", err)
				switch err {
				case bufferLine:
					debug("Buffering %#v", transformed)
					buffer = append(buffer, transformed)
					continue LINE
				case flushBuffer:
					debug("Flushing %#v", buffer)
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
func Range(start, end int, fs ...Transformation) func(i int, in interface{}) (string, interface{}, error) {
	name := "Range"
	return func(i int, in interface{}) (string, interface{}, error) {
		if start >= 0 && i < start {
			return name, nil, ignoreLine
		}
		if end >= 0 && i >= end {
			return name, nil, ignoreLine
		}

		line := in
		var err error
		for _, f := range fs {
			name, line, err = f(i, line)
			debug("   %s -> %#v", name, line)
			if err != nil {
				return name, line, err
			}
		}

		return name, line, nil
	}
}

// Block collects multiple lines into a single item, breaking on empty lines.
func Block(i int, in interface{}) (string, interface{}, error) {
	name := "Block"
	// Accept as delimiters empty strings, or empty lists
	l, ok := in.(string)
	if ok && len(l) == 0 {
		return name, nil, flushBuffer
	}
	ls, ok := in.([]string)
	if ok && len(ls) == 0 {
		return name, nil, flushBuffer
	}

	return name, in, bufferLine
}

// IgnoreBlankLines skips blank lines.
func IgnoreBlankLines(i int, in interface{}) (string, interface{}, error) {
	name := "IgnoreBlankLines"
	l, ok := in.(string)
	if !ok {
		return name, in, nil
	}
	if len(strings.TrimSpace(l)) == 0 {
		return name, "", skipLine
	}
	return name, l, nil
}

// Split will split on arbitrary strings, dropping empty segments.
func Split(sep string) func(i int, in interface{}) (string, interface{}, error) {
	name := fmt.Sprintf("Split(%q)", sep)
	return func(i int, in interface{}) (string, interface{}, error) {
		switch in.(type) {
		case []string:
			res := [][]string{}
			for _, str := range in.([]string) {
				line := []string{}
				for _, word := range strings.Split(str, sep) {
					if len(word) == 0 {
						continue
					}
					line = append(line, word)
				}
				res = append(res, line)
			}

			return name, res, nil
		case string:
			res := []string{}
			for _, word := range strings.Split(in.(string), sep) {
				if len(word) == 0 {
					continue
				}
				res = append(res, word)
			}
			return name, res, nil
		default:
			return name, in, nil
		}
	}
}

// SplitWords splits lines on spaces.
func SplitWords(i int, in interface{}) (string, interface{}, error) {
	name := "SplitWords"
	l, ok := in.(string)
	if !ok {
		return name, in, nil
	}
	return name, strings.Fields(l), nil
}

// ToUpper makes all strings upper case
func ToUpper(i int, in interface{}) (string, interface{}, error) {
	name := "ToUpper"
	l, ok := in.(string)
	if !ok {
		return name, in, nil
	}
	return name, strings.ToUpper(l), nil
}

// Ints transforms lines to ints
func Ints(i int, in interface{}) (string, interface{}, error) {
	name := "Ints"
	l, ok := in.(string)
	if !ok {
		return name, in, nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(l))
	if err != nil {
		return name, 0, fmt.Errorf("Can't convert %q at line %d to number: %v", l, i, err)
	}

	return name, n, nil
}

// MustInt _will_ convert a string into an int, or die trying.
func MustInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("cant convert %q to int: %v", s, err)
	}

	return n
}

// StringDiff returns 2 lines making it easier to find changes between 2 strings
func StringDiff(s1, s2 string) string {
	l1, l2 := len(s1), len(s2)
	out := []string{"", ""}
	if l1 == l2 {
		return strings.Join([]string{s1, s2}, "\n")
	}
	flipped := false
	if l1 < l2 {
		l1, l2 = l2, l1
		s1, s2 = s2, s1
		flipped = true
	}
	var prefix, suffix int
	for p := range s1 {
		if p >= l2 || s1[p] != s2[p] {
			prefix = p
			break
		}
	}
	for s := range s1 {
		suffix = s
		if l1-s <= prefix {
			break
		}
		if s1[l1-s-1] != s2[l2-s-1] {
			break
		}
	}

	diff := l1 - prefix - suffix
	tmpl := fmt.Sprintf("%%s %%%ds %%s", diff)
	out[0] = fmt.Sprintf(tmpl, s1[:prefix], s1[prefix:l1-suffix], s1[l1-suffix:])
	out[1] = fmt.Sprintf(tmpl, s2[:prefix], s2[prefix:l2-suffix], s2[l2-suffix:])

	if flipped {
		out[0], out[1] = out[1], out[0]
	}

	return strings.Join(out, "\n")
}
