package repositories

import (
	"database/sql"
	"errors"
	"time"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// ------------------------ CREATE USER ------------------------

func (r *UserRepository) CreateUser(username, email string, passwordHash []byte) (int64, error) {
	var id int64
	err := r.DB.QueryRow(`
        INSERT INTO users (username, email, password_hash, is_verified)
        VALUES ($1, $2, $3, false)
        RETURNING id
    `, username, email, passwordHash).Scan(&id)
	return id, err
}

// ------------------------ GET USER BY EMAIL ------------------------

func (r *UserRepository) GetUserByEmail(email string) (int64, []byte, bool, error) {
	var (
		id       int64
		password []byte
		verified bool
	)

	err := r.DB.QueryRow(`
        SELECT id, password_hash, is_verified 
        FROM users WHERE email = $1
    `, email).Scan(&id, &password, &verified)

	if err == sql.ErrNoRows {
		return 0, nil, false, errors.New("user not found")
	}

	return id, password, verified, err
}

// ------------------------ EMAIL VERIFICATION ------------------------

func (r *UserRepository) StoreVerificationCode(userID int64, code string, expires time.Time) error {
	_, err := r.DB.Exec(`
        INSERT INTO email_verifications (user_id, code, expires_at)
        VALUES ($1, $2, $3)
    `, userID, code, expires)
	return err
}

func (r *UserRepository) GetUserByVerificationCode(code string) (int64, time.Time, error) {
	var (
		userID    int64
		expiresAt time.Time
	)

	err := r.DB.QueryRow(`
        SELECT user_id, expires_at FROM email_verifications 
        WHERE code = $1
    `, code).Scan(&userID, &expiresAt)

	return userID, expiresAt, err
}

func (r *UserRepository) DeleteVerificationCode(userID int64) error {
	_, err := r.DB.Exec(`DELETE FROM email_verifications WHERE user_id = $1`, userID)
	return err
}

func (r *UserRepository) SetUserVerified(userID int64) error {
	_, err := r.DB.Exec(`UPDATE users SET is_verified = true WHERE id = $1`, userID)
	return err
}

// ------------------------ PASSWORD RESET ------------------------

func (r *UserRepository) CreatePasswordReset(userID int64, token string, expires time.Time) error {
	_, err := r.DB.Exec(`
        INSERT INTO password_resets (user_id, token, expires_at, used)
        VALUES ($1, $2, $3, false)
    `, userID, token, expires)
	return err
}

func (r *UserRepository) GetPasswordReset(token string) (int64, time.Time, bool, error) {
	var (
		userID  int64
		expires time.Time
		used    bool
	)

	err := r.DB.QueryRow(`
        SELECT user_id, expires_at, used 
        FROM password_resets WHERE token = $1
    `, token).Scan(&userID, &expires, &used)

	return userID, expires, used, err
}

func (r *UserRepository) MarkResetTokenUsed(token string) error {
	_, err := r.DB.Exec(`
        UPDATE password_resets SET used = true WHERE token = $1
    `, token)
	return err
}

func (r *UserRepository) UpdateUserPassword(userID int64, hashed []byte) error {
	_, err := r.DB.Exec(`
        UPDATE users SET password_hash = $1 WHERE id = $2
    `, hashed, userID)
	return err
}
