package services

import (
	"database/sql"
	"errors"
	"time"
)

type UserProfile struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Gender         string    `json:"gender"`
	BirthDate      string    `json:"birth_date"`
	Bio            string    `json:"bio"`
	ProfilePicture string    `json:"profile_picture"`
	Rating         int       `json:"rating"`
	Email          string    `json:"email"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProfileService struct {
	DB *sql.DB
}

// Получение профиля
func (s *ProfileService) GetProfile(userID int64) (*UserProfile, error) {
	var p UserProfile
	err := s.DB.QueryRow(`
		SELECT id, username, first_name, last_name, gender, birth_date,
		       bio, profile_picture, rating, email, updated_at
		FROM users WHERE id = $1
	`, userID).Scan(
		&p.ID, &p.Username, &p.FirstName, &p.LastName, &p.Gender,
		&p.BirthDate, &p.Bio, &p.ProfilePicture, &p.Rating, &p.Email, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}
	return &p, nil
}

// Обновление профиля
func (s *ProfileService) UpdateProfile(userID int64, firstName, lastName, gender, bio, birthDate string) error {
	_, err := s.DB.Exec(`
		UPDATE users
		SET first_name = $1, last_name = $2, gender = $3, bio = $4,
		    birth_date = $5, updated_at = $6
		WHERE id = $7
	`, firstName, lastName, gender, bio, birthDate, time.Now(), userID)
	return err
}

// Обновление фото
func (s *ProfileService) UpdateProfilePicture(userID int64, filePath string) error {
	_, err := s.DB.Exec(`
		UPDATE users SET profile_picture = $1, updated_at = $2 WHERE id = $3
	`, filePath, time.Now(), userID)
	return err
}


// Удаление профиля
func (s *ProfileService) DeleteProfile(userID int64) error {
	_, err := s.DB.Exec("DELETE FROM users WHERE id = $1", userID)
	return err
}

