package app

type InputOutputAdapter interface {
	GetMealy(filename string) (MealyAutomaton, error)
	GetMoore(filename string) (MooreAutomaton, error)
	GetGrammar(filename string) (GrammarAutomaton, error)

	WriteMealy(filename string, automaton MealyAutomaton) error
	WriteMoore(filename string, automaton MooreAutomaton) error
	WriteGrammar(filename string, automaton GrammarAutomaton) error
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

type GrammarAutomaton struct {
	States       []GrammarState
	InputSymbols []string
	Moves        GrammarMoves
}

type MealyMoves map[InitialStateAndInputSymbol]DestinationStateAndSignal

type MooreMoves = map[InitialStateAndInputSymbol]string

type GrammarMoves = map[InitialStateAndInputSymbol][]string

type InitialStateAndInputSymbol struct {
	State  string
	Symbol string
}

type DestinationStateAndSignal struct {
	State  string
	Signal string
}

type GrammarState struct {
	State   string
	IsFinal bool
}
