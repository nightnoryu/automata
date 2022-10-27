package app

func NewDeterminatorService() *DeterminatorService {
	return &DeterminatorService{}
}

type DeterminatorService struct{}

func (s *DeterminatorService) Determinate() error {
	// TODO: determinate by model
	return nil
}
