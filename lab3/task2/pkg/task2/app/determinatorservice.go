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
	automaton, err := s.inputOutputAdapter.GetWithEmpty(inputFilename)
	if err != nil {
		return err
	}

	// TODO: implement determination with empty symbols

	return s.inputOutputAdapter.WriteWithEmpty(outputFilename, automaton)
}
