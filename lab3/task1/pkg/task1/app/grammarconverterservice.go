package app

import (
	"automata/common/app"
	"fmt"
)

func NewGrammarConverterService(
	grammarInputAdapter app.GrammarInputAdapter,
	finiteInputOutputAdapter app.FiniteInputOutputAdapter,
) *GrammarConverterService {
	return &GrammarConverterService{
		grammarInputAdapter:      grammarInputAdapter,
		finiteInputOutputAdapter: finiteInputOutputAdapter,
	}
}

type GrammarConverterService struct {
	grammarInputAdapter      app.GrammarInputAdapter
	finiteInputOutputAdapter app.FiniteInputOutputAdapter
}

func (s *GrammarConverterService) ConvertToFinite(
	grammarSide app.GrammarSide,
	inputFilename, outputFilename string,
) error {
	grammar, err := s.grammarInputAdapter.GetGrammar(inputFilename, grammarSide)
	if err != nil {
		return err
	}

	fmt.Println(grammar)

	return nil
}
