package repositories

import (
    "database/sql"
    "dl/models"
    "errors"
    "time"
)

type RatingRepository struct {
    DB *sql.DB
}

func NewRatingRepository(db *sql.DB) *RatingRepository {
    return &RatingRepository{DB: db}
}

// ------------------------ GET ACTION POINTS ------------------------

func (r *RatingRepository) GetActionPoints(actionID int64) (int, error) {
    var points int
    err := r.DB.QueryRow(`SELECT points FROM eco_actions WHERE id = $1`, actionID).Scan(&points)
    if err == sql.ErrNoRows {
        return 0, errors.New("action not found")
    }
    return points, err
}

// ------------------------ ADD USER ACTION ------------------------

func (r *RatingRepository) AddUserAction(userID, actionID int64, points int) error {
    _, err := r.DB.Exec(`
        INSERT INTO user_actions (user_id, action_id, points, created_at)
        VALUES ($1, $2, $3, $4)
    `, userID, actionID, points, time.Now())
    return err
}

// ------------------------ UPDATE RATING ------------------------

func (r *RatingRepository) UpdateRating(userID int64, points int) error {
    _, err := r.DB.Exec(`UPDATE users SET rating = rating + $1 WHERE id = $2`, points, userID)
    return err
}

// ------------------------ GET LEVEL ------------------------

func (r *RatingRepository) GetUserLevel(userID int64) (int, error) {
    var level int
    err := r.DB.QueryRow(`SELECT level FROM users WHERE id = $1`, userID).Scan(&level)
    return level, err
}

// ------------------------ UPDATE LEVEL ------------------------

func (r *RatingRepository) UpdateLevel(userID int64, level int, league string) error {
    _, err := r.DB.Exec(`
        UPDATE users SET level = $1, league = $2 WHERE id = $3
    `, level, league, userID)
    return err
}

// ------------------------ GET USER ACTIONS ------------------------

func (r *RatingRepository) GetUserActions(userID int64) ([]models.UserAction, error) {
    rows, err := r.DB.Query(`
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
        if err := rows.Scan(&a.ActionName, &a.Points, &a.CreatedAt); err != nil {
            return nil, err
        }
        actions = append(actions, a)
    }
    return actions, rows.Err()
}

// ------------------------ LEADERBOARD ------------------------

func (r *RatingRepository) GetLeaderboard(limit int) ([]models.LeaderboardEntry, error) {
    rows, err := r.DB.Query(`
        SELECT username, rating, level, league
        FROM users ORDER BY rating DESC LIMIT $1
    `, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []models.LeaderboardEntry
    for rows.Next() {
        var e models.LeaderboardEntry
        if err := rows.Scan(&e.Username, &e.Rating, &e.Level, &e.League); err != nil {
            return nil, err
        }
        list = append(list, e)
    }
    return list, rows.Err()
}
