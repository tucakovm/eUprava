package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"users_module/handlers"
	"users_module/repositories"
	"users_module/services"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Router
	router := mux.NewRouter()

	// Repo init
	repo, err := repositories.NewUserRepository()
	if err != nil {
		log.Fatal("repository init error: ", err)
	}

	// Service init
	jwtSecret := mustEnv("JWT_SECRET")
	authSvc := services.NewAuthService(*repo, jwtSecret)

	// Handler init
	authHandler := handlers.NewAuthHandler(authSvc)

	// Routes
	router.Handle("/api/register", http.HandlerFunc(authHandler.Register)).Methods(http.MethodPost)
	router.Handle("/api/login", http.HandlerFunc(authHandler.Login)).Methods(http.MethodPost)

	router.Handle("/api/users/{id}", http.HandlerFunc(authHandler.GetUser)).Methods(http.MethodGet)

	// Wrap with CORS middleware
	handler := withCORS(router)

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("user-server starting on port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-quit
	log.Println("user-server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("shutdown error:", err)
	}
	log.Println("user-server stopped")
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

// withCORS dodaje potrebne CORS headere
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Ako je preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
