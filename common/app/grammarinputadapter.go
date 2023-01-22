package app

type GrammarInputAdapter interface {
	GetGrammar(filename string, side GrammarSide) (Grammar, error)
}

type Grammar struct {
	NonTerminalSymbols []string
	TerminalSymbols    []string
	Rules              Rules
	Side               GrammarSide
}

type Rules = map[NonTerminalWithTerminal][]string

type NonTerminalWithTerminal struct {
	NonTerminal string
	Terminal    string
}

type GrammarSide int

const (
	GrammarSideLeft = GrammarSide(iota)
	GrammarSideRight
)
