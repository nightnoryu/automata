package app

import "automata/common/app"

func NewDeterminatorService() *DeterminatorService {
	return &DeterminatorService{}
}

type DeterminatorService struct{}

func (s *DeterminatorService) Determinate(automaton app.NonDeterministicGrammarAutomaton) app.GrammarAutomaton {
	// TODO: determinate by model
	return app.GrammarAutomaton{}
}
