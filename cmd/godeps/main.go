package main

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	deps "github.com/matthewmueller/go-deps"
)

func main() {
	log.SetHandler(cli.Default)

	deps, err := deps.FindWithTests(os.Args[1:]...)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, dep := range deps {
		log.Infof(dep)
	}
}
