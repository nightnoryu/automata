package main

import (
	"log"

	"automata/common/infrastructure"
	"automata/lab3/task1/pkg/task1/app"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	grammarInputAdapter := infrastructure.NewGrammarInputAdapter()
	finiteInputOutputAdapter := infrastructure.NewFiniteInputOutputAdapter()
	service := app.NewGrammarConverterService(grammarInputAdapter, finiteInputOutputAdapter)

	if err = service.ConvertToFinite(args.GrammarSide, args.InputFilename, args.OutputFilename); err != nil {
		log.Fatal(err)
	}
}
