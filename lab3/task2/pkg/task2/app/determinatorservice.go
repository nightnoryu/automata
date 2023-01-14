package app

import (
	"fmt"

	"automata/common/app"
)

func NewDeterminatorService(inputOutputAdapter app.FiniteInputOutputAdapter) *DeterminatorService {
	return &DeterminatorService{
		inputOutputAdapter: inputOutputAdapter,
	}
}

type DeterminatorService struct {
	inputOutputAdapter app.FiniteInputOutputAdapter
}

func (s *DeterminatorService) Determinate(inputFilename, outputFilename string) error {
	automaton, err := s.inputOutputAdapter.GetNonDeterministicFinite(inputFilename)
	if err != nil {
		return err
	}

	fmt.Println(automaton)

	return nil
}
