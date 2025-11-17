package handlers

import (
	"database/sql"
	"dl/models"
	"dl/services"
	"dl/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ProfileHandler struct {
	Service *services.ProfileService
}

func nullStr(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

// ------------------------ GET PROFILE ------------------------

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, err := utils.UserIDFromContext(r.Context())
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.Service.GetProfile(userID)
	if err != nil {
		jsonError(w, http.StatusNotFound, err.Error())
		return
	}

	resp := models.ProfileResponse{
		ID:        profile.ID,
		Username:  profile.Username,
		Email:     profile.Email,
		FirstName: nullStr(profile.FirstName),
		LastName:  nullStr(profile.LastName),
		Gender:    nullStr(profile.Gender),
		BirthDate: nullStr(profile.BirthDate),
		Bio:       nullStr(profile.Bio),
		Avatar:    nullStr(profile.ProfilePicture),
		Rating:    profile.Rating,
		Level:     profile.Level,  // если добавишь в модель
		League:    profile.League, // если добавишь в модель
		UpdatedAt: profile.UpdatedAt,
	}

	jsonResponse(w, http.StatusOK, resp)
}

// ------------------------ UPDATE PROFILE ------------------------

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var data struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Gender    string `json:"gender"`
		Bio       string `json:"bio"`
		BirthDate string `json:"birth_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	userID, err := utils.UserIDFromContext(r.Context())
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.Service.UpdateProfile(userID, data.FirstName, data.LastName, data.Gender, data.Bio, data.BirthDate); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"message": "profile updated successfully",
	})
}

// ------------------------ DELETE PROFILE ------------------------

func (h *ProfileHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, err := utils.UserIDFromContext(r.Context())
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.Service.DeleteProfile(userID); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "profile deleted"})
}

// ------------------------ UPLOAD AVATAR ------------------------

func (h *ProfileHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 10MB limit
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		jsonError(w, http.StatusBadRequest, "could not parse form")
		return
	}

	file, handler, err := r.FormFile("avatar")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "could not read file")
		return
	}
	defer file.Close()

	// validate image type + size
	if err := utils.ValidateImage(file, handler); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := utils.UserIDFromContext(r.Context())
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	filename := fmt.Sprintf("user_%d_%s", userID, handler.Filename)
	filePath := fmt.Sprintf("./uploads/users/%s", filename)

	dst, err := os.Create(filePath)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "could not save file")
		return
	}
	defer dst.Close()

	// Reset file reader before copying (ValidateImage moved pointer)
	file.Seek(0, 0)

	if _, err := io.Copy(dst, file); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to save image")
		return
	}

	publicPath := "/uploads/users/" + filename

	if err := h.Service.UpdateProfilePicture(userID, publicPath); err != nil {
		jsonError(w, http.StatusInternalServerError, "database update failed")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "avatar uploaded successfully",
		"file":    publicPath,
	})
}
