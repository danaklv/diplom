package handlers

import (
	"dl/models"
	"dl/services"
	"dl/utils"
	"encoding/json"
	"net/http"
)

type EcoHandler struct {
	Service *services.EcoService
}

// ------------------------ SUBMIT ANSWERS ------------------------

func (h *EcoHandler) Submit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Проверка токена
	userID, err := utils.UserIDFromContext(r.Context())
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Парсим JSON
	var req models.EcoAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if len(req.Answers) == 0 {
		jsonError(w, http.StatusBadRequest, "answers are required")
		return
	}

	// Передаём в сервис
	result, err := h.Service.SubmitAnswers(userID, req.Answers)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, result)
}

// ------------------------ GET LATEST RESULT ------------------------

func (h *EcoHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, err := utils.UserIDFromContext(r.Context())
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.Service.GetLatest(userID)
	if err != nil {
		jsonError(w, http.StatusNotFound, "no results found")
		return
	}

	jsonResponse(w, http.StatusOK, result)
}

// ------------------------ GET QUESTIONS (опционально) ------------------------

func (h *EcoHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	questions, err := h.Service.GetQuestions()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, questions)
}

