package handlers

import (
	"dl/services"
	"encoding/json"
	"fmt"
	"net/http"
)

type NewsHandler struct {
	Service *services.NewsService
}

func NewNewsHandler(service *services.NewsService) *NewsHandler {
	return &NewsHandler{Service: service}
}

// news
func (h *NewsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	news, err := h.Service.GetAllNews()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("news1 = ", news)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}
