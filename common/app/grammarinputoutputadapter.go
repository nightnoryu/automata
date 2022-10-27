package app

type GrammarInputOutputAdapter interface {
	GetGrammar(filename string, side GrammarSide) (Grammar, error)
	GetWithEmpty(filename string) (EmptyMovesAutomaton, error)

	WriteWithEmpty(filename string, automaton EmptyMovesAutomaton) error
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
}

type Rules = map[string]Rule

type Rule struct {
	NonTerminalSymbol string
	TerminalSymbol    string
}

type EmptyMovesAutomaton struct {
	States       []StateWithFinalIndication
	InputSymbols []string
	Moves        MooreMoves
}

type StateWithFinalIndication struct {
	State   string
	IsFinal bool
}
