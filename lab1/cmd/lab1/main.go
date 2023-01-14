package main

import (
	"log"

	"automata/common/infrastructure"
	"automata/lab1/pkg/lab1/app"
)

const (
	mealyToMoore = "mealy-to-moore"
	mooreToMealy = "moore-to-mealy"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	adapter := infrastructure.NewMealyMooreInputOutputAdapter()
	service := app.NewTranslatorService(adapter)

	switch args.Mode {
	case mealyToMoore:
		err = service.MealyToMoore(args.InputFilename, args.OutputFilename)
	case mooreToMealy:
		err = service.MooreToMealy(args.InputFilename, args.OutputFilename)
	}

	if err != nil {
		log.Fatal(err)
	}
}
