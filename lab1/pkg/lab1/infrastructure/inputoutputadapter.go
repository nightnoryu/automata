package infrastructure

import (
	"encoding/csv"
	"os"
	"strings"

	"automata/lab1/pkg/lab1/app"
)

const stateAndSignalSeparator = "/"

func NewInputOutputAdapter() app.InputOutputAdapter {
	return &inputOutputAdapter{}
}

type inputOutputAdapter struct{}

func (a *inputOutputAdapter) GetMealy(filename string) (app.MealyAutomaton, error) {
	file, err := os.Open(filename)
	if err != nil {
		return app.MealyAutomaton{}, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	records, err := csvReader.ReadAll()
	if err != nil {
		return app.MealyAutomaton{}, err
	}

	states := getMealyStates(records)
	inputSymbols := getMealyInputSymbols(records)

	return app.MealyAutomaton{
		States:       states,
		InputSymbols: inputSymbols,
		Moves:        getMealyMoves(records, states, inputSymbols),
	}, nil
}

func (a *inputOutputAdapter) GetMoore(filename string) (app.MooreAutomaton, error) {
	file, err := os.Open(filename)
	if err != nil {
		return app.MooreAutomaton{}, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	records, err := csvReader.ReadAll()
	if err != nil {
		return app.MooreAutomaton{}, err
	}

	states := getMooreStates(records)
	inputSymbols := getMooreInputSymbols(records)
	stateSignals := getMooreStateSignals(records)

	return app.MooreAutomaton{
		States:       states,
		InputSymbols: inputSymbols,
		StateSignals: stateSignals,
		Moves:        getMooreMoves(records, states, inputSymbols),
	}, nil
}

func (a *inputOutputAdapter) WriteMealy(filename string, automaton app.MealyAutomaton) error {
	// TODO
	return nil
}

func (a *inputOutputAdapter) WriteMoore(filename string, automaton app.MooreAutomaton) error {
	// TODO
	return nil
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

func getMealyMoves(
	records [][]string,
	states, inputSymbols []string,
) map[app.InitialStateAndInputSymbol]app.DestinationStateAndSignal {
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

func getMooreStates(records [][]string) []string {
	return records[1][1:]
}

func getMooreInputSymbols(records [][]string) []string {
	result := make([]string, 0, len(records)-2)
	for _, row := range records[2:] {
		result = append(result, row[0])
	}

	return result
}

func getMooreStateSignals(records [][]string) map[string]string {
	states := getMooreStates(records)
	signals := records[0][1:]

	result := make(map[string]string)
	for i, state := range states {
		result[state] = signals[i]
	}

	return result
}

func getMooreMoves(
	records [][]string,
	states []string,
	inputSymbols []string,
) map[app.InitialStateAndInputSymbol]string {
	transposedRecords := transpose(records[2:])

	result := make(map[app.InitialStateAndInputSymbol]string)
	for i, stateAndMoves := range transposedRecords[1:] {
		for j, move := range stateAndMoves {
			stateAndInput := app.InitialStateAndInputSymbol{
				State:  states[i],
				Symbol: inputSymbols[j],
			}

			result[stateAndInput] = move
		}
	}

	return result
}

func transpose(matrix [][]string) [][]string {
	xl := len(matrix[0])
	yl := len(matrix)

	result := make([][]string, xl)
	for i := range result {
		result[i] = make([]string, yl)
	}

	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = matrix[j][i]
		}
	}

	return result
}
