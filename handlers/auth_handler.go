package handlers

import (
	"dl/services"
	"encoding/json"
	"net/http"
	"strings"
)

// Универсальный JSON ответ
func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// Универсальная ошибка
func jsonError(w http.ResponseWriter, status int, msg string) {
	jsonResponse(w, status, map[string]string{"error": msg})
}

type AuthHandler struct {
	Service *services.AuthService
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ------------------------ REGISTER ------------------------

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		jsonError(w, http.StatusBadRequest, "username, email and password are required")
		return
	}

	access, refresh, err := h.Service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	// После регистрации email не всегда подтверждён.
	// Если сервис вернул токены — отдаём. Если нет — просто success.
	if access == "" && refresh == "" {
		jsonResponse(w, http.StatusOK, map[string]string{
			"message": "User registered. Please check your email to verify your account.",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// ------------------------ LOGIN ------------------------

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" || req.Password == "" {
		jsonError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	access, refresh, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// ------------------------ EMAIL VERIFICATION ------------------------

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		jsonError(w, http.StatusBadRequest, "missing verification code")
		return
	}

	access, refresh, err := h.Service.VerifyEmail(code)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"message":       "Email verified successfully",
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// ------------------------ FORGOT PASSWORD ------------------------

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var data struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	data.Email = strings.TrimSpace(data.Email)

	if data.Email == "" {
		jsonError(w, http.StatusBadRequest, "email is required")
		return
	}

	if err := h.Service.RequestPasswordReset(data.Email); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"message": "If this email exists, a reset link has been sent",
	})
}

// ------------------------ RESET PASSWORD ------------------------

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var data struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if data.Token == "" || data.NewPassword == "" {
		jsonError(w, http.StatusBadRequest, "token and new_password are required")
		return
	}

	if err := h.Service.ResetPassword(data.Token, data.NewPassword); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"message": "Password successfully reset",
	})
}
