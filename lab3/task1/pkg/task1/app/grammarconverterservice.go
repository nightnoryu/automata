package app

import (
	"automata/common/app"
	task2app "automata/lab3/task2/pkg/task2/app"
	"errors"
)

func NewGrammarConverterService(
	grammarInputAdapter app.GrammarInputAdapter,
	finiteInputOutputAdapter app.FiniteInputOutputAdapter,
	determinator *task2app.Determinator,
) *GrammarConverterService {
	return &GrammarConverterService{
		grammarInputAdapter:      grammarInputAdapter,
		finiteInputOutputAdapter: finiteInputOutputAdapter,
		determinator:             determinator,
	}
}

type GrammarConverterService struct {
	grammarInputAdapter      app.GrammarInputAdapter
	finiteInputOutputAdapter app.FiniteInputOutputAdapter
	determinator             *task2app.Determinator
}

func (s *GrammarConverterService) ConvertToFinite(
	grammarSide app.GrammarSide,
	inputFilename, outputFilename string,
) error {
	grammar, err := s.grammarInputAdapter.GetGrammar(inputFilename, grammarSide)
	if err != nil {
		return err
	}

	automaton, err := grammarToAutomaton(grammar)
	if err != nil {
		return err
	}
	result := s.determinator.Determinate(automaton)

	return s.finiteInputOutputAdapter.WriteFinite(outputFilename, result)
}

func grammarToAutomaton(grammar app.Grammar) (app.NonDeterministicFiniteAutomaton, error) {
	switch grammar.Side {
	case app.GrammarSideLeft:
		return leftSideGrammarToAutomaton(grammar), nil
	case app.GrammarSideRight:
		return rightSideGrammarToAutomaton(grammar), nil
	default:
		return app.NonDeterministicFiniteAutomaton{}, errors.New("unsupported grammar type")
	}
}

func rightSideGrammarToAutomaton(grammar app.Grammar) app.NonDeterministicFiniteAutomaton {
	return app.NonDeterministicFiniteAutomaton{}
}

func leftSideGrammarToAutomaton(grammar app.Grammar) app.NonDeterministicFiniteAutomaton {
	return app.NonDeterministicFiniteAutomaton{}
}
