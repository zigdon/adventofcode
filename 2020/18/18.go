package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/alecthomas/participle"
)

type Expression struct {
	Subs Expression `"(" @@ ")"`
	Val  string     `|| @Int ( @Operator @Expression )*`
}

type Operator struct {
	Op string `@( "+" | "-" | "*" | "/" )`
}

func getParser() *participle.Parser {
	parser, err := participle.Build(&Syntax{})
	if err != nil {
		log.Fatalf("Error creating parser: %v", err)
	}

	return parser
}

func main() {
	input := os.Args[1]
	_, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("can't read input: %v", err)
	}
}
