package app

type MealyMooreInputOutputAdapter interface {
	GetMealy(filename string) (MealyAutomaton, error)
	GetMoore(filename string) (MooreAutomaton, error)

	WriteMealy(filename string, automaton MealyAutomaton) error
	WriteMoore(filename string, automaton MooreAutomaton) error
}

type MealyAutomaton struct {
	States       []string
	InputSymbols []string
	Moves        MovesWithSignals
}

type MooreAutomaton struct {
	States       []string
	InputSymbols []string
	StateSignals map[string]string
	Moves        DeterministicMoves
}
