package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shaolin/backend/internal/api"
	"shaolin/backend/internal/repository"
	"shaolin/backend/internal/repository/sqlite"
	"shaolin/backend/migrations"
)

func main() {
	log.Println("Starting Shaolin Academy API Server...")

	// 1. Database Configuration
	driver := getEnv("DB_DRIVER", "sqlite")
	dsn := getEnv("DB_DSN", "shaolin.db")

	log.Printf("Connecting to database: driver=%s, dsn=%s", driver, dsn)
	db, err := repository.Connect(repository.DBConfig{
		Driver:          driver,
		DSN:             dsn,
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		log.Println("Closing database connections...")
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// 2. Auto-run Schema Migrations
	log.Println("Running database migrations...")
	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("Migrations applied successfully.")

	// 3. Initialize Router & Middleware
	userRepo := sqlite.NewSQLiteUserRepository(db)
	termRepo := sqlite.NewSQLiteTermRepository(db)
	bookingRepo := sqlite.NewSQLiteBookingRepository(db)
	tokenRepo := sqlite.NewSQLiteTokenRepository(db)
	storeRepo := sqlite.NewSQLiteStoreRepository(db)
	gradingRepo := sqlite.NewSQLiteGradingRepository(db)
	forumRepo := sqlite.NewSQLiteForumRepository(db)
	mediaRepo := sqlite.NewSQLiteMediaRepository(db)
	guideRepo := sqlite.NewSQLiteGuideRepository(db)

	authHandler := api.NewAuthHandler(userRepo)
	termHandler := api.NewTermHandler(termRepo)
	classHandler := api.NewClassHandler(bookingRepo, userRepo)
	userHandler := api.NewUserHandler(userRepo)
	bookingHandler := api.NewBookingHandler(bookingRepo, tokenRepo)
	storeHandler := api.NewStoreHandler(storeRepo)
	gradingHandler := api.NewGradingHandler(gradingRepo)
	forumsMediaHandler := api.NewForumsMediaHandler(forumRepo, mediaRepo, guideRepo, gradingRepo)

	mux := http.NewServeMux()

	// Store & Cart Endpoints
	mux.HandleFunc("GET /api/store/products", storeHandler.ListProducts)
	mux.Handle("GET /api/cart", api.AuthMiddleware(http.HandlerFunc(storeHandler.GetCart)))
	mux.Handle("POST /api/cart", api.AuthMiddleware(http.HandlerFunc(storeHandler.AddToCart)))
	mux.Handle("DELETE /api/cart/{id}", api.AuthMiddleware(http.HandlerFunc(storeHandler.RemoveFromCart)))
	mux.Handle("POST /api/store/checkout", api.AuthMiddleware(http.HandlerFunc(storeHandler.Checkout)))
	mux.Handle("GET /api/store/orders/history", api.AuthMiddleware(http.HandlerFunc(storeHandler.GetOrderHistory)))

	// Token & Transactions Endpoints
	mux.Handle("GET /api/tokens/balance", api.AuthMiddleware(http.HandlerFunc(bookingHandler.GetBalance)))
	mux.HandleFunc("GET /api/tokens/packages", bookingHandler.ListPackages)
	mux.HandleFunc("GET /api/tokens/ping-pong-packages", bookingHandler.ListPingPongPackages)
	mux.Handle("GET /api/tokens/transactions", api.AuthMiddleware(http.HandlerFunc(bookingHandler.ListTransactions)))
	mux.Handle("POST /api/tokens/purchase", api.AuthMiddleware(http.HandlerFunc(bookingHandler.PurchaseTokens)))

	// Booking Endpoints
	mux.Handle("POST /api/bookings", api.AuthMiddleware(http.HandlerFunc(bookingHandler.CreateBooking)))
	mux.Handle("DELETE /api/bookings/{id}", api.AuthMiddleware(http.HandlerFunc(bookingHandler.CancelBooking)))

	// Grading Endpoints
	mux.HandleFunc("GET /api/grading/tracks", gradingHandler.GetTracksSyllabus)
	mux.Handle("GET /api/grading/exams/my-history", api.AuthMiddleware(http.HandlerFunc(gradingHandler.GetHistory)))
	mux.Handle("POST /api/grading/exams", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(gradingHandler.SubmitEvaluation))))

	// Forums, Media & Guides Endpoints
	mux.HandleFunc("GET /api/forums/boards", forumsMediaHandler.ListBoards)
	mux.HandleFunc("GET /api/forums/topics/{id}", forumsMediaHandler.ListTopics)
	mux.HandleFunc("GET /api/forums/topics/{id}/posts", forumsMediaHandler.ListPosts)
	mux.Handle("POST /api/forums/topics", api.AuthMiddleware(http.HandlerFunc(forumsMediaHandler.CreateTopic)))
	mux.Handle("POST /api/forums/posts", api.AuthMiddleware(http.HandlerFunc(forumsMediaHandler.CreatePost)))
	mux.Handle("GET /api/media", api.AuthMiddleware(http.HandlerFunc(forumsMediaHandler.GetMedia)))
	mux.HandleFunc("GET /api/guides", forumsMediaHandler.GetGuides)

	// Register Routes
	mux.HandleFunc("GET /api/health", handleHealth(db))
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.Handle("GET /api/auth/me", api.AuthMiddleware(http.HandlerFunc(authHandler.Profile)))

	// User, Family & Memberships Endpoints
	mux.Handle("GET /api/users/family", api.AuthMiddleware(http.HandlerFunc(userHandler.GetFamily)))
	mux.Handle("POST /api/users/family", api.AuthMiddleware(http.HandlerFunc(userHandler.CreateFamily)))
	mux.Handle("POST /api/users/family/members", api.AuthMiddleware(http.HandlerFunc(userHandler.AddFamilyMember)))
	mux.Handle("GET /api/users/memberships", api.AuthMiddleware(http.HandlerFunc(userHandler.GetMemberships)))
	mux.Handle("POST /api/users/memberships", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(userHandler.CreateMembership))))

	// Leaderboard Endpoints
	mux.HandleFunc("GET /api/users/leaderboard/weekly", userHandler.GetWeeklyLeaderboard)
	mux.HandleFunc("GET /api/users/leaderboard/overall", userHandler.GetOverallLeaderboard)

	// Terms Endpoints
	mux.HandleFunc("GET /api/terms", termHandler.List)
	mux.HandleFunc("GET /api/terms/{id}", termHandler.Get)
	mux.Handle("POST /api/terms", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(termHandler.Create))))
	mux.Handle("PUT /api/terms/{id}", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(termHandler.Update))))
	mux.Handle("DELETE /api/terms/{id}", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(termHandler.Delete))))

	// Term Breaks Endpoints
	mux.HandleFunc("GET /api/terms/{term_id}/breaks", termHandler.ListBreaks)
	mux.Handle("POST /api/terms/{term_id}/breaks", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(termHandler.CreateBreak))))
	mux.Handle("DELETE /api/terms/breaks/{id}", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(termHandler.DeleteBreak))))

	// Classes Endpoints
	mux.HandleFunc("GET /api/classes", classHandler.List)
	mux.HandleFunc("GET /api/classes/{id}", classHandler.Get)
	mux.Handle("POST /api/classes", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(classHandler.Create))))
	mux.Handle("PUT /api/classes/{id}", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(classHandler.Update))))
	mux.Handle("DELETE /api/classes/{id}", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(classHandler.Delete))))

	// Halls & Event Types Endpoints
	mux.HandleFunc("GET /api/halls", classHandler.ListHalls)
	mux.HandleFunc("GET /api/event-types", classHandler.ListEventTypes)

	// Class Occurrences Endpoints
	mux.HandleFunc("GET /api/occurrences", classHandler.ListOccurrences)
	mux.Handle("PUT /api/occurrences/{id}", api.AuthMiddleware(api.RequireRole("instructor", "admin")(http.HandlerFunc(classHandler.UpdateOccurrence))))
	mux.Handle("POST /api/occurrences/{id}/attendance", api.AuthMiddleware(api.RequireRole("instructor", "admin")(http.HandlerFunc(classHandler.SubmitAttendance))))

	// Recurring Auto-reservations Endpoints
	mux.Handle("GET /api/reservations", api.AuthMiddleware(http.HandlerFunc(classHandler.GetAutoReservations)))
	mux.Handle("POST /api/reservations", api.AuthMiddleware(http.HandlerFunc(classHandler.CreateAutoReservation)))
	mux.Handle("DELETE /api/reservations/{id}", api.AuthMiddleware(http.HandlerFunc(classHandler.DeleteAutoReservation)))

	// Class Notes Endpoints
	mux.Handle("GET /api/classes/{id}/notes", api.AuthMiddleware(http.HandlerFunc(classHandler.GetClassNote)))
	mux.Handle("POST /api/classes/{id}/notes", api.AuthMiddleware(http.HandlerFunc(classHandler.SaveClassNote)))

	// Wrap Mux with Middleware (CORS & Logging/Recovery)
	handler := loggingMiddleware(corsMiddleware(mux))

	// 4. Server Configuration
	port := getEnv("PORT", "8081")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 5. Graceful Shutdown Setup
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("API Server listening on http://localhost:%s", port)
		serverErrors <- server.ListenAndServe()
	}()

	// Block until a signal or server error occurs
	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	case sig := <-shutdownChan:
		log.Printf("Received signal: %v. Initiating graceful shutdown...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v. Forcing close...", err)
			if err := server.Close(); err != nil {
				log.Printf("Error closing server: %v", err)
			}
		}
		log.Println("Server gracefully stopped.")
	}
}

// Handler: API Health Check
func handleHealth(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status := "OK"
		dbError := ""
		if err := db.Ping(); err != nil {
			status = "Degraded"
			dbError = fmt.Sprintf("Database ping failed: %v", err)
		}

		resp := map[string]string{
			"status":    status,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		if dbError != "" {
			resp["database_error"] = dbError
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		json.NewEncoder(w).Encode(resp)
	}
}

// Middleware: CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:3000"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Middleware: Logger & Panic Recovery
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC recovered: %v", rec)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
			}
		}()

		next.ServeHTTP(w, r)
		log.Printf("%s %s %s took %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

// Helper: Get Environment Variable with Default
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
