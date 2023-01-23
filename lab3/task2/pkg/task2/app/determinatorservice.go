package app

import (
	"automata/common/app"
)

func NewDeterminatorService(inputOutputAdapter app.FiniteInputOutputAdapter, determinator *Determinator) *DeterminatorService {
	return &DeterminatorService{
		inputOutputAdapter: inputOutputAdapter,
		determinator:       determinator,
	}
}

type DeterminatorService struct {
	inputOutputAdapter app.FiniteInputOutputAdapter
	determinator       *Determinator
}

func (s *DeterminatorService) Determinate(inputFilename, outputFilename string) error {
	automaton, err := s.inputOutputAdapter.GetNonDeterministicFinite(inputFilename)
	if err != nil {
		return err
	}

	result := s.determinator.Determinate(automaton)

	return s.inputOutputAdapter.WriteFinite(outputFilename, result)
}
