package main

import (
  "io/ioutil"
  "log"
  "os"

  "github.com/alecthomas/participle"
)

type Expr struct {
  Left *Value `@@`
  Right []*Ops `@@*`
}

type Value struct {
  Num int   `  @Int`
  Sub *Expr `| "(" @@ ")"`
}

type Ops struct {
  Op string `@("+" | "-" | "*" | "/")`
  Val *Value `@@`
}

func getParser() *participle.Parser {
  p, err := participle.Build(&Expr{})
  if err != nil {
    log.Fatalf("can't build parser: %v", err)
  }
  return p
}

func main() {
  input := os.Args[1]
  _, err := ioutil.ReadFile(input)
  if err != nil {
    log.Fatalf("can't read input: %v", err)
  }

  _ = getParser()
}
