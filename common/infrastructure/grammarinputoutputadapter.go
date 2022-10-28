package infrastructure

import (
	"bufio"
	"encoding/csv"
	"os"
	"strings"

	"automata/common/app"
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
		rule := strings.Split(scanner.Text(), " -> ")
		nonTerminals = append(nonTerminals, rule[0])

		for _, resultSymbols := range strings.Split(rule[1], " | ") {
			if len(resultSymbols) == 2 {
				if side == app.GrammarSideLeft {
					uniqueTerminals[string(resultSymbols[1])] = true
					rules[rule[0]] = append(rules[rule[0]], app.Rule{
						NonTerminalSymbol: string(resultSymbols[0]),
						TerminalSymbol:    string(resultSymbols[1]),
					})
				} else {
					uniqueTerminals[string(resultSymbols[0])] = true
					rules[rule[0]] = append(rules[rule[0]], app.Rule{
						NonTerminalSymbol: string(resultSymbols[1]),
						TerminalSymbol:    string(resultSymbols[0]),
					})
				}
			} else {
				rules[rule[0]] = append(rules[rule[0]], app.Rule{
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

	// TODO

	return app.GrammarAutomaton{}, nil
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

	return csvWriter.WriteAll(serializeGrammar(automaton))
}

func serializeGrammar(automaton app.GrammarAutomaton) [][]string {
	// TODO
	return nil
}
