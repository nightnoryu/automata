package app

import (
	"fmt"
	"strings"

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

	rightSideGrammar := leftSideToRightSideGrammar(grammar)
	automaton := rightSideGrammarToAutomaton(rightSideGrammar)

	fmt.Println("left: ", grammar)
	fmt.Println("right: ", rightSideGrammar)

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
	invertedRules := make(map[app.NonTerminalWithTerminal][]string)
	for nonTerminal, rules := range grammar.Rules {
		for _, rule := range rules {
			if len(rule.NonTerminalSymbol) != 0 {
				invertedRules[rule] = append(invertedRules[rule], nonTerminal)
			}
		}
	}

	newNonTerminals := make([]string, len(grammar.NonTerminalSymbols))
	copy(newNonTerminals, grammar.NonTerminalSymbols)
	newNonTerminals = append(newNonTerminals, startStateFromLeft)

	newRules := make(app.Rules)
	for sourceNonTerminal, rules := range grammar.Rules {
		for _, rule := range rules {
			key := app.NonTerminalWithTerminal{
				NonTerminalSymbol: sourceNonTerminal,
				TerminalSymbol:    rule.TerminalSymbol,
			}

			nonTerminal := rule.NonTerminalSymbol
			if len(nonTerminal) == 0 {
				nonTerminal = startStateFromLeft
			}

			newRules[nonTerminal] = append(newRules[nonTerminal], app.NonTerminalWithTerminal{
				NonTerminalSymbol: strings.Join(invertedRules[key], ""),
				TerminalSymbol:    rule.TerminalSymbol,
			})
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
