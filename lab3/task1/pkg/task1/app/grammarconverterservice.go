package app

import (
	"fmt"

	"automata/common/app"
)

const startStateFromLeft = "H"

func NewGrammarConverterService(inputOutputAdapter app.GrammarInputOutputAdapter) *GrammarConverterService {
	return &GrammarConverterService{
		inputOutputAdapter: inputOutputAdapter,
	}
}

type GrammarConverterService struct {
	inputOutputAdapter app.GrammarInputOutputAdapter
}

func (s *GrammarConverterService) ConvertLeftSideGrammarToAutomaton(inputFilename, outputFilename string) error {
	grammar, err := s.inputOutputAdapter.GetGrammar(inputFilename, app.GrammarSideLeft)
	if err != nil {
		return err
	}

	fmt.Println(grammar)

	rightSideGrammar := leftSideToRightSideGrammar(grammar)
	automaton := rightSideGrammarToAutomaton(rightSideGrammar)

	fmt.Println(rightSideGrammar)

	return s.inputOutputAdapter.WriteFinite(outputFilename, automaton)
}

func (s *GrammarConverterService) ConvertRightSideGrammarToAutomaton(inputFilename, outputFilename string) error {
	grammar, err := s.inputOutputAdapter.GetGrammar(inputFilename, app.GrammarSideRight)
	if err != nil {
		return err
	}

	automaton := rightSideGrammarToAutomaton(grammar)

	return s.inputOutputAdapter.WriteFinite(outputFilename, automaton)
}

func leftSideToRightSideGrammar(grammar app.Grammar) app.Grammar {
	newNonTerminals := make([]string, 0, len(grammar.NonTerminalSymbols))
	newNonTerminals = append(newNonTerminals, startStateFromLeft)
	for _, nonTerminal := range grammar.NonTerminalSymbols {
		newNonTerminals = append(newNonTerminals, nonTerminal)
	}

	newRules := make(app.Rules)
	for sourceNonTerminalWithTerminal, destinationNonTerminals := range grammar.Rules {
		for _, destinationNonTerminal := range destinationNonTerminals {
			nonTerminal := destinationNonTerminal
			if len(nonTerminal) == 0 {
				nonTerminal = startStateFromLeft
			}

			key := app.NonTerminalWithTerminal{
				NonTerminalSymbol: nonTerminal,
				TerminalSymbol:    sourceNonTerminalWithTerminal.TerminalSymbol,
			}

			newRules[key] = append(newRules[key], sourceNonTerminalWithTerminal.NonTerminalSymbol)
		}
	}

	return app.Grammar{
		NonTerminalSymbols: newNonTerminals,
		TerminalSymbols:    grammar.TerminalSymbols,
		Rules:              newRules,
		Side:               app.GrammarSideRight,
	}
}

func rightSideGrammarToAutomaton(grammar app.Grammar) app.FiniteAutomaton {
	// TODO
	return app.FiniteAutomaton{}
}
