package infrastructure

import (
	"strings"

	"automata/common/app"
)

const (
	stateAndSignalSeparator = "/"
)

func NewMealyMooreInputOutputAdapter() app.MealyMooreInputOutputAdapter {
	return &mealyMooreInputOutputAdapter{}
}

type mealyMooreInputOutputAdapter struct{}

func (a *mealyMooreInputOutputAdapter) GetMealy(filename string) (app.MealyAutomaton, error) {
	records, err := readCsv(filename)
	if err != nil {
		return app.MealyAutomaton{}, err
	}

	states := getMealyStates(records)
	inputSymbols := getMealyInputSymbols(records)

	return app.MealyAutomaton{
		States:       states,
		InputSymbols: inputSymbols,
		Moves:        getMovesWithSignals(records, states, inputSymbols),
	}, nil
}

func (a *mealyMooreInputOutputAdapter) GetMoore(filename string) (app.MooreAutomaton, error) {
	records, err := readCsv(filename)
	if err != nil {
		return app.MooreAutomaton{}, err
	}

	states, stateSignals := getStatesWithSignals(records)
	inputSymbols := getInputSymbols(records)

	return app.MooreAutomaton{
		States:       states,
		InputSymbols: inputSymbols,
		StateSignals: stateSignals,
		Moves:        getDeterministicMoves(records, states, inputSymbols),
	}, nil
}

func (a *mealyMooreInputOutputAdapter) WriteMealy(filename string, automaton app.MealyAutomaton) error {
	return writeCsv(filename, serializeMealy(automaton))
}

func (a *mealyMooreInputOutputAdapter) WriteMoore(filename string, automaton app.MooreAutomaton) error {
	return writeCsv(filename, serializeMoore(automaton))
}

func getMealyStates(records [][]string) []string {
	return records[0][1:]
}

func getMealyInputSymbols(records [][]string) []string {
	result := make([]string, 0, len(records)-1)
	for _, row := range records[1:] {
		result = append(result, row[0])
	}

	return result
}

func getMovesWithSignals(
	records [][]string,
	states, inputSymbols []string,
) app.MovesWithSignals {
	transposedRecords := transpose(records[1:])

	result := make(map[app.InitialStateAndInputSymbol]app.DestinationStateAndSignal)
	for i, stateAndMoves := range transposedRecords[1:] {
		for j, move := range stateAndMoves {
			stateAndInput := app.InitialStateAndInputSymbol{
				State:  states[i],
				Symbol: inputSymbols[j],
			}

			split := strings.Split(move, stateAndSignalSeparator)

			result[stateAndInput] = app.DestinationStateAndSignal{
				State:  split[0],
				Signal: split[1],
			}
		}
	}

	return result
}

func serializeMealy(automaton app.MealyAutomaton) [][]string {
	result := make([][]string, len(automaton.InputSymbols)+1)
	for i := range result {
		result[i] = make([]string, 0, len(automaton.States)+1)
	}

	result[0] = append(result[0], "")
	for _, state := range automaton.States {
		result[0] = append(result[0], state)
	}

	for i, inputSymbol := range automaton.InputSymbols {
		result[i+1] = append(result[i+1], inputSymbol)

		for _, state := range automaton.States {
			key := app.InitialStateAndInputSymbol{
				State:  state,
				Symbol: inputSymbol,
			}

			result[i+1] = append(result[i+1], serializeMealyMove(automaton.Moves[key]))
		}
	}

	return result
}

func serializeMoore(automaton app.MooreAutomaton) [][]string {
	result := make([][]string, len(automaton.InputSymbols)+2)
	for i := range result {
		result[i] = make([]string, 0, len(automaton.States)+1)
	}

	result[0] = append(result[0], "")
	result[1] = append(result[1], "")
	for _, state := range automaton.States {
		result[0] = append(result[0], automaton.StateSignals[state])
		result[1] = append(result[1], state)
	}

	for i, inputSymbol := range automaton.InputSymbols {
		result[i+2] = append(result[i+2], inputSymbol)

		for _, state := range automaton.States {
			key := app.InitialStateAndInputSymbol{
				State:  state,
				Symbol: inputSymbol,
			}

			result[i+2] = append(result[i+2], automaton.Moves[key])
		}
	}

	return result
}

func serializeMealyMove(stateAndSignal app.DestinationStateAndSignal) string {
	return stateAndSignal.State + stateAndSignalSeparator + stateAndSignal.Signal
}
