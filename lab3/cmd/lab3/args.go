package main

import (
	"errors"
	"os"
)

func parseArgs() (*args, error) {
	if len(os.Args) != 3 {
		return nil, errors.New("usage: lab3 [input csv filename] [output csv filename]")
	}

	return &args{
		InputFilename:  os.Args[2],
		OutputFilename: os.Args[3],
	}, nil
}

type args struct {
	InputFilename  string
	OutputFilename string
}
