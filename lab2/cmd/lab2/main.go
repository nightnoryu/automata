package main

import (
	"log"

	"automata/common/infrastructure"
	"automata/lab2/pkg/lab2/app"
)

const (
	mealy = "mealy"
	moore = "moore"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	adapter := infrastructure.NewInputOutputAdapter()
	service := app.NewMinimizerService(adapter)

	switch args.Mode {
	case mealy:
		err = service.MinimizeMealy(args.InputFilename, args.OutputFilename)
	case moore:
		err = service.MinimizeMoore(args.InputFilename, args.OutputFilename)
	}

	if err != nil {
		log.Fatal(err)
	}
}
