package main

import (
	"log"

	"automata/common/infrastructure"
	"automata/lab3/task1/pkg/task1/app"
	task2app "automata/lab3/task2/pkg/task2/app"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	grammarInputAdapter := infrastructure.NewGrammarInputAdapter()
	finiteInputOutputAdapter := infrastructure.NewFiniteInputOutputAdapter()
	determinator := task2app.NewDeterminator()
	service := app.NewGrammarConverterService(grammarInputAdapter, finiteInputOutputAdapter, determinator)

	if err = service.ConvertToFinite(args.GrammarSide, args.InputFilename, args.OutputFilename); err != nil {
		log.Fatal(err)
	}
}
