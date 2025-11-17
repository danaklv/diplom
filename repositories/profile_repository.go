package repositories

import (
	"database/sql"
	"dl/models"
	"errors"
	"time"
)

type ProfileRepository struct {
	DB *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

// ------------------------ GET PROFILE ------------------------

func (r *ProfileRepository) GetProfileByID(userID int64) (*models.UserProfile, error) {
	var p models.UserProfile
	err := r.DB.QueryRow(`
        SELECT id, username, first_name, last_name, gender, birth_date,
               bio, profile_picture, rating, email, updated_at
        FROM users WHERE id = $1
    `, userID).Scan(
		&p.ID, &p.Username, &p.FirstName, &p.LastName, &p.Gender,
		&p.BirthDate, &p.Bio, &p.ProfilePicture, &p.Rating,
		&p.Email, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &p, err
}

// ------------------------ UPDATE PROFILE ------------------------

func (r *ProfileRepository) UpdateProfile(userID int64, first, last, gender, bio, birth string) error {
	_, err := r.DB.Exec(`
        UPDATE users
        SET first_name = $1, last_name = $2, gender = $3, bio = $4,
            birth_date = $5, updated_at = $6
        WHERE id = $7
    `, first, last, gender, bio, birth, time.Now(), userID)
	return err
}

// ------------------------ UPDATE AVATAR ------------------------

func (r *ProfileRepository) UpdateAvatar(userID int64, filePath string) error {
	_, err := r.DB.Exec(`
        UPDATE users SET profile_picture = $1, updated_at = $2 WHERE id = $3
    `, filePath, time.Now(), userID)
	return err
}

// ------------------------ DELETE ------------------------

func (r *ProfileRepository) DeleteUser(userID int64) error {
	_, err := r.DB.Exec(`DELETE FROM users WHERE id = $1`, userID)
	return err
}
