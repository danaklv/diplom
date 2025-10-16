package tests

import (
	"bytes"
	"database/sql"
	"dl/handlers"
	"dl/services"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	connStr := "postgres://postgres:dana1234@localhost:5432/ecofoot?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect test DB: %v", err)
	}
	db.Exec("TRUNCATE users, email_verifications RESTART IDENTITY CASCADE;")
	return db
}

func TestRegisterIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := &services.AuthService{DB: db}
	handler := &handlers.AuthHandler{Service: service}

	body := map[string]string{
		"username": "TestUser",
		"email":    "testuser@gmail.com",
		"password": "StrongPass123!",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Register(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, hot %d", rr.Code)
	}

	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", "testuser@gmail.com").Scan(&count)
	if err != nil || count == 0 {
		fmt.Println("errr=", err)
		t.Errorf("user not inserted in DB")
	}

}
