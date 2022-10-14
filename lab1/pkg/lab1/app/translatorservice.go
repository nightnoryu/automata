package app

import (
	"log"
	"sort"
	"strconv"
)

const newStatesIdentifier = "S"

func NewTranslatorService(inputOutputAdapter InputOutputAdapter) *TranslatorService {
	return &TranslatorService{
		inputOutputAdapter: inputOutputAdapter,
	}
}

type TranslatorService struct {
	inputOutputAdapter InputOutputAdapter
}

func (s *TranslatorService) MealyToMoore(inputFilename, outputFilename string) error {
	mealyAutomaton, err := s.inputOutputAdapter.GetMealy(inputFilename)
	if err != nil {
		return err
	}

	newStateToOldStateAndSignalMap := buildNewMooreStates(mealyAutomaton)
	states := getMooreStates(newStateToOldStateAndSignalMap)

	mooreAutomaton := MooreAutomaton{
		States:       states,
		InputSymbols: mealyAutomaton.InputSymbols,
		StateSignals: getMooreStateSignals(newStateToOldStateAndSignalMap),
		Moves:        getMooreMoves(mealyAutomaton, states, newStateToOldStateAndSignalMap),
	}

	return s.inputOutputAdapter.WriteMoore(outputFilename, mooreAutomaton)
}

func (s *TranslatorService) MooreToMealy(inputFilename, outputFilename string) error {
	mooreAutomaton, err := s.inputOutputAdapter.GetMoore(inputFilename)
	if err != nil {
		return err
	}

	mealyAutomaton := MealyAutomaton{
		States:       mooreAutomaton.States,
		InputSymbols: mooreAutomaton.InputSymbols,
		Moves:        getMealyMoves(mooreAutomaton.Moves, mooreAutomaton.StateSignals),
	}

	return s.inputOutputAdapter.WriteMealy(outputFilename, mealyAutomaton)
}

func buildNewMooreStates(mealyAutomaton MealyAutomaton) map[string]DestinationStateAndSignal {
	processedStates := make(map[DestinationStateAndSignal]bool)

	result := make(map[string]DestinationStateAndSignal)
	counter := 0
	for _, inputSymbol := range mealyAutomaton.InputSymbols {
		for _, state := range mealyAutomaton.States {
			key := InitialStateAndInputSymbol{
				State:  state,
				Symbol: inputSymbol,
			}

			destinationStateAndSignal := mealyAutomaton.Moves[key]

			if processedStates[destinationStateAndSignal] {
				continue
			}

			stateName := getNewStateName(counter)
			result[stateName] = destinationStateAndSignal

			counter++
			processedStates[destinationStateAndSignal] = true

			log.Printf("%s = %s/%s", stateName, destinationStateAndSignal.State, destinationStateAndSignal.Signal)
		}
	}

	return result
}

func getNewStateName(number int) string {
	return newStatesIdentifier + strconv.Itoa(number)
}

func getMooreStates(newStateToOldStateAndSignalMap map[string]DestinationStateAndSignal) []string {
	result := make([]string, 0, len(newStateToOldStateAndSignalMap))
	for state := range newStateToOldStateAndSignalMap {
		result = append(result, state)
	}

	sort.Strings(result)

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
	mealyAutomaton MealyAutomaton,
	states []string,
	stateToOldStateAndSignalMap map[string]DestinationStateAndSignal,
) map[InitialStateAndInputSymbol]string {
	oldStateToStateMap := getOldStateToStateMap(stateToOldStateAndSignalMap)

	result := make(map[InitialStateAndInputSymbol]string)
	for _, state := range states {
		oldState := stateToOldStateAndSignalMap[state].State
		for _, symbol := range mealyAutomaton.InputSymbols {
			oldDestination := mealyAutomaton.Moves[InitialStateAndInputSymbol{
				State:  oldState,
				Symbol: symbol,
			}]

			result[InitialStateAndInputSymbol{
				State:  state,
				Symbol: symbol,
			}] = oldStateToStateMap[oldDestination.State]
		}
	}

	return result
}

func getOldStateToStateMap(stateToOldStateAndSignalMap map[string]DestinationStateAndSignal) map[string]string {
	result := make(map[string]string)
	for state, oldStateAndSignal := range stateToOldStateAndSignalMap {
		result[oldStateAndSignal.State] = state
	}

	return result
}

func getMealyMoves(
	moves map[InitialStateAndInputSymbol]string,
	stateToSignalMap map[string]string,
) map[InitialStateAndInputSymbol]DestinationStateAndSignal {
	result := make(map[InitialStateAndInputSymbol]DestinationStateAndSignal)
	for initialStateAndInputSymbol, destinationState := range moves {
		result[initialStateAndInputSymbol] = DestinationStateAndSignal{
			State:  destinationState,
			Signal: stateToSignalMap[destinationState],
		}
	}

	return result
}
