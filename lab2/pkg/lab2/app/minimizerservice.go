package app

import (
	"fmt"
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

	// TODO

	return s.inputOutputAdapter.WriteMealy(outputFilename, mealyAutomaton)
}

func (s *MinimizerService) MinimizeMoore(inputFilename, outputFilename string) error {
	mooreAutomaton, err := s.inputOutputAdapter.GetMoore(inputFilename)
	if err != nil {
		return err
	}

	stateToGroupMap, groupAmount := buildZeroEquivalencyGroups(mooreAutomaton.StateSignals)
	for {
		previousGroupAmount := groupAmount

		stateToGroupMap, groupAmount = buildNextEquivalencyGroups(
			stateToGroupMap,
			mooreAutomaton.InputSymbols,
			mooreAutomaton.Moves,
		)

		if previousGroupAmount == groupAmount {
			break
		}
	}

	fmt.Println(stateToGroupMap, groupAmount)

	return s.inputOutputAdapter.WriteMoore(outputFilename, mooreAutomaton)
}

func buildZeroEquivalencyGroups(stateSignals map[string]string) (stateToGroupMap map[string]int, groupAmount int) {
	signalToStatesMap := make(map[string][]string)
	for state, signal := range stateSignals {
		signalToStatesMap[signal] = append(signalToStatesMap[signal], state)
	}

	stateToGroupMap = make(map[string]int)
	for _, states := range signalToStatesMap {
		for _, state := range states {
			stateToGroupMap[state] = groupAmount
		}
		groupAmount++
	}

	return stateToGroupMap, groupAmount
}

func buildNextEquivalencyGroups(
	stateToGroupMap map[string]int,
	inputSymbols []string,
	moves map[InitialStateAndInputSymbol]string,
) (stateToNewGroupMap map[string]int, groupAmount int) {
	groupToStatesMap := make(map[int][]string)
	for state, group := range stateToGroupMap {
		groupToStatesMap[group] = append(groupToStatesMap[group], state)
	}

	stateToNewGroupMap = make(map[string]int)
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

		groupHashToStatesMap := make(map[string][]string)
		for state, groupHash := range stateToGroupHashMap {
			groupHashToStatesMap[groupHash] = append(groupHashToStatesMap[groupHash], state)
		}

		for _, newStates := range groupHashToStatesMap {
			for _, state := range newStates {
				stateToNewGroupMap[state] = groupAmount
			}
			groupAmount++
		}
	}

	return stateToNewGroupMap, groupAmount
}
