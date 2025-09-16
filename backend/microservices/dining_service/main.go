package main

import (
	"context"
	"dining/handler"
	"dining/repo"
	"dining/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	router := mux.NewRouter()

	// Repository Init
	repository, err := repo.NewDiningRepo()
	if err != nil {
		log.Fatal("Creating repository error: ", err)
	}

	// Service Init
	diningService := service.NewDiningService(repository)

	// Handler Init
	diningHandler := handler.NewDiningHandler(*diningService)

	// Rute
	router.HandleFunc("/api/canteens/", diningHandler.GetAllCanteens).Methods(http.MethodGet)
	router.HandleFunc("/api/canteens/{id}", diningHandler.GetCanteen).Methods(http.MethodGet)
	router.HandleFunc("/api/canteens/{id}", diningHandler.DeleteCanteen).Methods(http.MethodDelete)
	router.HandleFunc("/api/canteens/", diningHandler.CreateCanteen).Methods(http.MethodPost)
	router.HandleFunc("/api/canteens/popular-meals/{id}", diningHandler.GetPopularMeals).Methods(http.MethodGet)
	router.HandleFunc("/api/canteens/meal-history/{id}", diningHandler.GetMealHistory).Methods(http.MethodGet)

	router.HandleFunc("/api/menus/{id}", diningHandler.GetMenusByCanteenID).Methods(http.MethodGet)
	router.HandleFunc("/api/menus/{id}", diningHandler.DeleteMenu).Methods(http.MethodDelete)
	router.HandleFunc("/api/menus/", diningHandler.CreateMenu).Methods(http.MethodPost)

	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}), // Angular frontend
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8001"
	}

	server := http.Server{
		Addr:         ":" + port,
		Handler:      corsObj(router), // <-- ovde koristimo CORS middleware
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server u go rutini
	go func() {
		log.Println("server_starting on :" + port)
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("service_shutting_down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
