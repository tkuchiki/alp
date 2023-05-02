package main

import (
	"log"

	"github.com/tkuchiki/alp/cmd/alp/cmd"
)

var version string

func main() {
	if err := cmd.Execute(version); err != nil {
		log.Fatal(err)
	}
}
