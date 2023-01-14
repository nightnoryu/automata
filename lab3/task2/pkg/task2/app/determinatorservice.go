package app

import (
	"automata/common/app"
	"fmt"
)

const emptyMoveIndicator = "e"

func NewDeterminatorService(inputOutputAdapter app.FiniteInputOutputAdapter) *DeterminatorService {
	return &DeterminatorService{
		inputOutputAdapter: inputOutputAdapter,
	}
}

type DeterminatorService struct {
	inputOutputAdapter app.FiniteInputOutputAdapter
}

func (s *DeterminatorService) Determinate(inputFilename, outputFilename string) error {
	automaton, err := s.inputOutputAdapter.GetNonDeterministicFinite(inputFilename)
	if err != nil {
		return err
	}

	closures := buildClosures(automaton)
	fmt.Println(closures)

	return nil
}

func buildClosures(automaton app.NonDeterministicFiniteAutomaton) map[string][]string {
	result := make(map[string][]string)
	for _, state := range automaton.States {
		key := app.InitialStateAndInputSymbol{
			State:  state,
			Symbol: emptyMoveIndicator,
		}

		for _, destinationState := range automaton.Moves[key] {
			result[state] = append(result[state], destinationState)
		}
	}

	for recurseClosures(result) {
	}

	return result
}

func recurseClosures(result map[string][]string) bool {
	foundDeeper := false
	for state, closure := range result {
		for _, closureState := range closure {
			for _, transitiveState := range result[closureState] {
				if inSlice(result[state], transitiveState) {
					continue
				}

				result[state] = append(result[state], transitiveState)
				foundDeeper = true
			}
		}
	}

	return foundDeeper
}

func inSlice(haystack []string, needle string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}
