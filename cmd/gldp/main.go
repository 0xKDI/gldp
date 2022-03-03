package main

import (
	"log"
	"github.com/0qq/gldp/pkg/parser"
)


func main() {
	p := parser.Parser{}
	err := p.Run()

	if err != nil {
		log.Fatal(err)
	}
}
