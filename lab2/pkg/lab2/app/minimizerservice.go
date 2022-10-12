package app

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

	// TODO

	return s.inputOutputAdapter.WriteMoore(outputFilename, mooreAutomaton)
}
