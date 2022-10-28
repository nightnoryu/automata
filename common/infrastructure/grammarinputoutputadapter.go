package infrastructure

import (
	"bufio"
	"encoding/csv"
	"os"
	"strings"

	"automata/common/app"
)

const (
	nonTerminalRuleSeparator = " -> "
	ruleSeparator            = " | "

	finalStateIndication = "F"
)

func NewGrammarInputOutputAdapter() app.GrammarInputOutputAdapter {
	return &grammarInputOutputAdapter{}
}

type grammarInputOutputAdapter struct{}

func (a *grammarInputOutputAdapter) GetGrammar(filename string, side app.GrammarSide) (app.Grammar, error) {
	file, err := os.Open(filename)
	if err != nil {
		return app.Grammar{}, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	var nonTerminals []string
	uniqueTerminals := make(map[string]bool)
	rules := make(app.Rules)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rule := strings.Split(scanner.Text(), nonTerminalRuleSeparator)

		nonTerminal := rule[0]
		nonTerminals = append(nonTerminals, nonTerminal)

		for _, resultSymbols := range strings.Split(rule[1], ruleSeparator) {
			if len(resultSymbols) == 2 {
				left := resultSymbols[0]
				right := resultSymbols[1]

				if side == app.GrammarSideLeft {
					uniqueTerminals[string(right)] = true
					rules[nonTerminal] = append(rules[nonTerminal], app.Rule{
						NonTerminalSymbol: string(left),
						TerminalSymbol:    string(right),
					})
				} else {
					uniqueTerminals[string(left)] = true
					rules[nonTerminal] = append(rules[nonTerminal], app.Rule{
						NonTerminalSymbol: string(right),
						TerminalSymbol:    string(left),
					})
				}
			} else {
				rules[nonTerminal] = append(rules[nonTerminal], app.Rule{
					TerminalSymbol: resultSymbols,
				})
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return app.Grammar{}, err
	}

	terminals := make([]string, 0, len(uniqueTerminals))
	for terminal := range uniqueTerminals {
		terminals = append(terminals, terminal)
	}

	return app.Grammar{
		NonTerminalSymbols: nonTerminals,
		TerminalSymbols:    terminals,
		Rules:              rules,
		Side:               side,
	}, nil
}

func (a *grammarInputOutputAdapter) GetWithEmpty(filename string) (app.GrammarAutomaton, error) {
	file, err := os.Open(filename)
	if err != nil {
		return app.GrammarAutomaton{}, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = csvValuesSeparator

	records, err := csvReader.ReadAll()
	if err != nil {
		return app.GrammarAutomaton{}, err
	}

	states := getStatesWithFinalIndication(records)
	inputSymbols := getStateSignalsDependentInputSymbols(records)

	return app.GrammarAutomaton{
		States:       states,
		InputSymbols: inputSymbols,
		Moves:        getMooreMoves(records, getPlainStatesFromGrammarStates(states), inputSymbols),
	}, nil
}

func (a *grammarInputOutputAdapter) WriteWithEmpty(filename string, automaton app.GrammarAutomaton) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Comma = csvValuesSeparator

	return csvWriter.WriteAll(serializeGrammarAutomaton(automaton))
}

func serializeGrammarAutomaton(automaton app.GrammarAutomaton) [][]string {
	result := make([][]string, len(automaton.InputSymbols)+2)
	for i := range result {
		result[i] = make([]string, 0, len(automaton.States)+1)
	}

	result[0] = append(result[0], "")
	result[1] = append(result[1], "")
	for _, state := range automaton.States {
		if state.IsFinal {
			result[0] = append(result[0], finalStateIndication)
		} else {
			result[0] = append(result[0], "")
		}
		result[1] = append(result[1], state.State)
	}

	for i, inputSymbol := range automaton.InputSymbols {
		result[i+2] = append(result[i+2], inputSymbol)

		for _, state := range automaton.States {
			key := app.InitialStateAndInputSymbol{
				State:  state.State,
				Symbol: inputSymbol,
			}

			result[i+2] = append(result[i+2], automaton.Moves[key])
		}
	}

	return result
}

func getStatesWithFinalIndication(records [][]string) []app.StateWithFinalIndication {
	states := records[1][1:]
	finalIndicators := records[0][1:]

	result := make([]app.StateWithFinalIndication, 0, len(states))
	for i, state := range states {
		result = append(result, app.StateWithFinalIndication{
			State:   state,
			IsFinal: finalIndicators[i] == grammarFinalStateIndicator,
		})
	}

	return result
}

func getStateSignalsDependentInputSymbols(records [][]string) []string {
	result := make([]string, 0, len(records)-2)
	for _, row := range records[2:] {
		result = append(result, row[0])
	}

	return result
}

func getPlainStatesFromGrammarStates(states []app.StateWithFinalIndication) []string {
	result := make([]string, 0, len(states))
	for _, state := range states {
		result = append(result, state.State)
	}

	return result
}
