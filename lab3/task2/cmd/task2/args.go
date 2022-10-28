package main

import (
	"errors"
	"os"
)

func parseArgs() (*args, error) {
	if len(os.Args) != 3 {
		return nil, errors.New("usage: task2 [input csv filename] [output csv filename]")
	}

	return &args{
		InputFilename:  os.Args[1],
		OutputFilename: os.Args[2],
	}, nil
}

type args struct {
	InputFilename  string
	OutputFilename string
}
