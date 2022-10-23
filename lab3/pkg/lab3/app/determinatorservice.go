package app

import "automata/common/app"

const newStatesIdentifier = "S"

func NewDeterminatorService(inputOutputAdapter app.InputOutputAdapter) *DeterminatorService {
	return &DeterminatorService{
		inputOutputAdapter: inputOutputAdapter,
	}
}

type DeterminatorService struct {
	inputOutputAdapter app.InputOutputAdapter
}

func (s *DeterminatorService) Determinate(inputFilename, outputFilename string) error {
	// TODO: implement
	return nil
}
