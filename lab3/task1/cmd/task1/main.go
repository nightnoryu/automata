package main

import (
	"log"

	"automata/lab3/task1/pkg/task1/app"

	"automata/common/infrastructure"
)

const (
	leftSide  = "left"
	rightSide = "right"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	adapter := infrastructure.NewInputOutputAdapter()
	determinatorService := app.NewDeterminatorService()
	service := app.NewGrammarConverterService(adapter, determinatorService)

	switch args.Mode {
	case leftSide:
		err = service.ConvertLeftSideGrammarToAutomaton(args.InputFilename, args.OutputFilename)
	case rightSide:
		err = service.ConvertRightSideGrammarToAutomaton(args.InputFilename, args.OutputFilename)
	}

	if err != nil {
		log.Fatal(err)
	}
}
