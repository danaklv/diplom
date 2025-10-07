package services

import (
	"database/sql"
	"dl/models"
	"fmt"
	"time"
)

type RatingService struct {
	DB *sql.DB
}

func (s *RatingService) AddEcoAction(userID, actionID int64) error {
	var points int
	err := s.DB.QueryRow("SELECT points FROM eco_actions WHERE id = $1", actionID).Scan(&points)
	if err == sql.ErrNoRows {
		return fmt.Errorf("action not found")
	} else if err != nil {
		return err
	}

	_, err = s.DB.Exec(`
		INSERT INTO user_actions (user_id, action_id, points, created_at)
		VALUES ($1, $2, $3, $4)
	`, userID, actionID, points, time.Now())
	if err != nil {
		return err
	}

	// Старый уровень до обновления
	var oldLevel int
	s.DB.QueryRow("SELECT level FROM users WHERE id = $1", userID).Scan(&oldLevel)

	_, err = s.DB.Exec("UPDATE users SET rating = rating + $1 WHERE id = $2", points, userID)
	if err != nil {
		return err
	}

	// Обновляем уровень
	if err := s.UpdateUserLevel(userID); err != nil {
		return err
	}

	// Проверяем новый уровень
	var newLevel int
	s.DB.QueryRow("SELECT level FROM users WHERE id = $1", userID).Scan(&newLevel)

	if newLevel > oldLevel {
		msg := fmt.Sprintf("🎉 Congratulations! You've reached Level %d!", newLevel)
		_, _ = s.DB.Exec("INSERT INTO notifications (user_id, message) VALUES ($1, $2)", userID, msg)
	}

	return nil
}

// История действий пользователя
func (s *RatingService) GetUserActions(userID int64) ([]models.UserAction, error) {
	rows, err := s.DB.Query(`
		SELECT a.name, ua.points, ua.created_at
		FROM user_actions ua
		JOIN eco_actions a ON ua.action_id = a.id
		WHERE ua.user_id = $1
		ORDER BY ua.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []models.UserAction
	for rows.Next() {
		var a models.UserAction
		if err := rows.Scan(&a.ActionName, &a.Points, &a.CreatedAt); err == nil {
			actions = append(actions, a)
		}
	}
	return actions, nil
}

// Таблица лидеров
func (s *RatingService) GetLeaderboard(limit int) ([]models.LeaderboardEntry, error) {
	rows, err := s.DB.Query(`
		SELECT username, rating, level, league
		FROM users
		ORDER BY rating DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []models.LeaderboardEntry
	for rows.Next() {
		var entry models.LeaderboardEntry
		if err := rows.Scan(&entry.Username, &entry.Rating, &entry.Level, &entry.League); err == nil {
			leaderboard = append(leaderboard, entry)
		}
	}
	return leaderboard, nil
}

// Определяем уровень и лигу на основе рейтинга
func (s *RatingService) UpdateUserLevel(userID int64) error {
	var rating int
	err := s.DB.QueryRow("SELECT rating FROM users WHERE id = $1", userID).Scan(&rating)
	if err != nil {
		return err
	}

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

	_, err = s.DB.Exec("UPDATE users SET level = $1, league = $2 WHERE id = $3", level, league, userID)
	return err
}
