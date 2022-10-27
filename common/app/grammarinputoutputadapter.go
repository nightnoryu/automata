package app

type GrammarInputOutputAdapter interface {
	GetGrammar(filename string) (Grammar, error)
	GetWithEmpty(filename string) (EmptyMovesAutomaton, error)

	WriteWithEmpty(filename string, automaton EmptyMovesAutomaton) error
}

type Grammar struct {
	// TODO: describe model
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
