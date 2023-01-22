package infrastructure

import (
	"bufio"
	"os"
	"sort"
	"strings"

	"automata/common/app"
)

const (
	nonTerminalRuleSeparator = " -> "
	ruleSeparator            = " | "
)

func NewGrammarInputAdapter() app.GrammarInputAdapter {
	return &grammarInputAdapter{}
}

type grammarInputAdapter struct{}

func (a *grammarInputAdapter) GetGrammar(filename string, side app.GrammarSide) (app.Grammar, error) {
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

		sourceNonTerminal := rule[0]
		nonTerminals = append(nonTerminals, sourceNonTerminal)

		for _, resultSymbols := range strings.Split(rule[1], ruleSeparator) {
			var destinationNonTerminal, terminal string

			if len(resultSymbols) == 1 {
				terminal = resultSymbols
			} else if side == app.GrammarSideLeft {
				destinationNonTerminal = string(resultSymbols[0])
				terminal = string(resultSymbols[1])
			} else if side == app.GrammarSideRight {
				destinationNonTerminal = string(resultSymbols[1])
				terminal = string(resultSymbols[0])
			}

			uniqueTerminals[terminal] = true

			key := app.NonTerminalWithTerminal{
				NonTerminal: sourceNonTerminal,
				Terminal:    terminal,
			}

			rules[key] = append(rules[key], destinationNonTerminal)
		}
	}

	if err = scanner.Err(); err != nil {
		return app.Grammar{}, err
	}

	return app.Grammar{
		NonTerminalSymbols: nonTerminals,
		TerminalSymbols:    uniqueTerminalsToFinalTerminals(uniqueTerminals),
		Rules:              rules,
		Side:               side,
	}, nil
}

func uniqueTerminalsToFinalTerminals(uniqueTerminals map[string]bool) []string {
	result := make([]string, 0, len(uniqueTerminals))
	for terminal := range uniqueTerminals {
		result = append(result, terminal)
	}

	sort.Strings(result)

	return result
}
