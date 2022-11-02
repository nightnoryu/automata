package app

import "automata/common/app"

func NewDeterminatorService(inputOutputAdapter app.GrammarInputOutputAdapter) *DeterminatorService {
	return &DeterminatorService{
		inputOutputAdapter: inputOutputAdapter,
	}
}

type DeterminatorService struct {
	inputOutputAdapter app.GrammarInputOutputAdapter
}

func (s *DeterminatorService) Determinate(inputFilename, outputFilename string) error {
	automaton, err := s.inputOutputAdapter.GetFinite(inputFilename)
	if err != nil {
		return err
	}

	// TODO: implement determination with empty symbols

	return s.inputOutputAdapter.WriteFinite(outputFilename, automaton)
}

func buildStateToClosureMap(automaton app.FiniteAutomaton) map[string][]string {
	result := make(map[string][]string)
	for _, state := range automaton.States {
		key := app.InitialStateAndInputSymbol{
			State:  state.State,
			Symbol: app.EmptySymbol,
		}

		result[state.State] = append(result[state.State], automaton.Moves[key])
	}

	return result
}
