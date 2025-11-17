package handlers

import (
    "dl/services"
    "dl/utils"
    "encoding/json"
    "net/http"
)

type RatingHandler struct {
    Service *services.RatingService
}



// ------------------------ ADD ECO ACTION ------------------------

func (h *RatingHandler) AddAction(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    var data struct {
        ActionID int64 `json:"action_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        jsonError(w, http.StatusBadRequest, "invalid JSON body")
        return
    }

    if data.ActionID == 0 {
        jsonError(w, http.StatusBadRequest, "action_id is required")
        return
    }

    userID, err := utils.UserIDFromContext(r.Context())
    if err != nil {
        jsonError(w, http.StatusUnauthorized, "unauthorized")
        return
    }

    if err := h.Service.AddEcoAction(userID, data.ActionID); err != nil {
        jsonError(w, http.StatusBadRequest, err.Error())
        return
    }

    jsonResponse(w, http.StatusOK, map[string]string{
        "message": "action recorded successfully",
    })
}

// ------------------------ GET USER ACTION HISTORY ------------------------

func (h *RatingHandler) GetUserActions(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    userID, err := utils.UserIDFromContext(r.Context())
    if err != nil {
        jsonError(w, http.StatusUnauthorized, "unauthorized")
        return
    }

    actions, err := h.Service.GetUserActions(userID)
    if err != nil {
        jsonError(w, http.StatusInternalServerError, err.Error())
        return
    }

    jsonResponse(w, http.StatusOK, actions)
}

// ------------------------ GET LEADERBOARD ------------------------

func (h *RatingHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    leaderboard, err := h.Service.GetLeaderboard(10)
    if err != nil {
        jsonError(w, http.StatusInternalServerError, err.Error())
        return
    }

    jsonResponse(w, http.StatusOK, leaderboard)
}
