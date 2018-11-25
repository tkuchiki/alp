package main

import (
	"log"
	"os"

	"github.com/tkuchiki/alp"
)

func main() {
	p := alp.NewProfiler(os.Stdout, os.Stderr)

	err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
