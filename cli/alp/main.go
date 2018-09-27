package main

import (
	"log"

	"github.com/tkuchiki/alp"
)

func main() {
	err := alp.Run()
	if err != nil {
		log.Fatal(err)
	}
}
