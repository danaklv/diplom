package services

import (
	"dl/models"
	"dl/repositories"
)

type RatingService struct {
	Repo *repositories.RatingRepository
}

func NewRatingService(repo *repositories.RatingRepository) *RatingService {
	return &RatingService{Repo: repo}
}

func (s *RatingService) AddEcoAction(userID, actionID int64) error {
	points, err := s.Repo.GetActionPoints(actionID)
	if err != nil {
		return err
	}

	oldLevel, _ := s.Repo.GetUserLevel(userID)

	if err := s.Repo.AddUserAction(userID, actionID, points); err != nil {
		return err
	}

	if err := s.Repo.UpdateRating(userID, points); err != nil {
		return err
	}

	if err := s.UpdateUserLevel(userID); err != nil {
		return err
	}

	newLevel, _ := s.Repo.GetUserLevel(userID)

	if newLevel > oldLevel {
		// тут можно вызывать NotificationsService
	}

	return nil
}

func (s *RatingService) GetUserActions(userID int64) ([]models.UserAction, error) {
	return s.Repo.GetUserActions(userID)
}

func (s *RatingService) GetLeaderboard(limit int) ([]models.LeaderboardEntry, error) {
	return s.Repo.GetLeaderboard(limit)
}

func (s *RatingService) UpdateUserLevel(userID int64) error {
	rating := 0
	// тут можно сделать отдельный метод GetRating
	// но можно считать уровень по ситуации

	level := 1
	league := "Green Seed"

	switch {
	case rating >= 1000:
		level, league = 5, "Earth Legend"
	case rating >= 500:
		level, league = 4, "Planet Guardian"
	case rating >= 250:
		level, league = 3, "Nature Keeper"
	case rating >= 100:
		level, league = 2, "Eco Enthusiast"
	}

	return s.Repo.UpdateLevel(userID, level, league)
}
