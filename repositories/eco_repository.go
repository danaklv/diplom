package repositories

import (
	"database/sql"
	"dl/models"
)

type EcoRepository struct {
	DB *sql.DB
}

func NewEcoRepository(db *sql.DB) *EcoRepository {
	return &EcoRepository{DB: db}
}

//
// ---------------------------------------------------------------
// QUESTIONS
// ---------------------------------------------------------------
//

func (r *EcoRepository) GetQuestions() ([]models.EcoQuestion, error) {
	rows, err := r.DB.Query(`
        SELECT id, category, question, max_value
        FROM eco_questions
        ORDER BY id ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.EcoQuestion

	for rows.Next() {
		var q models.EcoQuestion
		if err := rows.Scan(&q.ID, &q.Category, &q.Question, &q.MaxValue); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, rows.Err()
}

//
// ---------------------------------------------------------------
// SAVE ANSWERS
// ---------------------------------------------------------------
//

func (r *EcoRepository) SaveAnswers(userID int64, answers map[int]int) error {
	for qID, value := range answers {
		_, err := r.DB.Exec(`
            INSERT INTO eco_answers (user_id, question_id, value)
            VALUES ($1, $2, $3)
        `, userID, qID, value)

		if err != nil {
			return err
		}
	}
	return nil
}

//
// ---------------------------------------------------------------
// SAVE RESULT
// ---------------------------------------------------------------
//

func (r *EcoRepository) SaveResult(userID int64, total int, category, description string) error {
	_, err := r.DB.Exec(`
        INSERT INTO eco_results (user_id, total_score, category, description)
        VALUES ($1, $2, $3, $4)
    `, userID, total, category, description)
	return err
}

//
// ---------------------------------------------------------------
// GET LATEST RESULT
// ---------------------------------------------------------------
//

func (r *EcoRepository) GetLatestResult(userID int64) (*models.EcoResult, error) {
	var result models.EcoResult

	err := r.DB.QueryRow(`
        SELECT total_score, category, description
        FROM eco_results
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT 1
    `, userID).Scan(
		&result.TotalScore,
		&result.Category,
		&result.Description,
	)

	if err != nil {
		return nil, err
	}

	result.UserID = userID
	return &result, nil
}
