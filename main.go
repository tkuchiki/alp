package main

import (
	"log"

	"github.com/tkuchiki/alp/cmd/alp/cmd"
	buildversion "github.com/tkuchiki/alp/internal/version"
)

var version string

func main() {
	command := cmd.NewCommand(buildversion.Resolve(version))
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
