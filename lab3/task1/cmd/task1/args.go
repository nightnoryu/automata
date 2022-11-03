package main

import (
	"errors"
	"os"
)

func parseArgs() (*args, error) {
	if len(os.Args) != 4 {
		return nil, errors.New("usage: task1 [left|right] [input filename] [output csv filename]")
	}

	return &args{
		Mode:           os.Args[1],
		InputFilename:  os.Args[2],
		OutputFilename: os.Args[3],
	}, nil
}

type args struct {
	Mode           string
	InputFilename  string
	OutputFilename string
}
