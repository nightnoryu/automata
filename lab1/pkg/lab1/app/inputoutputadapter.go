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
	Moves        map[string]map[string]DestinationStateAndSignal // State -> Input symbol -> New state and signal
}

type DestinationStateAndSignal struct {
	State  string
	Signal string
}

type MooreAutomaton struct {
	States       []string
	InputSymbols []string
	StateSignals map[string]string            // State -> signal
	Moves        map[string]map[string]string // State -> Input symbol -> New state
}
