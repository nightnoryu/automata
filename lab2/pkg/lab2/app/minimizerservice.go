package app

import (
	"strconv"
)

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

	// TODO: build new automaton

	return s.inputOutputAdapter.WriteMealy(outputFilename, mealyAutomaton)
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

	// TODO: build new automaton

	return s.inputOutputAdapter.WriteMoore(outputFilename, mooreAutomaton)
}

func buildOneEquivalencyGroups(mealyAutomaton MealyAutomaton) (stateToGroupMap map[int][]string, groupAmount int) {
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
	stateToGroupMap = make(map[int][]string)

	for _, newStates := range groupHashToStatesMap {
		for _, state := range newStates {
			stateToGroupMap[groupAmount] = append(stateToGroupMap[groupAmount], state)
		}
		groupAmount++
	}

	return stateToGroupMap, groupAmount
}

func buildZeroEquivalencyGroups(stateSignals map[string]string) (stateToGroupMap map[int][]string, groupAmount int) {
	signalToStatesMap := buildSignalToStatesMap(stateSignals)
	stateToGroupMap = make(map[int][]string)

	for _, states := range signalToStatesMap {
		for _, state := range states {
			stateToGroupMap[groupAmount] = append(stateToGroupMap[groupAmount], state)
		}
		groupAmount++
	}

	return stateToGroupMap, groupAmount
}

func buildNextEquivalencyGroups(
	groupToStatesMap map[int][]string,
	inputSymbols []string,
	moves map[InitialStateAndInputSymbol]string,
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

func simplifyMealyMoves(
	mealyMoves map[InitialStateAndInputSymbol]DestinationStateAndSignal,
) map[InitialStateAndInputSymbol]string {
	result := make(map[InitialStateAndInputSymbol]string)
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
