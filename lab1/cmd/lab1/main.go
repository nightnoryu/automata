package main

import (
	"automata/lab1/pkg/lab1/app"
	"automata/lab1/pkg/lab1/infrastructure"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	adapter := infrastructure.NewInputOutputAdapter()
	service := app.NewTranslatorService(adapter)

	application := &cli.App{
		Name:  "lab1",
		Usage: "Mealy to Moore and vice-versa automaton converter",
		Commands: []*cli.Command{
			{
				Name: "mealy-to-moore",
				Action: func(ctx *cli.Context) error {
					return service.MealyToMoore(ctx.String("input"), ctx.String("output"))
				},
			},
			{
				Name: "moore-to-mealy",
				Action: func(ctx *cli.Context) error {
					return service.MooreToMealy(ctx.String("input"), ctx.String("output"))
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Usage:    "input CSV filename",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Usage:    "output CSV filename",
				Required: true,
			},
		},
	}

	if err := application.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
