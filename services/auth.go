package services

import (
	"crypto/rand"
	"database/sql"
	"dl/utils"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AuthServicer interface {
	Register(username, email, password string) (string, string, error)
	Login(email, password string) (string, string, error)
	// VerifyEmail(code string) (string, string, error)
	// ResendVerification(email string) error
	// RequestPasswordReset(email string) error
	// ResetPassword(token, newPassword string) error
}

type AuthService struct {
	DB *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{DB: db}
}


// --------------------------- REGISTER ---------------------------
func (s *AuthService) Register(username, email, password string) (string, string, error) {
	// Валидации
	if err := utils.ValidateEmail(email); err != nil {
		return "", "", err
	}
	if err := utils.ValidateUsername(username); err != nil {
		return "", "", err
	}
	if err := utils.ValidatePassword(password); err != nil {
		return "", "", err
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	// Создаём пользователя
	var userID int64
	err = s.DB.QueryRow(`
		INSERT INTO users (username, email, password_hash, is_verified)
		VALUES ($1, $2, $3, false)
		RETURNING id
	`, username, email, hashedPassword).Scan(&userID)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			if strings.Contains(pgErr.Message, "users_email_key") {
				return "", "", errors.New("email already exists")
			}
			if strings.Contains(pgErr.Message, "users_username_key") {
				return "", "", errors.New("username already exists")
			}
		}
		return "", "", err
	}

	// Генерируем код подтверждения
	code := generateVerificationCode()
	expiration := time.Now().Add(10 * time.Minute)

	_, err = s.DB.Exec(`
		INSERT INTO email_verifications (user_id, code, expires_at)
		VALUES ($1, $2, $3)
	`, userID, code, expiration)
	if err != nil {
		return "", "", err
	}

	// Отправляем письмо
	go utils.SendVerificationEmail(email, code)

	return "", "", nil
}

// логин
func (s *AuthService) Login(email, password string) (string, string, error) {
	var userID int64
	var dbPassword []byte
	var isVerified bool

	err := s.DB.QueryRow(`
		SELECT id, password_hash, is_verified FROM users WHERE email = $1
	`, email).Scan(&userID, &dbPassword, &isVerified)

	if err == sql.ErrNoRows {
		return "", "", errors.New("invalid email or password")
	}
	if err != nil {
		return "", "", err
	}

	if bcrypt.CompareHashAndPassword(dbPassword, []byte(password)) != nil {
		return "", "", errors.New("invalid email or password")
	}

	if !isVerified {
		return "", "", errors.New("email not verified")
	}

	return utils.GenerateTokens(userID)
}

// подтв почты
func (s *AuthService) VerifyEmail(code string) (string, string, error) {
	var userID int64
	var expiresAt time.Time

	err := s.DB.QueryRow(`
		SELECT user_id, expires_at FROM email_verifications WHERE code = $1
	`, code).Scan(&userID, &expiresAt)
	if err == sql.ErrNoRows {
		return "", "", errors.New("invalid or expired verification link")
	}
	if err != nil {
		return "", "", err
	}

	if time.Now().After(expiresAt) {
		_ = s.deleteVerification(userID)
		return "", "", errors.New("verification link has expired")
	}

	_, err = s.DB.Exec(`UPDATE users SET is_verified = true WHERE id = $1`, userID)
	if err != nil {
		return "", "", err
	}

	_ = s.deleteVerification(userID)

	return utils.GenerateTokens(userID)
}

// 
func (s *AuthService) ResendVerification(email string) error {
	var userID int64
	var isVerified bool

	err := s.DB.QueryRow(`
		SELECT id, is_verified FROM users WHERE email = $1
	`, email).Scan(&userID, &isVerified)

	if err == sql.ErrNoRows {
		// Чтобы не раскрывать, есть ли пользователь, возвращаем nil
		return nil
	}
	if err != nil {
		return err
	}

	if isVerified {
		return nil // Уже подтверждён
	}

	// Удаляем старый код
	_ = s.deleteVerification(userID)

	// Создаём новый
	code := generateVerificationCode()
	expiration := time.Now().Add(10 * time.Minute)

	_, err = s.DB.Exec(`
		INSERT INTO email_verifications (user_id, code, expires_at)
		VALUES ($1, $2, $3)
	`, userID, code, expiration)
	if err != nil {
		return err
	}

	go utils.SendVerificationEmail(email, code)
	return nil
}


func (s *AuthService) deleteVerification(userID int64) error {
	_, err := s.DB.Exec(`DELETE FROM email_verifications WHERE user_id = $1`, userID)
	return err
}

func generateVerificationCode() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *AuthService) RequestPasswordReset(email string) error {
	var userID int64
	err := s.DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err == sql.ErrNoRows {
		return errors.New("user not found")
	} else if err != nil {
		return err
	}

	token := uuid.New().String()
	expires := time.Now().Add(15 * time.Minute)

	_, err = s.DB.Exec(`
		INSERT INTO password_resets (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, userID, token, expires)
	if err != nil {
		return err
	}

	resetLink := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", token)
	go utils.SendResetPasswordEmail(email, resetLink)

	return nil
}

// сброс пароля

func (s *AuthService) ResetPassword(token, newPassword string) error {
	var userID int64
	var expires time.Time
	var used bool

	err := s.DB.QueryRow(`
		SELECT user_id, expires_at, used FROM password_resets WHERE token = $1
	`, token).Scan(&userID, &expires, &used)
	if err == sql.ErrNoRows {
		return errors.New("invalid or expired token")
	} else if err != nil {
		return err
	}

	if used || time.Now().After(expires) {
		return errors.New("token expired or already used")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users SET password_hash = $1 WHERE id = $2", hashed, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE password_resets SET used = true WHERE token = $1", token)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
