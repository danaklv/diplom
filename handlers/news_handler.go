package handlers

import (
	"dl/services"
	"net/http"
)

type NewsHandler struct {
	Service *services.NewsService
}

func NewNewsHandler(service *services.NewsService) *NewsHandler {
	return &NewsHandler{Service: service}
}

// ------------------------ GET ALL NEWS ------------------------

func (h *NewsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	news, err := h.Service.GetAllNews()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, news)
}
