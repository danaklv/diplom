package services

import (
	"dl/models"
	"dl/repositories"
	"dl/utils"
)

type EcoService struct {
	Repo *repositories.EcoRepository
}

func NewEcoService(repo *repositories.EcoRepository) *EcoService {
	return &EcoService{Repo: repo}
}

func (s *EcoService) SubmitAnswers(userID int64, answers map[int]int) (*models.EcoResult, error) {
	// Sum score
	total := 0
	for _, v := range answers {
		total += v
	}

	// Determine category
	category, description := utils.CalculateEcoCategory(total)

	// Save answers
	if err := s.Repo.SaveAnswers(userID, answers); err != nil {
		return nil, err
	}

	// Save result
	if err := s.Repo.SaveResult(userID, total, category, description); err != nil {
		return nil, err
	}

	// Return result
	return &models.EcoResult{
		UserID:      userID,
		TotalScore:  total,
		Category:    category,
		Description: description,
	}, nil
}

func (s *EcoService) GetLatest(userID int64) (*models.EcoResult, error) {
	return s.Repo.GetLatestResult(userID)
}

func (s *EcoService) GetQuestions() ([]models.EcoQuestion, error) {
	return s.Repo.GetQuestions()
}
