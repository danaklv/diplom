package handlers

import (
	"database/sql"
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

// Получение профиля
// func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {

// 	userID, err := utils.UserIDFromContext(r.Context())

// 	fmt.Println("USERID = ", userID)
// 	if err != nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	profile, err := h.Service.GetProfile(userID)
// 	fmt.Println(profile, err)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(profile)
// }

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := h.Service.GetProfile(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// ✅ Преобразуем NullString → string
	resp := map[string]interface{}{
		"id":         profile.ID,
		"username":   profile.Username,
		"email":      profile.Email,
		"first_name": nullToString(profile.FirstName),
		"last_name":  nullToString(profile.LastName),
		"gender":     nullToString(profile.Gender),
		"birth_date": nullToString(profile.BirthDate),
		"bio":        nullToString(profile.Bio),
		"avatar":     nullToString(profile.ProfilePicture),
		"rating":     profile.Rating,
		"updated_at": profile.UpdatedAt,
	}

	json.NewEncoder(w).Encode(resp)
}

func nullToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// Обновление профиля
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var data struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Gender    string `json:"gender"`
		Bio       string `json:"bio"`
		BirthDate string `json:"birth_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := utils.UserIDFromContext(r.Context())
	fmt.Println("USERID = ", userID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.Service.UpdateProfile(userID, data.FirstName, data.LastName, data.Gender, data.Bio, data.BirthDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

// Удаление профиля
func (h *ProfileHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID := int64(1) // TODO: из JWT
	err := h.Service.DeleteProfile(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile deleted"})
}

func (h *ProfileHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // до 10 МБ
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Could not read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// ✅ Проверяем тип и размер файла
	if err := utils.ValidateImage(file, handler); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: заменить на ID из JWT
	userID := int64(1)

	// Формируем путь
	filename := fmt.Sprintf("user_%d_%s", userID, handler.Filename)
	filePath := fmt.Sprintf("./uploads/users/%s", filename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Копируем содержимое файла
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	publicPath := fmt.Sprintf("/uploads/users/%s", filename)
	if err := h.Service.UpdateProfilePicture(userID, publicPath); err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Avatar uploaded successfully",
		"file":    publicPath,
	})
}
