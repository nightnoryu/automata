package app

import (
	"fmt"
	"strconv"
)

const mooreStatesLetter = "q"

func NewTranslatorService(inputOutputAdapter InputOutputAdapter) *TranslatorService {
	return &TranslatorService{
		inputOutputAdapter: inputOutputAdapter,
	}
}

type TranslatorService struct {
	inputOutputAdapter InputOutputAdapter
}

func (s *TranslatorService) MealyToMoore(inputFilename, outputFilename string) error {
	automaton, err := s.inputOutputAdapter.GetMealy(inputFilename)
	if err != nil {
		return err
	}

	newStateToOldStateAndSignalMap := buildNewMooreStates(automaton.Moves)

	newStates := getMooreStates(newStateToOldStateAndSignalMap)

	newAutomaton := MooreAutomaton{
		States:       newStates,
		InputSymbols: automaton.InputSymbols,
		StateSignals: getMooreStateSignals(newStateToOldStateAndSignalMap),
		Moves:        getMooreMoves(newStates, newStateToOldStateAndSignalMap, automaton.InputSymbols, automaton.Moves),
	}

	return s.inputOutputAdapter.WriteMoore(outputFilename, newAutomaton)
}

func (s *TranslatorService) MooreToMealy(inputFilename, outputFilename string) error {
	automaton, err := s.inputOutputAdapter.GetMoore(inputFilename)
	if err != nil {
		return err
	}

	fmt.Println(automaton)

	newAutomaton := MealyAutomaton{
		States:       automaton.States,
		InputSymbols: automaton.InputSymbols,
	}

	// TODO

	return s.inputOutputAdapter.WriteMealy(outputFilename, newAutomaton)
}

func buildNewMooreStates(
	moves map[InitialStateAndInputSymbol]DestinationStateAndSignal,
) map[string]DestinationStateAndSignal {
	result := make(map[string]DestinationStateAndSignal)
	processedStates := make(map[DestinationStateAndSignal]bool)
	counter := 1
	for _, destinationStateAndSignal := range moves {
		if processedStates[destinationStateAndSignal] {
			continue
		}
		stateName := getMooreStateName(counter)
		result[stateName] = destinationStateAndSignal
		counter++
		processedStates[destinationStateAndSignal] = true
	}

	return result
}

func getMooreStateName(number int) string {
	return mooreStatesLetter + strconv.Itoa(number)
}

func getMooreStates(newStateToOldStateAndSignalMap map[string]DestinationStateAndSignal) []string {
	result := make([]string, 0, len(newStateToOldStateAndSignalMap))
	for state := range newStateToOldStateAndSignalMap {
		result = append(result, state)
	}

	return result
}

func getMooreStateSignals(newStateToOldStateAndSignalMap map[string]DestinationStateAndSignal) map[string]string {
	result := make(map[string]string)
	for newState, oldStateAndSignal := range newStateToOldStateAndSignalMap {
		result[newState] = oldStateAndSignal.Signal
	}

	return result
}

func getMooreMoves(
	newStates []string,
	newStateToOldStateAndSignalMap map[string]DestinationStateAndSignal,
	inputSymbols []string,
	oldMoves map[InitialStateAndInputSymbol]DestinationStateAndSignal,
) map[InitialStateAndInputSymbol]string {
	oldStateToNewStateMap := make(map[string]string)
	for newState, oldStateAndSignal := range newStateToOldStateAndSignalMap {
		oldStateToNewStateMap[oldStateAndSignal.State] = newState
	}

	result := make(map[InitialStateAndInputSymbol]string)
	for _, newState := range newStates {
		oldState := newStateToOldStateAndSignalMap[newState].State
		for _, symbol := range inputSymbols {
			key := InitialStateAndInputSymbol{
				State:  oldState,
				Symbol: symbol,
			}

			oldDestination := oldMoves[key]

			newKey := InitialStateAndInputSymbol{
				State:  newState,
				Symbol: symbol,
			}
			result[newKey] = oldStateToNewStateMap[oldDestination.State]
		}
	}

	return result
}
