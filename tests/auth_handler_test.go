package tests

import (
	"bytes"
	"dl/handlers"
	"dl/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRegisterHandler(t *testing.T) {
	// мок бд
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("UserTest", "testemail@gmail.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("INSERT INTO email_verifications").
		WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	authService := &services.AuthService{DB: db}
	handler := &handlers.AuthHandler{Service: authService}

	body := map[string]string{
		"username": "UserTest",
		"email":    "testemail@gmail.com",
		"password": "StrongPass123!",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Register(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rr.Code)
	}
}
