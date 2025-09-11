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

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	router := mux.NewRouter()

	//Repository Init
	repository, err := repo.NewDiningRepo()
	if err != nil {
		log.Fatal("Creating repository error: ", err)
	}

	////Create tables for DB
	//err = repository.Migrate()
	//if err != nil {
	//	return
	//}

	//Service Init
	diningService := service.NewDiningService(repository)

	//Handler Init
	diningHandler := handler.NewDiningHandler(*diningService)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router.Handle("/api/canteen/", http.HandlerFunc(diningHandler.GetAllCanteens)).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8001"
	}

	server := http.Server{
		Addr:         ":" + port,
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("server_starting")
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service_shutting_down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
