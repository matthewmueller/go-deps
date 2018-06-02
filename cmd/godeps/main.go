package main

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/matthewmueller/deps"
)

func main() {
	log.SetHandler(cli.Default)

	deps, err := deps.Find(os.Args[1:]...)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, dep := range deps {
		log.Infof(dep)
	}
}
