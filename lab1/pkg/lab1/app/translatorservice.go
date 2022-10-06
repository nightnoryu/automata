package app

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

	newAutomaton := MooreAutomaton{
		States:       automaton.States,
		InputSymbols: automaton.InputSymbols,
	}

	return s.inputOutputAdapter.WriteMoore(outputFilename, newAutomaton)
}

func (s *TranslatorService) MooreToMealy(inputFilename, outputFilename string) error {
	automaton, err := s.inputOutputAdapter.GetMoore(inputFilename)
	if err != nil {
		return err
	}

	newAutomaton := MealyAutomaton{
		States:       automaton.States,
		InputSymbols: automaton.InputSymbols,
	}

	return s.inputOutputAdapter.WriteMealy(outputFilename, newAutomaton)
}
