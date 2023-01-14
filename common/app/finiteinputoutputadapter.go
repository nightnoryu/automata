package app

type FiniteInputOutputAdapter interface {
	GetNonDeterministicFinite(filename string) (NonDeterministicFiniteAutomaton, error)

	WriteFinite(filename string, automaton FiniteAutomaton) error
}

type FiniteAutomaton struct {
	States       []string
	InputSymbols []string
	FinalStates  map[string]bool
	Moves        DeterministicMoves
}

type NonDeterministicFiniteAutomaton struct {
	States       []string
	InputSymbols []string
	FinalStates  map[string]bool
	Moves        NonDeterministicMoves
}
