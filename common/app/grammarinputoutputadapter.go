package app

type GrammarInputOutputAdapter interface {
	GetGrammar(filename string, side GrammarSide) (Grammar, error)
	GetFinite(filename string) (FiniteAutomaton, error)

	WriteFinite(filename string, automaton FiniteAutomaton) error
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

type Rules = map[string][]NonTerminalWithTerminal

type NonTerminalWithTerminal struct {
	NonTerminalSymbol string
	TerminalSymbol    string
}

type FiniteAutomaton struct {
	States       []StateWithFinalIndication
	InputSymbols []string
	Moves        MooreMoves
}

type NonDeterministicFiniteAutomaton struct {
	States       []StateWithFinalIndication
	InputSymbols []string
	Moves        NonDeterministicMoves
}

type NonDeterministicMoves = map[InitialStateAndInputSymbol][]string

type StateWithFinalIndication struct {
	State   string
	IsFinal bool
}
