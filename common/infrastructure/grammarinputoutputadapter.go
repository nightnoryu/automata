package infrastructure

import "automata/common/app"

func NewGrammarInputOutputAdapter() app.GrammarInputOutputAdapter {
	return &grammarInputOutputAdapter{}
}

type grammarInputOutputAdapter struct{}

func (a *grammarInputOutputAdapter) GetGrammar(filename string, side app.GrammarSide) (app.Grammar, error) {
	// TODO implement me
	panic("implement me")
}

func (a *grammarInputOutputAdapter) GetWithEmpty(filename string) (app.EmptyMovesAutomaton, error) {
	// TODO implement me
	panic("implement me")
}

func (a *grammarInputOutputAdapter) WriteWithEmpty(filename string, automaton app.EmptyMovesAutomaton) error {
	// TODO implement me
	panic("implement me")
}
