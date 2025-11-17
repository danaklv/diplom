package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"-"`
	IsVerified   bool   `json:"is_verified"`
}

type Session struct {
	UserID     int64
	Cookie     string
	Expiration time.Time
}

type UserAction struct {
	ActionName string    `json:"action_name"`
	Points     int       `json:"points"`
	CreatedAt  time.Time `json:"created_at"`
}

type LeaderboardEntry struct {
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Level    int    `json:"level"`
	League   string `json:"league"`
	Avatar   string `json:"avatar"`
}

type UserProfile struct {
	ID             int64          `json:"id"`
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	FirstName      sql.NullString `json:"first_name"`
	LastName       sql.NullString `json:"last_name"`
	Gender         sql.NullString `json:"gender"`
	BirthDate      sql.NullString `json:"birth_date"`
	Bio            sql.NullString `json:"bio"`
	ProfilePicture sql.NullString `json:"profile_picture"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Rating         int            `json:"rating"`
	Level          int            `json:"level"`
	League         string         `json:"league"`
}

type UserStats struct {
	UserID int64  `json:"user_id"`
	Rating int    `json:"rating"`
	Level  int    `json:"level"`
	League string `json:"league"`
}

type ProfileResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Bio       string `json:"bio"`

	Avatar    string    `json:"avatar"`
	Rating    int       `json:"rating"`
	Level     int       `json:"level"`
	League    string    `json:"league"`
	UpdatedAt time.Time `json:"updated_at"`
}
