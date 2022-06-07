package main

import (
	"log"
	"os"

	"github.com/tkuchiki/alp"
)

var version string

func main() {
	p := alp.NewProfiler(os.Stdout, os.Stderr)
	p.SetVersion(version)

	err := p.Run(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
