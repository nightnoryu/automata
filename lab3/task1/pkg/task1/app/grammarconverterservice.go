package app

import (
	"fmt"
	"sort"
	"strings"

	"automata/common/app"
)

const (
	startStateFromLeft = "H"
	finalStateForRight = "F"
)

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

	return s.inputOutputAdapter.WriteFinite(outputFilename, automaton)
}

func (s *GrammarConverterService) ConvertRightSideGrammarToAutomaton(inputFilename, outputFilename string) error {
	grammar, err := s.inputOutputAdapter.GetGrammar(inputFilename, app.GrammarSideRight)
	if err != nil {
		return err
	}

	fmt.Println("grammar: ", grammar)

	automaton := rightSideGrammarToAutomaton(grammar)

	fmt.Println("automaton: ", automaton)

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
	uniqueStates := make(map[string]bool)
	moves := make(app.MooreMoves)

	var queue []string
	queue = append(queue, grammar.NonTerminalSymbols[0])

	for len(queue) > 0 {
		sourceNonTerminal := queue[0]
		queue = queue[1:]

		if uniqueStates[sourceNonTerminal] {
			continue
		}
		uniqueStates[sourceNonTerminal] = true

		for _, terminalSymbol := range grammar.TerminalSymbols {
			key := app.NonTerminalWithTerminal{
				NonTerminalSymbol: sourceNonTerminal,
				TerminalSymbol:    terminalSymbol,
			}

			destinationNonTerminals, ok := grammar.Rules[key]
			if !ok {
				continue
			}

			newState := finalStateForRight
			combinedDestinationNonTerminals := combineSymbols(destinationNonTerminals)
			if len(destinationNonTerminals) > 0 {
				newState = strings.Join(combinedDestinationNonTerminals, "")
			}

			movesKey := app.InitialStateAndInputSymbol{
				State:  sourceNonTerminal,
				Symbol: terminalSymbol,
			}
			moves[movesKey] = newState

			queue = append(queue, newState)
		}
	}

	return app.FiniteAutomaton{
		States:       uniqueStatesToFinalStates(uniqueStates),
		InputSymbols: grammar.TerminalSymbols,
		Moves:        moves,
	}
}

func uniqueStatesToFinalStates(uniqueStates map[string]bool) []app.StateWithFinalIndication {
	result := make([]app.StateWithFinalIndication, 0, len(uniqueStates))
	for state := range uniqueStates {
		result = append(result, app.StateWithFinalIndication{
			State:   state,
			IsFinal: false, // TODO: fill with actual final states
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].State < result[j].State
	})

	return result
}

func combineSymbols(symbols []string) []string {
	symbolsMap := make(map[string]bool)
	for _, terminal := range symbols {
		symbolsMap[terminal] = true
	}

	var result []string
	for symbol := range symbolsMap {
		resultingSymbol := finalStateForRight
		if len(symbol) > 0 {
			resultingSymbol = symbol
		}

		result = append(result, resultingSymbol)
	}

	sort.Strings(result)

	return result
}
