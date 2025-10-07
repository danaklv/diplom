package models

import "time"

type User struct {
	ID       int64
	Username string
	Email    string
	Password []byte
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
}
