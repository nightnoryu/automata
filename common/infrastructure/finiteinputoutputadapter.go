package infrastructure

import (
	"automata/common/app"
	"strings"
)

const (
	finalStateIndicator            = "F"
	nonDeterministicMovesSeparator = ","
)

func NewFiniteInputOutputAdapter() app.FiniteInputOutputAdapter {
	return &finiteInputOutputAdapter{}
}

type finiteInputOutputAdapter struct{}

func (a *finiteInputOutputAdapter) GetNonDeterministicFinite(
	filename string,
) (app.NonDeterministicFiniteAutomaton, error) {
	records, err := readCsv(filename)
	if err != nil {
		return app.NonDeterministicFiniteAutomaton{}, err
	}

	states, stateSignals := getStatesWithSignals(records)
	inputSymbols := getInputSymbols(records)

	return app.NonDeterministicFiniteAutomaton{
		States:       states,
		InputSymbols: inputSymbols,
		FinalStates:  buildFinalStates(stateSignals),
		Moves:        getNonDeterministicMoves(records, states, inputSymbols),
	}, nil
}

func (a *finiteInputOutputAdapter) WriteFinite(filename string, automaton app.FiniteAutomaton) error {
	return writeCsv(filename, serializeFinite(automaton))
}

func buildFinalStates(stateSignals map[string]string) map[string]bool {
	result := make(map[string]bool)
	for state, signal := range stateSignals {
		if signal == finalStateIndicator {
			result[state] = true
		}
	}

	return result
}

func getNonDeterministicMoves(records [][]string, states []string, inputSymbols []string) app.NonDeterministicMoves {
	transposedRecords := transpose(records[2:])

	result := make(app.NonDeterministicMoves)
	for i, stateAndMoves := range transposedRecords[1:] {
		for j, moves := range stateAndMoves {
			if moves == emptyMoveIndicator {
				continue
			}

			stateAndInput := app.InitialStateAndInputSymbol{
				State:  states[i],
				Symbol: inputSymbols[j],
			}

			for _, move := range strings.Split(moves, nonDeterministicMovesSeparator) {
				result[stateAndInput] = append(result[stateAndInput], move)
			}
		}
	}

	return result
}

func serializeFinite(automaton app.FiniteAutomaton) [][]string {
	result := make([][]string, len(automaton.InputSymbols)+2)
	for i := range result {
		result[i] = make([]string, 0, len(automaton.States)+1)
	}

	result[0] = append(result[0], "")
	result[1] = append(result[1], "")
	for _, state := range automaton.States {
		if automaton.FinalStates[state] {
			result[0] = append(result[0], finalStateIndicator)
		} else {
			result[0] = append(result[0], "")
		}
		result[1] = append(result[1], state)
	}

	for i, inputSymbol := range automaton.InputSymbols {
		result[i+2] = append(result[i+2], inputSymbol)

		for _, state := range automaton.States {
			key := app.InitialStateAndInputSymbol{
				State:  state,
				Symbol: inputSymbol,
			}

			destination := automaton.Moves[key]
			if destination == "" {
				destination = emptyMoveIndicator
			}

			result[i+2] = append(result[i+2], destination)
		}
	}

	return result
}
