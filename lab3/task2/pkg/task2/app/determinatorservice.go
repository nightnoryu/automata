package app

import (
	"automata/common/app"
	"sort"
	"strconv"
	"strings"
)

const (
	emptyMoveIndicator = "e"
	newStatesName      = "S"
)

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

	stateHashMap := make(map[string]stateCombination)

	var newStates []string
	newFinalStates := make(map[string]bool)
	newMoves := make(app.DeterministicMoves)

	var stateQueue [][]string
	stateQueue = append(stateQueue, []string{automaton.States[0]})

	for len(stateQueue) > 0 {
		states := stateQueue[len(stateQueue)-1]
		stateQueue = stateQueue[1:]

		currentState := getFullState(states, closures, automaton.FinalStates)

		stateHash := getStatesHash(currentState.States)
		if _, ok := stateHashMap[stateHash]; ok {
			continue
		}
		stateHashMap[stateHash] = currentState
		newStates = append(newStates, stateHash)
		newFinalStates[stateHash] = currentState.IsFinal

		for _, symbol := range automaton.InputSymbols {
			if symbol == emptyMoveIndicator {
				continue
			}

			newKey := app.InitialStateAndInputSymbol{
				State:  stateHash,
				Symbol: symbol,
			}

			var destinationStates []string
			for _, state := range currentState.States {
				key := app.InitialStateAndInputSymbol{
					State:  state,
					Symbol: symbol,
				}

				for _, initialDestinationState := range automaton.Moves[key] {
					destinationStates = append(destinationStates, initialDestinationState)
				}
			}

			if len(destinationStates) != 0 {
				stateQueue = append(stateQueue, destinationStates)
				destinationState := getFullState(destinationStates, closures, automaton.FinalStates)
				newMoves[newKey] = strings.Join(destinationState.States, ",")
			}
		}
	}

	result := app.FiniteAutomaton{
		States:       newStates,
		InputSymbols: removeEmptyInputSymbol(automaton.InputSymbols),
		FinalStates:  newFinalStates,
		Moves:        newMoves,
	}

	return s.inputOutputAdapter.WriteFinite(outputFilename, result)
}

func removeEmptyInputSymbol(symbols []string) []string {
	result := make([]string, 0, len(symbols)-1)
	for _, symbol := range symbols {
		if symbol == emptyMoveIndicator {
			continue
		}
		result = append(result, symbol)
	}

	return result
}

func getStatesHash(states []string) string {
	result := ""
	for _, state := range states {
		result += state
	}

	return result
}

func getFullState(
	states []string,
	closures map[string]stateCombination,
	finalStates map[string]bool,
) stateCombination {
	stateMap := make(map[string]bool)
	for _, state := range states {
		stateMap[state] = true
		if closure, ok := closures[state]; ok {
			for _, closureState := range closure.States {
				stateMap[closureState] = true
			}
		}
	}

	resultStates := make([]string, 0, len(stateMap))
	isFinal := false
	for state := range stateMap {
		resultStates = append(resultStates, state)
		if finalStates[state] {
			isFinal = true
		}
	}

	sort.Strings(resultStates)

	return stateCombination{
		IsFinal: isFinal,
		States:  resultStates,
	}
}

func buildClosures(automaton app.NonDeterministicFiniteAutomaton) map[string]stateCombination {
	flatClosures := make(map[string][]string)
	for _, state := range automaton.States {
		key := app.InitialStateAndInputSymbol{
			State:  state,
			Symbol: emptyMoveIndicator,
		}

		for _, destinationState := range automaton.Moves[key] {
			flatClosures[state] = append(flatClosures[state], destinationState)
		}
	}

	for recurseClosures(flatClosures) {
	}

	result := make(map[string]stateCombination)
	for state, closureStates := range flatClosures {
		isFinal := false
		for _, closureState := range closureStates {
			if automaton.FinalStates[closureState] {
				isFinal = true
				break
			}
		}

		result[state] = stateCombination{
			IsFinal: isFinal,
			States:  closureStates,
		}
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

func buildStateName(number int) string {
	return newStatesName + strconv.Itoa(number)
}

func inSlice(haystack []string, needle string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}

type stateCombination struct {
	IsFinal bool
	States  []string
}
