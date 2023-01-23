package app

import (
	"errors"
	"fmt"

	"automata/common/app"
	task2app "automata/lab3/task2/pkg/task2/app"
)

const (
	startingState = "H"
	finalState    = "F"
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
	fmt.Println(automaton)
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
	states := make([]string, 0, len(grammar.NonTerminalSymbols)+1)
	for _, symbol := range grammar.NonTerminalSymbols {
		states = append(states, symbol)
	}
	states = append(states, finalState)

	finalStates := make(map[string]bool)
	finalStates[finalState] = true

	moves := make(app.NonDeterministicMoves)
	for nonTerminalWithTerminal, destinationNonTerminals := range grammar.Rules {
		key := app.InitialStateAndInputSymbol{
			State:  nonTerminalWithTerminal.NonTerminal,
			Symbol: nonTerminalWithTerminal.Terminal,
		}

		for _, destinationNonTerminal := range destinationNonTerminals {
			destination := destinationNonTerminal
			if destination == "" {
				destination = finalState
			}

			moves[key] = append(moves[key], destination)
		}
	}

	return app.NonDeterministicFiniteAutomaton{
		States:       states,
		InputSymbols: grammar.TerminalSymbols,
		FinalStates:  finalStates,
		Moves:        moves,
	}
}

func leftSideGrammarToAutomaton(grammar app.Grammar) app.NonDeterministicFiniteAutomaton {
	states := make([]string, 0, len(grammar.NonTerminalSymbols)+1)
	states = append(states, startingState)
	for _, symbol := range grammar.NonTerminalSymbols {
		states = append(states, symbol)
	}

	finalStates := make(map[string]bool)
	finalStates[states[1]] = true

	moves := make(app.NonDeterministicMoves)
	for nonTerminalWithTerminal, destinationNonTerminals := range grammar.Rules {
		for _, destinationNonTerminal := range destinationNonTerminals {
			initialState := destinationNonTerminal
			if initialState == "" {
				initialState = startingState
			}

			key := app.InitialStateAndInputSymbol{
				State:  initialState,
				Symbol: nonTerminalWithTerminal.Terminal,
			}

			moves[key] = append(moves[key], nonTerminalWithTerminal.NonTerminal)
		}
	}

	return app.NonDeterministicFiniteAutomaton{
		States:       states,
		InputSymbols: grammar.TerminalSymbols,
		FinalStates:  finalStates,
		Moves:        moves,
	}
}
