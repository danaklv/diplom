package models

type EcoAnswerRequest struct {
	Answers map[int]int `json:"answers"` // question_id → value (0–5)
}

type EcoQuestion struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Question string `json:"question"`
	MaxValue int    `json:"max_value"`
}



type EcoResult struct {
	UserID      int64  `json:"user_id"`
	TotalScore  int    `json:"total_score"`
	Category    string `json:"category"`
	Description string `json:"description"`
}
