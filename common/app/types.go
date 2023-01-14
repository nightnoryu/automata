package app

type NonDeterministicMoves = map[InitialStateAndInputSymbol][]string

type DeterministicMoves = map[InitialStateAndInputSymbol]string

type MovesWithSignals map[InitialStateAndInputSymbol]DestinationStateAndSignal

type InitialStateAndInputSymbol struct {
	State  string
	Symbol string
}

type DestinationStateAndSignal struct {
	State  string
	Signal string
}
