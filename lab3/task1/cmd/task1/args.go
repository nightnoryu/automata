package main

import (
	"automata/common/app"
	"errors"
	"os"
)

func parseArgs() (*args, error) {
	if len(os.Args) != 4 {
		return nil, errors.New("usage: task1 [left|right] [input grammar filename] [output csv filename]")
	}

	grammarSide, err := getGrammarSide(os.Args[1])
	if err != nil {
		return nil, err
	}

	return &args{
		GrammarSide:    grammarSide,
		InputFilename:  os.Args[2],
		OutputFilename: os.Args[3],
	}, nil
}

type args struct {
	GrammarSide    app.GrammarSide
	InputFilename  string
	OutputFilename string
}

func getGrammarSide(side string) (app.GrammarSide, error) {
	switch os.Args[1] {
	case "left":
		return app.GrammarSideLeft, nil
	case "right":
		return app.GrammarSideRight, nil
	default:
		return 0, errors.New("invalid grammar side")
	}
}
