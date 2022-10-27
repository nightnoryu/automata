package app

type InputOutputAdapter interface {
	GetMealy(filename string) (MealyAutomaton, error)
	GetMoore(filename string) (MooreAutomaton, error)
	GetWithEmpty(filename string) (EmptyMovesAutomaton, error)

	WriteMealy(filename string, automaton MealyAutomaton) error
	WriteMoore(filename string, automaton MooreAutomaton) error
	WriteWithEmpty(filename string, automaton EmptyMovesAutomaton) error
}

type MealyAutomaton struct {
	States       []string
	InputSymbols []string
	Moves        MealyMoves
}

type MooreAutomaton struct {
	States       []string
	InputSymbols []string
	StateSignals map[string]string
	Moves        MooreMoves
}

type EmptyMovesAutomaton struct {
	States       []StateWithFinalIndication
	InputSymbols []string
	Moves        MooreMoves
}

type MealyMoves map[InitialStateAndInputSymbol]DestinationStateAndSignal

type MooreMoves = map[InitialStateAndInputSymbol]string

type InitialStateAndInputSymbol struct {
	State  string
	Symbol string
}

type DestinationStateAndSignal struct {
	State  string
	Signal string
}

type StateWithFinalIndication struct {
	State   string
	IsFinal bool
}
