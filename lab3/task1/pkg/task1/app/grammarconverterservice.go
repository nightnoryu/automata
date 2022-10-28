package app

import "automata/common/app"

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
	// TODO: model and conversion + determination
	return nil
}

func (s *GrammarConverterService) ConvertRightSideGrammarToAutomaton(inputFilename, outputFilename string) error {
	// TODO: model and conversion + determination
	return nil
}
