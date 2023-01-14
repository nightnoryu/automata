package infrastructure

import (
	"encoding/csv"
	"os"

	"automata/common/app"
)

const (
	csvValuesSeparator = ';'
	emptyMoveIndicator = "-"
)

func readCsv(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = csvValuesSeparator

	return csvReader.ReadAll()
}

func writeCsv(filename string, records [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Comma = csvValuesSeparator

	return csvWriter.WriteAll(records)
}

func getStatesWithSignals(records [][]string) (states []string, signals map[string]string) {
	states = records[1][1:]
	rawSignals := records[0][1:]

	signals = make(map[string]string)
	for i, state := range states {
		signals[state] = rawSignals[i]
	}

	return states, signals
}

func getInputSymbols(records [][]string) []string {
	result := make([]string, 0, len(records)-2)
	for _, row := range records[2:] {
		result = append(result, row[0])
	}

	return result
}

func getDeterministicMoves(
	records [][]string,
	states []string,
	inputSymbols []string,
) app.DeterministicMoves {
	transposedRecords := transpose(records[2:])

	result := make(app.DeterministicMoves)
	for i, stateAndMoves := range transposedRecords[1:] {
		for j, move := range stateAndMoves {
			if move == emptyMoveIndicator {
				continue
			}

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
