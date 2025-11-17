package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"dl/handlers"
	"dl/middleware"
	"dl/repositories"
	"dl/seeders"
	"dl/services"

	_ "github.com/lib/pq"
)

func main() {
	// --- Конфигурация (env с fallback) ---
	dbURL := getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ecofoot?sslmode=disable")
	addr := getenv("HTTP_ADDR", ":8080")
	uploadsDir := getenv("UPLOADS_DIR", "./uploads")
	newsIntervalMin := getenvInt("NEWS_INTERVAL_MIN", 30)

	// --- Создать папку uploads если нет ---
	if err := ensureDir(uploadsDir); err != nil {
		log.Fatalf("failed to ensure uploads dir: %v", err)
	}

	// --- DB init ---
	db := InitDB(dbURL)
	defer db.Close()

	if err := seeders.RunAllSeeders(db); err != nil {
		log.Fatal("Failed to run seeders: ", err)
	}

	// --- AUTH ---
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := &handlers.AuthHandler{Service: authService}

	// --- PROFILE ---
	profileRepo := repositories.NewProfileRepository(db)
	profileService := services.NewProfileService(profileRepo)
	profileHandler := &handlers.ProfileHandler{Service: profileService}

	// --- RATING ---
	ratingRepo := repositories.NewRatingRepository(db)
	ratingService := services.NewRatingService(ratingRepo)
	ratingHandler := &handlers.RatingHandler{Service: ratingService}

	// --- NEWS ---
	newsRepo := repositories.NewNewsRepository(db)
	newsService := services.NewNewsService(newsRepo)
	newsHandler := handlers.NewNewsHandler(newsService)

	// -- ECO
	ecoRepo := repositories.NewEcoRepository(db)
	ecoService := services.NewEcoService(ecoRepo)
	ecoHandler := handlers.EcoHandler{Service: ecoService}

	// --- Router ---
	mux := http.NewServeMux()

	// Public auth routes
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/verify", authHandler.Verify)
	mux.HandleFunc("/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("/reset-password", authHandler.ResetPassword)
	// TODO: add /refresh, /logout endpoints in AuthHandler (and implement refresh token storage)

	// Static uploads
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	// Protected profile routes (JWTAuth wrapper uses current signature: middleware.JWTAuth(next http.HandlerFunc) http.HandlerFunc)
	mux.Handle("/eco", middleware.JWTAuth(http.HandlerFunc(ecoHandler.GetQuestions)))
	mux.Handle("/profile", middleware.JWTAuth(http.HandlerFunc(profileHandler.GetProfile)))
	mux.Handle("/update-profile", middleware.JWTAuth(http.HandlerFunc(profileHandler.UpdateProfile)))
	mux.Handle("/delete-profile", middleware.JWTAuth(http.HandlerFunc(profileHandler.DeleteProfile)))
	mux.Handle("/upload-avatar", middleware.JWTAuth(http.HandlerFunc(profileHandler.UploadAvatar)))

	mux.Handle("/add-action", middleware.JWTAuth(http.HandlerFunc(ratingHandler.AddAction)))
	mux.Handle("/user-actions", middleware.JWTAuth(http.HandlerFunc(ratingHandler.GetUserActions)))
	mux.Handle("/leaderboard", middleware.JWTAuth(http.HandlerFunc(ratingHandler.GetLeaderboard)))

	// News (public)
	mux.HandleFunc("/news", newsHandler.GetAll)

	// Middleware chain: CORS -> (optionally Logging/Recovery) -> mux
	handler := middleware.EnableCORS(mux)
	// TODO: add middleware.Recovery(handler) and middleware.RequestLogger(handler) if добавите реализации

	// --- Background job: обновление новостей по расписанию ---
	go func() {
		ticker := time.NewTicker(time.Duration(newsIntervalMin) * time.Minute)
		defer ticker.Stop()

		// Запуск сразу при старте
		if err := newsService.UpdateNews(); err != nil {
			log.Println("news update error:", err)
		}

		for range ticker.C {
			if err := newsService.UpdateNews(); err != nil {
				log.Println("news update error:", err)
			}
		}
	}()

	// --- HTTP Server с таймаутами и graceful shutdown ---
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Server listening on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Graceful shutdown при SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited properly")
}

// InitDB теперь принимает строку подключения
func InitDB(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	return db
}

// Помощники
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if t, err := strconv.Atoi(v); err == nil {
			return t
		}
	}
	return fallback
}

func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	// защитим от относительных путей и вернём абсолютный путь (по желанию)
	_, err := filepath.Abs(path)
	return err
}
