package app

type InputOutputAdapter interface {
	GetMealy(filename string) (MealyAutomaton, error)
	GetMoore(filename string) (MooreAutomaton, error)

	WriteMealy(filename string, automaton MealyAutomaton) error
	WriteMoore(filename string, automaton MooreAutomaton) error
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

type InitialStateAndInputSymbol struct {
	State  string
	Symbol string
}

type MealyMoves map[InitialStateAndInputSymbol]DestinationStateAndSignal

type MooreMoves = map[InitialStateAndInputSymbol]string

type DestinationStateAndSignal struct {
	State  string
	Signal string
}
