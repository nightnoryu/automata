package app

import (
	"fmt"

	"automata/common/app"
)

func NewGrammarConverterService(
	inputOutputAdapter app.GrammarInputOutputAdapter,
	determinatorService *DeterminatorService,
) *GrammarConverterService {
	return &GrammarConverterService{
		inputOutputAdapter:  inputOutputAdapter,
		determinatorService: determinatorService,
	}
}

type GrammarConverterService struct {
	inputOutputAdapter  app.GrammarInputOutputAdapter
	determinatorService *DeterminatorService
}

func (s *GrammarConverterService) ConvertLeftSideGrammarToAutomaton(inputFilename, outputFilename string) error {
	grammar, err := s.inputOutputAdapter.GetGrammar(inputFilename, app.GrammarSideLeft)
	if err != nil {
		return err
	}

	fmt.Println(grammar)
	// TODO: convert

	return nil
}

func (s *GrammarConverterService) ConvertRightSideGrammarToAutomaton(inputFilename, outputFilename string) error {
	grammar, err := s.inputOutputAdapter.GetGrammar(inputFilename, app.GrammarSideRight)
	if err != nil {
		return err
	}

	fmt.Println(grammar)
	// TODO: convert

	return nil
}
