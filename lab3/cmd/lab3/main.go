package main

import (
	"log"

	"automata/common/infrastructure"
	"automata/lab3/pkg/lab3/app"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	adapter := infrastructure.NewInputOutputAdapter()
	service := app.NewDeterminatorService(adapter)

	if err = service.Determinate(args.InputFilename, args.OutputFilename); err != nil {
		log.Fatal(err)
	}
}
