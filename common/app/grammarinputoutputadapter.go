package app

type GrammarInputOutputAdapter interface {
	GetGrammar(filename string, side GrammarSide) (Grammar, error)
	GetWithEmpty(filename string) (GrammarAutomaton, error)

	WriteWithEmpty(filename string, automaton GrammarAutomaton) error
}

type GrammarSide int

const (
	GrammarSideLeft = GrammarSide(iota)
	GrammarSideRight
)

type Grammar struct {
	NonTerminalSymbols []string
	TerminalSymbols    []string
	Rules              Rules
	Side               GrammarSide
}

type Rules = map[string][]Rule

type Rule struct {
	NonTerminalSymbol string
	TerminalSymbol    string
}

type GrammarAutomaton struct {
	States       []StateWithFinalIndication
	InputSymbols []string
	Moves        MooreMoves
}

type StateWithFinalIndication struct {
	State   string
	IsFinal bool
}
