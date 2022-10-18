package app

import (
	"log"
	"strconv"
	"strings"
)

const newStatesIdentifier = "q"

func NewMinimizerService(inputOutputAdapter InputOutputAdapter) *MinimizerService {
	return &MinimizerService{
		inputOutputAdapter: inputOutputAdapter,
	}
}

type MinimizerService struct {
	inputOutputAdapter InputOutputAdapter
}

func (s *MinimizerService) MinimizeMealy(inputFilename, outputFilename string) error {
	mealyAutomaton, err := s.inputOutputAdapter.GetMealy(inputFilename)
	if err != nil {
		return err
	}

	groupToStatesMap, groupAmount := buildOneEquivalencyGroups(mealyAutomaton)
	for {
		previousGroupAmount := groupAmount

		groupToStatesMap, groupAmount = buildNextEquivalencyGroups(
			groupToStatesMap,
			mealyAutomaton.InputSymbols,
			simplifyMealyMoves(mealyAutomaton.Moves),
		)

		if previousGroupAmount == groupAmount {
			break
		}
	}

	minimizedAutomaton := buildMinimizedMealy(mealyAutomaton, groupToStatesMap)
	return s.inputOutputAdapter.WriteMealy(outputFilename, minimizedAutomaton)
}

func (s *MinimizerService) MinimizeMoore(inputFilename, outputFilename string) error {
	mooreAutomaton, err := s.inputOutputAdapter.GetMoore(inputFilename)
	if err != nil {
		return err
	}

	groupToStatesMap, groupAmount := buildZeroEquivalencyGroups(mooreAutomaton.StateSignals)
	for {
		previousGroupAmount := groupAmount

		groupToStatesMap, groupAmount = buildNextEquivalencyGroups(
			groupToStatesMap,
			mooreAutomaton.InputSymbols,
			mooreAutomaton.Moves,
		)

		if previousGroupAmount == groupAmount {
			break
		}
	}

	minimizedAutomaton := buildMinimizedMoore(mooreAutomaton, groupToStatesMap)
	return s.inputOutputAdapter.WriteMoore(outputFilename, minimizedAutomaton)
}

func buildOneEquivalencyGroups(mealyAutomaton MealyAutomaton) (groupToStatesMap map[int][]string, groupAmount int) {
	stateToGroupHashMap := make(map[string]string)

	for _, sourceState := range mealyAutomaton.States {
		for _, inputSymbol := range mealyAutomaton.InputSymbols {
			key := InitialStateAndInputSymbol{
				State:  sourceState,
				Symbol: inputSymbol,
			}

			destinationSignal := mealyAutomaton.Moves[key].Signal
			stateToGroupHashMap[sourceState] += destinationSignal
		}
	}

	groupHashToStatesMap := buildGroupHashToStatesMap(stateToGroupHashMap)
	groupToStatesMap = make(map[int][]string)

	for _, newStates := range groupHashToStatesMap {
		for _, state := range newStates {
			groupToStatesMap[groupAmount] = append(groupToStatesMap[groupAmount], state)
		}
		groupAmount++
	}

	return groupToStatesMap, groupAmount
}

func buildZeroEquivalencyGroups(stateSignals map[string]string) (groupToStatesMap map[int][]string, groupAmount int) {
	signalToStatesMap := buildSignalToStatesMap(stateSignals)
	groupToStatesMap = make(map[int][]string)

	for _, states := range signalToStatesMap {
		for _, state := range states {
			groupToStatesMap[groupAmount] = append(groupToStatesMap[groupAmount], state)
		}
		groupAmount++
	}

	return groupToStatesMap, groupAmount
}

func buildNextEquivalencyGroups(
	groupToStatesMap map[int][]string,
	inputSymbols []string,
	moves MooreMoves,
) (stateToNewGroupMap map[int][]string, groupAmount int) {
	stateToNewGroupMap = make(map[int][]string)
	stateToGroupMap := buildStateToGroupMap(groupToStatesMap)

	for _, groupStates := range groupToStatesMap {
		stateToGroupHashMap := make(map[string]string)

		for _, sourceState := range groupStates {
			for _, inputSymbol := range inputSymbols {
				key := InitialStateAndInputSymbol{
					State:  sourceState,
					Symbol: inputSymbol,
				}

				destinationState := moves[key]
				destinationGroup := stateToGroupMap[destinationState]

				stateToGroupHashMap[sourceState] += strconv.Itoa(destinationGroup)
			}
		}

		groupHashToStatesMap := buildGroupHashToStatesMap(stateToGroupHashMap)

		for _, newStates := range groupHashToStatesMap {
			for _, state := range newStates {
				stateToNewGroupMap[groupAmount] = append(stateToNewGroupMap[groupAmount], state)
			}
			groupAmount++
		}
	}

	return stateToNewGroupMap, groupAmount
}

func buildMinimizedMealy(mealyAutomaton MealyAutomaton, groupToStatesMap map[int][]string) MealyAutomaton {
	oldStateToNewStateMap := make(map[string]string)
	newStates := make([]string, 0, len(groupToStatesMap))

	for group, oldStates := range groupToStatesMap {
		baseState := oldStates[0]
		newState := getNewStateName(group)

		for _, oldState := range oldStates {
			oldStateToNewStateMap[oldState] = newState
		}

		newStates = append(newStates, newState)

		log.Printf(
			"group %d = { %s }; %s = %s",
			group,
			strings.Join(oldStates, ", "),
			newState,
			baseState,
		)
	}

	newMoves := make(MealyMoves)

	for _, states := range groupToStatesMap {
		baseState := states[0]

		for _, inputSymbol := range mealyAutomaton.InputSymbols {
			key := InitialStateAndInputSymbol{
				State:  baseState,
				Symbol: inputSymbol,
			}
			oldDestinationStateAndSignal := mealyAutomaton.Moves[key]

			newKey := InitialStateAndInputSymbol{
				State:  oldStateToNewStateMap[baseState],
				Symbol: inputSymbol,
			}
			newMoves[newKey] = DestinationStateAndSignal{
				State:  oldStateToNewStateMap[oldDestinationStateAndSignal.State],
				Signal: oldDestinationStateAndSignal.Signal,
			}
		}
	}

	return MealyAutomaton{
		States:       newStates,
		InputSymbols: mealyAutomaton.InputSymbols,
		Moves:        newMoves,
	}
}

func buildMinimizedMoore(mooreAutomaton MooreAutomaton, groupToStatesMap map[int][]string) MooreAutomaton {
	oldStateToNewStateMap := make(map[string]string)
	newStates := make([]string, 0, len(groupToStatesMap))
	newStateSignals := make(map[string]string)

	for group, oldStates := range groupToStatesMap {
		baseState := oldStates[0]
		newState := getNewStateName(group)

		for _, oldState := range oldStates {
			oldStateToNewStateMap[oldState] = newState
		}

		newStates = append(newStates, newState)
		newStateSignals[newState] = mooreAutomaton.StateSignals[baseState]

		log.Printf(
			"group %d = { %s }; %s = %s",
			group,
			strings.Join(oldStates, ", "),
			newState,
			baseState,
		)
	}

	newMoves := make(MooreMoves)

	for _, states := range groupToStatesMap {
		baseState := states[0]

		for _, inputSymbol := range mooreAutomaton.InputSymbols {
			key := InitialStateAndInputSymbol{
				State:  baseState,
				Symbol: inputSymbol,
			}
			oldDestinationState := mooreAutomaton.Moves[key]

			newKey := InitialStateAndInputSymbol{
				State:  oldStateToNewStateMap[baseState],
				Symbol: inputSymbol,
			}
			newMoves[newKey] = oldStateToNewStateMap[oldDestinationState]
		}
	}

	return MooreAutomaton{
		States:       newStates,
		InputSymbols: mooreAutomaton.InputSymbols,
		StateSignals: newStateSignals,
		Moves:        newMoves,
	}
}

func getNewStateName(number int) string {
	return newStatesIdentifier + strconv.Itoa(number)
}

func simplifyMealyMoves(mealyMoves MealyMoves) MooreMoves {
	result := make(MooreMoves)
	for initialStateAndInputSymbol, destinationStateAndSignal := range mealyMoves {
		result[initialStateAndInputSymbol] = destinationStateAndSignal.State
	}

	return result
}

func buildSignalToStatesMap(stateSignals map[string]string) map[string][]string {
	result := make(map[string][]string)
	for state, signal := range stateSignals {
		result[signal] = append(result[signal], state)
	}

	return result
}

func buildStateToGroupMap(groupToStatesMap map[int][]string) map[string]int {
	result := make(map[string]int)
	for group, states := range groupToStatesMap {
		for _, state := range states {
			result[state] = group
		}
	}

	return result
}

func buildGroupHashToStatesMap(stateToGroupHashMap map[string]string) map[string][]string {
	result := make(map[string][]string)
	for state, groupHash := range stateToGroupHashMap {
		result[groupHash] = append(result[groupHash], state)
	}

	return result
}
