package handlers

import (
	"dl/services"
	"encoding/json"
	"net/http"
)

type RatingHandler struct {
	Service *services.RatingService
}

func (h *RatingHandler) AddAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		ActionID int64 `json:"action_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// TODO: userID из JWT
	userID := int64(1)

	if err := h.Service.AddEcoAction(userID, data.ActionID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Action recorded successfully"})
}
func (h *RatingHandler) GetUserActions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: userID из JWT
	userID := int64(1)

	actions, err := h.Service.GetUserActions(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(actions)
}

func (h *RatingHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	leaderboard, err := h.Service.GetLeaderboard(10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(leaderboard)
}
