package main

import (
	"database/sql"
	"dl/handlers"
	"dl/middleware"
	"dl/services"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// func main() {
// 	DB := InitDB()

// 	authService := services.NewAuthService(DB)
// 	authHandler := &handlers.AuthHandler{Service: authService}
// 	profileService := &services.ProfileService{DB: DB}
// 	profileHandler := &handlers.ProfileHandler{Service: profileService}

// 	http.HandleFunc("/register", authHandler.Register)
// 	http.HandleFunc("/login", authHandler.Login)
// 	http.HandleFunc("/verify", authHandler.Verify)
// 	http.HandleFunc("/forgot-password", authHandler.ForgotPassword)
// 	http.HandleFunc("/reset-password", authHandler.ResetPassword)

// 	http.HandleFunc("/profile", profileHandler.GetProfile)           //
// 	http.HandleFunc("/update-profile", profileHandler.UpdateProfile) // PUT
// 	http.HandleFunc("/delete-profile", profileHandler.DeleteProfile)

// 	http.HandleFunc("/upload-avatar", profileHandler.UploadAvatar)
// 	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

// 	ratingService := &services.RatingService{DB: DB}
// 	ratingHandler := &handlers.RatingHandler{Service: ratingService}

// 	http.HandleFunc("/add-action", ratingHandler.AddAction)
// 	http.HandleFunc("/user-actions", ratingHandler.GetUserActions)
// 	http.HandleFunc("/leaderboard", ratingHandler.GetLeaderboard)

// 	http.ListenAndServe(":8080", nil)
// }

func main() {
	DB := InitDB()

	authService := services.NewAuthService(DB)
	authHandler := &handlers.AuthHandler{Service: authService}
	profileService := &services.ProfileService{DB: DB}
	profileHandler := &handlers.ProfileHandler{Service: profileService}
	ratingService := &services.RatingService{DB: DB}
	ratingHandler := &handlers.RatingHandler{Service: ratingService}

	// ✅ создаём новый маршрутизатор (mux)
	mux := http.NewServeMux()

	// ---------- Auth ----------
	// mux.HandleFunc("/register", authHandler.Register)
	// mux.HandleFunc("/login", authHandler.Login)
	// mux.HandleFunc("/verify", authHandler.Verify)
	// mux.HandleFunc("/forgot-password", authHandler.ForgotPassword)
	// mux.HandleFunc("/reset-password", authHandler.ResetPassword)

	// // ---------- Profile ----------
	// mux.HandleFunc("/profile", profileHandler.GetProfile)
	// mux.HandleFunc("/update-profile", profileHandler.UpdateProfile)
	// mux.HandleFunc("/delete-profile", profileHandler.DeleteProfile)
	// mux.HandleFunc("/upload-avatar", profileHandler.UploadAvatar)

	// отдача файлов из папки uploads/
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// ---------- Rating ----------
	// mux.HandleFunc("/add-action", ratingHandler.AddAction)
	// mux.HandleFunc("/user-actions", ratingHandler.GetUserActions)
	// mux.HandleFunc("/leaderboard", ratingHandler.GetLeaderboard)

	// // ✅ Оборачиваем всё в CORS middleware
	// handler := middleware.EnableCORS(mux)

	// // ✅ Запускаем сервер
	// http.ListenAndServe(":8080", handler)
	// mux := http.NewServeMux()

	// Публичные маршруты
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/verify", authHandler.Verify)
	mux.HandleFunc("/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("/reset-password", authHandler.ResetPassword)

	// ✅ Защищённые маршруты (нужен токен)
	mux.HandleFunc("/profile", middleware.JWTAuth(profileHandler.GetProfile))
	mux.HandleFunc("/update-profile", middleware.JWTAuth(profileHandler.UpdateProfile))
	mux.HandleFunc("/delete-profile", middleware.JWTAuth(profileHandler.DeleteProfile))
	mux.HandleFunc("/upload-avatar", middleware.JWTAuth(profileHandler.UploadAvatar))

	mux.HandleFunc("/add-action", middleware.JWTAuth(ratingHandler.AddAction))
	mux.HandleFunc("/user-actions", middleware.JWTAuth(ratingHandler.GetUserActions))
	mux.HandleFunc("/leaderboard", middleware.JWTAuth(ratingHandler.GetLeaderboard))

	handler := middleware.EnableCORS(mux)
	http.ListenAndServe(":8080", handler)

}

func InitDB() *sql.DB {
	connStr := "postgres://postgres:dana1234@localhost:5432/ecofoot?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database ping error:", err)
	}

	return db
}
