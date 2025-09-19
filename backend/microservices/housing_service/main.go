package main

import (
	"context"
	"housing/handler"
	"housing/repository"
	"housing/service"
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

	// === Repository init (konekcija na DB) ===
	repositor, err := repository.NewHousingRepo()
	if err != nil {
		log.Fatal("Creating repository error: ", err)
	}
	defer repositor.Close()

	// === Service init ===
	svcs := service.New(
		repositor.DB,
		repository.NewDomRepo(),
		repository.NewSobaRepo(),
		repository.NewStudentRepo(),
		repository.NewRecRepo(),
		repository.NewKvarRepo(),
		repository.NewStudentskaKarticaRepo(), // NOVO: repo za studentske kartice
	)

	// === Handler init (housing) ===
	hh := handler.NewHousingHandler(svcs)

	// === Routes (housing) ===
	// Doms
	router.Handle("/api/housing/doms", http.HandlerFunc(hh.ListDomovi)).Methods(http.MethodGet) // svi domovi
	router.Handle("/api/housing/dom", http.HandlerFunc(hh.GetDom)).Methods(http.MethodGet)      // jedan dom po ID-u (query param id)
	// Students
	router.Handle("/api/housing/students", http.HandlerFunc(hh.CreateStudent)).Methods(http.MethodPost)
	router.Handle("/api/housing/students/release", http.HandlerFunc(hh.ReleaseStudentRoom)).Methods(http.MethodPost)

	// Studentska kartica (NOVO)
	router.Handle("/api/housing/students/cards", http.HandlerFunc(hh.CreateStudentCardIfMissing)).Methods(http.MethodPost) // create-if-missing
	router.Handle("/api/housing/students/cards", http.HandlerFunc(hh.GetStudentCard)).Methods(http.MethodGet)              // get by studentId (query param)
	router.Handle("/api/housing/students/cards/balance", http.HandlerFunc(hh.UpdateStudentCardBalance)).Methods(http.MethodPost)

	// Rooms
	router.Handle("/api/housing/rooms", http.HandlerFunc(hh.GetRoom)).Methods(http.MethodGet)
	router.Handle("/api/housing/rooms/detail", http.HandlerFunc(hh.GetRoomDetail)).Methods(http.MethodGet)
	router.Handle("/api/housing/rooms/assign", http.HandlerFunc(hh.AssignStudentToRoom)).Methods(http.MethodPost)
	router.Handle("/api/housing/rooms/free", http.HandlerFunc(hh.ListFreeRooms)).Methods(http.MethodGet) // NOVO: slobodne sobe
	router.Handle("/api/housing/rooms/checkStudent/{userId}", http.HandlerFunc(hh.IsStudentAssignedToAnySoba)).Methods(http.MethodGet)
	router.Handle("/api/housing/rooms/meal-history/", http.HandlerFunc(hh.GetRoomMealHistory)).Methods(http.MethodPost)

	// Reviews & Faults
	router.Handle("/api/housing/rooms/reviews", http.HandlerFunc(hh.AddRoomReview)).Methods(http.MethodPost)
	router.Handle("/api/housing/rooms/faults", http.HandlerFunc(hh.ReportFault)).Methods(http.MethodPost)
	router.Handle("/api/housing/faults/status", http.HandlerFunc(hh.ChangeFaultStatus)).Methods(http.MethodPost)

	// === Server setup ===
	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	appHandler := withCORS(router)

	server := http.Server{
		Addr:         ":" + port,
		Handler:      appHandler,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("server_starting on :" + port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
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
