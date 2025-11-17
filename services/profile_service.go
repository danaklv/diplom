package services

import (
	"dl/models"
	"dl/repositories"
)

type ProfileService struct {
	Repo *repositories.ProfileRepository
}

func NewProfileService(repo *repositories.ProfileRepository) *ProfileService {
	return &ProfileService{Repo: repo}
}

func (s *ProfileService) GetProfile(userID int64) (*models.UserProfile, error) {
	return s.Repo.GetProfileByID(userID)
}

func (s *ProfileService) UpdateProfile(userID int64, first, last, gender, bio, birth string) error {
	return s.Repo.UpdateProfile(userID, first, last, gender, bio, birth)
}

func (s *ProfileService) UpdateProfilePicture(userID int64, filePath string) error {
	return s.Repo.UpdateAvatar(userID, filePath)
}

func (s *ProfileService) DeleteProfile(userID int64) error {
	return s.Repo.DeleteUser(userID)
}
