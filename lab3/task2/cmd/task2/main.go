package main

import (
	"log"

	"automata/common/infrastructure"
	"automata/lab3/task2/pkg/task2/app"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	adapter := infrastructure.NewFiniteInputOutputAdapter()
	determinator := app.NewDeterminator()
	service := app.NewDeterminatorService(adapter, determinator)

	if err = service.Determinate(args.InputFilename, args.OutputFilename); err != nil {
		log.Fatal(err)
	}
}
