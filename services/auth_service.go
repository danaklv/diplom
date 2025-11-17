package services

import (
	"dl/repositories"
	"dl/utils"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

// --------------------------------------------------------
// REGISTER
// --------------------------------------------------------

func (s *AuthService) Register(username, email, password string) (string, string, error) {
	// trim
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	// validation
	if err := utils.ValidateEmail(email); err != nil {
		return "", "", err
	}
	if err := utils.ValidateUsername(username); err != nil {
		return "", "", err
	}
	if err := utils.ValidatePassword(password); err != nil {
		return "", "", err
	}

	// hash
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	// try to create user
	userID, err := s.Repo.CreateUser(username, email, hashed)

	// unique errors
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

	// generate verification code
	code := utils.GenerateVerificationCode()
	expires := time.Now().Add(10 * time.Minute)

	// store code
	if err := s.Repo.StoreVerificationCode(userID, code, expires); err != nil {
		return "", "", err
	}

	// send email
	go utils.SendVerificationEmail(email, code)

	// do not return tokens until email is verified
	return "", "", nil
}

// --------------------------------------------------------
// LOGIN
// --------------------------------------------------------

func (s *AuthService) Login(email, password string) (string, string, error) {
	email = strings.TrimSpace(email)

	userID, hashed, verified, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	if bcrypt.CompareHashAndPassword(hashed, []byte(password)) != nil {
		return "", "", errors.New("invalid email or password")
	}

	if !verified {
		return "", "", errors.New("email not verified")
	}

	return utils.GenerateTokens(userID)
}

// --------------------------------------------------------
// VERIFY EMAIL
// --------------------------------------------------------

func (s *AuthService) VerifyEmail(code string) (string, string, error) {
	// userID, expires, err := s.Repo.GetUserByVerificationCode(code)
	userID, _, err := s.Repo.GetUserByVerificationCode(code)
	if err != nil {
		return "", "", errors.New("invalid or expired verification link")
	}

	// if time.Now().After(expires) {
	// 	_ = s.Repo.DeleteVerificationCode(userID)
	// 	return "", "", errors.New("verification link expired")
	// }

	// if err := s.Repo.SetUserVerified(userID); err != nil {
	// 	return "", "", err
	// }

	// _ = s.Repo.DeleteVerificationCode(userID)

	return utils.GenerateTokens(userID)
}

// --------------------------------------------------------
// REQUEST PASSWORD RESET
// --------------------------------------------------------

func (s *AuthService) RequestPasswordReset(email string) error {
	userID, _, _, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	token := uuid.New().String()
	expires := time.Now().Add(15 * time.Minute)

	if err := s.Repo.CreatePasswordReset(userID, token, expires); err != nil {
		return err
	}

	resetLink := "http://localhost:5173/reset-password?token=" + token
	go utils.SendResetPasswordEmail(email, resetLink)

	return nil
}

// --------------------------------------------------------
// RESET PASSWORD
// --------------------------------------------------------

func (s *AuthService) ResetPassword(token, newPassword string) error {
	userID, expires, used, err := s.Repo.GetPasswordReset(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if used || time.Now().After(expires) {
		return errors.New("token expired or already used")
	}

	// validate new password
	if err := utils.ValidatePassword(newPassword); err != nil {
		return err
	}

	// hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// update password
	if err := s.Repo.UpdateUserPassword(userID, hashed); err != nil {
		return err
	}

	// mark token as used
	return s.Repo.MarkResetTokenUsed(token)
}
