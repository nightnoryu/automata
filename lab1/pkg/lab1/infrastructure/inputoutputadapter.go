package infrastructure

import "automata/lab1/pkg/lab1/app"

func NewInputOutputAdapter() app.InputOutputAdapter {
	return &inputOutputAdapter{}
}

type inputOutputAdapter struct{}

func (a *inputOutputAdapter) GetMealy(filename string) (app.MealyAutomaton, error) {
	panic("unimplemented")
}

func (a *inputOutputAdapter) GetMoore(filename string) (app.MooreAutomaton, error) {
	panic("unimplemented")
}

func (a *inputOutputAdapter) WriteMealy(filename string, automaton app.MealyAutomaton) error {
	panic("unimplemented")
}

func (a *inputOutputAdapter) WriteMoore(filename string, automaton app.MooreAutomaton) error {
	panic("unimplemented")
}
