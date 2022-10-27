package app

import (
	"log"
	"sort"
	"strconv"

	"automata/common/app"
)

const newStatesIdentifier = "S"

func NewTranslatorService(inputOutputAdapter app.AutomataInputOutputAdapter) *TranslatorService {
	return &TranslatorService{
		inputOutputAdapter: inputOutputAdapter,
	}
}

type TranslatorService struct {
	inputOutputAdapter app.AutomataInputOutputAdapter
}

func (s *TranslatorService) MealyToMoore(inputFilename, outputFilename string) error {
	mealyAutomaton, err := s.inputOutputAdapter.GetMealy(inputFilename)
	if err != nil {
		return err
	}

	newStateToOldStateAndSignalMap := buildNewMooreStates(mealyAutomaton)
	states := getMooreStates(newStateToOldStateAndSignalMap)

	mooreAutomaton := app.MooreAutomaton{
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

	mealyAutomaton := app.MealyAutomaton{
		States:       mooreAutomaton.States,
		InputSymbols: mooreAutomaton.InputSymbols,
		Moves:        getMealyMoves(mooreAutomaton.Moves, mooreAutomaton.StateSignals),
	}

	return s.inputOutputAdapter.WriteMealy(outputFilename, mealyAutomaton)
}

func buildNewMooreStates(mealyAutomaton app.MealyAutomaton) map[string]app.DestinationStateAndSignal {
	processedStates := make(map[app.DestinationStateAndSignal]bool)

	result := make(map[string]app.DestinationStateAndSignal)
	counter := 0
	for _, inputSymbol := range mealyAutomaton.InputSymbols {
		for _, state := range mealyAutomaton.States {
			key := app.InitialStateAndInputSymbol{
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

func getMooreStates(newStateToOldStateAndSignalMap map[string]app.DestinationStateAndSignal) []string {
	result := make([]string, 0, len(newStateToOldStateAndSignalMap))
	for state := range newStateToOldStateAndSignalMap {
		result = append(result, state)
	}

	sort.Strings(result)

	return result
}

func getMooreStateSignals(newStateToOldStateAndSignalMap map[string]app.DestinationStateAndSignal) map[string]string {
	result := make(map[string]string)
	for newState, oldStateAndSignal := range newStateToOldStateAndSignalMap {
		result[newState] = oldStateAndSignal.Signal
	}

	return result
}

func getMooreMoves(
	mealyAutomaton app.MealyAutomaton,
	states []string,
	stateToOldStateAndSignalMap map[string]app.DestinationStateAndSignal,
) map[app.InitialStateAndInputSymbol]string {
	oldStateToStateMap := getOldStateAndSignalToStateMap(stateToOldStateAndSignalMap)

	result := make(map[app.InitialStateAndInputSymbol]string)
	for _, state := range states {
		oldState := stateToOldStateAndSignalMap[state].State
		for _, symbol := range mealyAutomaton.InputSymbols {
			oldDestination := mealyAutomaton.Moves[app.InitialStateAndInputSymbol{
				State:  oldState,
				Symbol: symbol,
			}]

			result[app.InitialStateAndInputSymbol{
				State:  state,
				Symbol: symbol,
			}] = oldStateToStateMap[oldDestination]
		}
	}

	return result
}

func getOldStateAndSignalToStateMap(
	stateToOldStateAndSignalMap map[string]app.DestinationStateAndSignal,
) map[app.DestinationStateAndSignal]string {
	result := make(map[app.DestinationStateAndSignal]string)
	for state, oldStateAndSignal := range stateToOldStateAndSignalMap {
		result[oldStateAndSignal] = state
	}

	return result
}

func getMealyMoves(
	moves map[app.InitialStateAndInputSymbol]string,
	stateToSignalMap map[string]string,
) map[app.InitialStateAndInputSymbol]app.DestinationStateAndSignal {
	result := make(map[app.InitialStateAndInputSymbol]app.DestinationStateAndSignal)
	for initialStateAndInputSymbol, destinationState := range moves {
		result[initialStateAndInputSymbol] = app.DestinationStateAndSignal{
			State:  destinationState,
			Signal: stateToSignalMap[destinationState],
		}
	}

	return result
}
